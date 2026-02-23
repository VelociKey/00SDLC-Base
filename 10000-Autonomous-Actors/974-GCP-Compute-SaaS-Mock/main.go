package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/connect"
	computev1 "OlympusGCP-Compute/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/compute/v1"
	"OlympusGCP-Compute/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/compute/v1/computev1connect"
	"Olympus2/90000-Enablement-Labs/P0000-pkg/000-whisper"
)

func main() {
	logger := whisper.New("SaaS-Compute-Mock", "gcp_compute.lpsv")
	client := computev1connect.NewComputeServiceClient(
		http.DefaultClient,
		"http://localhost:8095",
	)

	slog.Info("SaaS-Compute-Mock: Initializing...")

	// Phase 1: Cloud Run
	slog.Info("SaaS-Compute-Mock: Phase 1 - Cloud Run Provisioning")
	runReq := &computev1.RunServiceRequest{
		ServiceName: "order-processor",
		Image:       "gcr.io/olympus/order-proc:v1.0",
		EnvVars: map[string]string{
			"MODE": "PROD",
		},
	}
	start1 := time.Now()
	res1, err := client.RunService(context.Background(), connect.NewRequest(runReq))
	if err != nil {
		slog.Error("FAIL: Could not provision Cloud Run service", "error", err)
		return
	}
	logger.Log("CloudRun_Provision", "VERIFIED", runReq.ServiceName, res1.Msg.EndpointUrl, time.Since(start1))

	// Phase 2: Cloud Functions
	slog.Info("SaaS-Compute-Mock: Phase 2 - Cloud Functions Trigger")
	funcReq := &computev1.TriggerFunctionRequest{
		FunctionName: "thumbnail-gen",
		Payload:      `{"bucket": "images", "file": "logo.png"}`,
	}
	start2 := time.Now()
	res2, err := client.TriggerFunction(context.Background(), connect.NewRequest(funcReq))
	if err != nil {
		slog.Error("FAIL: Could not trigger Cloud Function", "error", err)
		return
	}
	logger.Log("CloudFunc_Trigger", "VERIFIED", funcReq.FunctionName, res2.Msg.Result, time.Since(start2))

	slog.Info("SaaS-Compute-Mock: All Compute Assurance checks passed.")
	time.Sleep(1 * time.Second)
}
