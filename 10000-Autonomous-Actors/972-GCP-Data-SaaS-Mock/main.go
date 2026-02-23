package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	datav1 "OlympusGCP-Data/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/data/v1"
	"OlympusGCP-Data/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/data/v1/datav1connect"
	"Olympus2/90000-Enablement-Labs/P0000-pkg/000-whisper"
)

func main() {
	logger := whisper.New("SaaS-Data-Mock", "gcp_data.lpsv")
	client := datav1connect.NewDataServiceClient(
		http.DefaultClient,
		"http://localhost:8093",
	)

	log.Printf("INFO SaaS-Data-Mock: Initializing...")

	// Phase 1: Firestore Simulation
	log.Printf("INFO SaaS-Data-Mock: Phase 1 - Firestore Document Ops")
	upsertReq := &datav1.UpsertRequest{
		Collection: "users",
		DocId:      "user_123",
		DataJson:   `{"name": "Olympus Dev", "status": "active"}`, 
	}
	start1 := time.Now()
	_, err := client.Upsert(context.Background(), connect.NewRequest(upsertReq))
	if err != nil {
		log.Fatalf("FAIL: Could not upsert document: %v", err)
	}
	logger.Log("Firestore_Upsert", "VERIFIED", "users/user_123", "", time.Since(start1))

	// Phase 2: BigTable Simulation
	log.Printf("INFO SaaS-Data-Mock: Phase 2 - BigTable Row Ops")
	writeReq := &datav1.WriteRowRequest{
		Table:        "telemetry",
		RowKey:       "sensor_001",
		ColumnFamily: "metrics",
		Column:       "cpu_usage",
		Value:        []byte("45.5"),
	}
	start2 := time.Now()
	_, err = client.WriteRow(context.Background(), connect.NewRequest(writeReq))
	if err != nil {
		log.Fatalf("FAIL: Could not write row: %v", err)
	}
	logger.Log("BigTable_Write", "VERIFIED", "telemetry/sensor_001", "", time.Since(start2))

	log.Printf("INFO SaaS-Data-Mock: All Data Assurance checks passed.")
	time.Sleep(1 * time.Second) // Let whisper flush
}
