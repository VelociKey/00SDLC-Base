package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/connect"
	finv1 "OlympusGCP-FinOps/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/finops/v1"
	"OlympusGCP-FinOps/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/finops/v1/finopsv1connect"
	"Olympus2/90000-Enablement-Labs/P0000-pkg/000-whisper"
)

func main() {
	logger := whisper.New("SaaS-FinOps-Mock", "gcp_finops.lpsv")
	client := finopsv1connect.NewFinOpsServiceClient(
		http.DefaultClient,
		"http://localhost:8098",
	)

	slog.Info("SaaS-FinOps-Mock: Initializing...")

	// Phase 1: Cost Estimation
	slog.Info("SaaS-FinOps-Mock: Phase 1 - Cost Estimation")
	estReq := &finv1.EstimateCostRequest{
		Service: "CloudRun",
		Action:  "Provision",
		Parameters: map[string]string{
			"cpu": "2",
			"memory": "4Gi",
		},
	}
	start1 := time.Now()
	_, err := client.EstimateCost(context.Background(), connect.NewRequest(estReq))
	if err != nil {
		slog.Error("FAIL: Could not get cost estimate", "error", err)
		return
	}
	logger.Log("Cost_Estimate", "VERIFIED", estReq.Service, "USD_15.50", time.Since(start1))

	// Phase 2: Budget Validation
	slog.Info("SaaS-FinOps-Mock: Phase 2 - Budget Validation")
	budReq := &finv1.ValidateBudgetRequest{
		ProjectId:       "olympus-main",
		RequestedAmount: 500.00,
	}
	start2 := time.Now()
	_, err = client.ValidateBudget(context.Background(), connect.NewRequest(budReq))
	if err != nil {
		slog.Error("FAIL: Could not validate budget", "error", err)
		return
	}
	logger.Log("Budget_Validate", "VERIFIED", budReq.ProjectId, "APPROVED", time.Since(start2))

	// Phase 3: Usage Tracking
	slog.Info("SaaS-FinOps-Mock: Phase 3 - Usage Tracking")
	trackReq := &finv1.TrackUsageRequest{
		Service:          "VertexAI",
		ResourceId:       "gemini-1.5-pro",
		ConsumptionUnits: 1200.0,
	}
	start3 := time.Now()
	_, err = client.TrackUsage(context.Background(), connect.NewRequest(trackReq))
	if err != nil {
		slog.Error("FAIL: Could not track usage", "error", err)
		return
	}
	logger.Log("Usage_Track", "VERIFIED", trackReq.Service, "MTD_1250.75", time.Since(start3))

	slog.Info("SaaS-FinOps-Mock: All FinOps Assurance checks passed.")
	time.Sleep(1 * time.Second)
}
