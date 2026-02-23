package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/connect"
	obsv1 "OlympusGCP-Observability/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/observability/v1"
	observabilityv1connect "OlympusGCP-Observability/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/observability/v1/observabilityv1connect"
	"Olympus2/90000-Enablement-Labs/P0000-pkg/000-whisper"
)

func main() {
	logger := whisper.New("SaaS-Observability-Mock", "gcp_observability.lpsv")
	client := observabilityv1connect.NewObservabilityServiceClient(
		http.DefaultClient,
		"http://localhost:8097",
	)

	slog.Info("SaaS-Observability-Mock: Initializing...")

	// Phase 1: Cloud Logging
	slog.Info("SaaS-Observability-Mock: Phase 1 - Cloud Logging")
	logReq := &obsv1.WriteLogRequest{
		LogName:  "system-events",
		Severity: "INFO",
		Message:  "User login successful",
		Labels:   map[string]string{"env": "prod"},
	}
	start1 := time.Now()
	_, err := client.WriteLog(context.Background(), connect.NewRequest(logReq))
	if err != nil {
		slog.Error("FAIL: Could not write log", "error", err)
		return
	}
	logger.Log("Logging_Write", "VERIFIED", "system-events", "INFO", time.Since(start1))

	// Phase 2: Cloud Monitoring
	slog.Info("SaaS-Observability-Mock: Phase 2 - Cloud Monitoring")
	metricReq := &obsv1.RecordMetricRequest{
		MetricType: "custom.googleapis.com/request_count",
		Value:      1.0,
		Labels:     map[string]string{"status": "200"},
	}
	start2 := time.Now()
	_, err = client.RecordMetric(context.Background(), connect.NewRequest(metricReq))
	if err != nil {
		slog.Error("FAIL: Could not record metric", "error", err)
		return
	}
	logger.Log("Monitoring_Record", "VERIFIED", metricReq.MetricType, "1.0", time.Since(start2))

	// Phase 3: Cloud Trace
	slog.Info("SaaS-Observability-Mock: Phase 3 - Cloud Trace")
	spanReq := &obsv1.StartSpanRequest{
		Name:    "DatabaseQuery",
		TraceId: "trace-xyz-123",
	}
	start3 := time.Now()
	res3, err := client.StartSpan(context.Background(), connect.NewRequest(spanReq))
	if err != nil {
		slog.Error("FAIL: Could not start span", "error", err)
		return
	}
	logger.Log("Trace_Span_Start", "VERIFIED", spanReq.Name, res3.Msg.SpanId, time.Since(start3))

	_, err = client.EndSpan(context.Background(), connect.NewRequest(&obsv1.EndSpanRequest{SpanId: res3.Msg.SpanId}))
	if err != nil {
		slog.Error("FAIL: Could not end span", "error", err)
		return
	}
	logger.Log("Trace_Span_End", "VERIFIED", res3.Msg.SpanId, "SUCCESS", 0)

	slog.Info("SaaS-Observability-Mock: All Observability Assurance checks passed.")
	time.Sleep(1 * time.Second)
}
