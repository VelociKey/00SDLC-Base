package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"
	"os"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel"

	intelligencev1 "olympus.fleet/00SDLC/OlympusGCP-Intelligence/50000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/intelligence"
	"olympus.fleet/00SDLC/OlympusGCP-Intelligence/50000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/intelligence/intelligencev1connect"
	"olympus.fleet/00SDLC/Olympus2/90000-Enablement-Labs/90200-Logic-Libraries/220-Whisper"
	econotel "olympus.fleet/00SDLC/Olympus2/90000-Enablement-Labs/90200-Logic-Libraries/120-Econotel"
)

func main() {
	// Centralized OTel Output (SaaS Side)
	otelFile, _ := os.Create("C:/aAntigravitySpace/olympus.fleet/00SDLC/Olympus2/C0500-Agent-Intelligence-Outputs/LPSV/gcp_saas_mock_otel.json")
	tp, _ := econotel.InitTracer("SaaS-Intelligence-Mock", otelFile)
	defer tp.Shutdown(context.Background())

	logger := whisper.New("SaaS-Intelligence-Mock", "gcp_intelligence.lpsv")
	client := intelligencev1connect.NewIntelligenceServiceClient(
		http.DefaultClient,
		"http://localhost:8096",
	)

	tracer := otel.GetTracerProvider().Tracer("SaaS-Intelligence-Mock")
	ctx, span := tracer.Start(context.Background(), "IntelligenceAssuranceCycle")
	defer span.End()

	slog.Info("SaaS-Intelligence-Mock: Initializing with OTel context...")

	// Phase 1: Vertex AI Prediction
	slog.Info("SaaS-Intelligence-Mock: Phase 1 - Model Prediction")
	predReq := &intelligencev1.PredictRequest{
		Model:  "gemini-1.5-pro",
		Prompt: "How can OpenTelemetry predict GCP costs?",
	}
	start1 := time.Now()
	res1, err := client.Predict(ctx, connect.NewRequest(predReq))
	if err != nil {
		slog.Error("FAIL: Could not get prediction", "error", err)
		return
	}
	logger.Log("Model_Predict", "VERIFIED", predReq.Model, res1.Msg.Prediction, time.Since(start1))

	slog.Info("SaaS-Intelligence-Mock: All Intelligence Assurance checks passed with Economic Tracing.")
	time.Sleep(2 * time.Second) // Let OTel batcher flush
}
