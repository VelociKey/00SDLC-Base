package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"connectrpc.com/connect"

	// Cluster Clients
	vaultv1 "olympus.fleet/00SDLC/OlympusGCP-Vault/gen/v1/vault"
	"olympus.fleet/00SDLC/OlympusGCP-Vault/gen/v1/vault/vaultv1connect"

	storagev1 "olympus.fleet/00SDLC/OlympusGCP-Storage/gen/v1/storage"
	"olympus.fleet/00SDLC/OlympusGCP-Storage/gen/v1/storage/storagev1connect"

	intv1 "olympus.fleet/00SDLC/OlympusGCP-Intelligence/gen/v1/intelligence"
	"olympus.fleet/00SDLC/OlympusGCP-Intelligence/gen/v1/intelligence/intelligencev1connect"

	finopsv1 "olympus.fleet/00SDLC/OlympusGCP-FinOps/gen/v1/finops"
	"olympus.fleet/00SDLC/OlympusGCP-FinOps/gen/v1/finops/finopsv1connect"

	computev1 "olympus.fleet/00SDLC/OlympusGCP-Compute/gen/v1/compute"
	"olympus.fleet/00SDLC/OlympusGCP-Compute/gen/v1/compute/computev1connect"

	datav1 "olympus.fleet/00SDLC/OlympusGCP-Data/gen/v1/data"
	"olympus.fleet/00SDLC/OlympusGCP-Data/gen/v1/data/datav1connect"

	eventsv1 "olympus.fleet/00SDLC/OlympusGCP-Events/gen/v1/events"
	"olympus.fleet/00SDLC/OlympusGCP-Events/gen/v1/events/eventsv1connect"
)

func main() {
	fmt.Println("🛡️  Olympus Assurance: Starting Fleet Validation Suite")
	fmt.Println("-----------------------------------------------------")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	validateVault(ctx)
	validateStorage(ctx)
	validateIntelligence(ctx)
	validateFinOps(ctx)
	validateCompute(ctx)
	validateData(ctx)
	validateEvents(ctx)

	fmt.Println("\n✅ Fleet Validation Complete.")
}

func validateVault(ctx context.Context) {
	fmt.Print("Checking Vault Cluster... ")
	client := vaultv1connect.NewVaultServiceClient(http.DefaultClient, "http://localhost:8092")
	res, err := client.VaultRead(ctx, connect.NewRequest(&vaultv1.VaultReadRequest{
		Key: "olympus.fleet/00SDLC/Olympus2/90000-Enablement-Labs/000-Tools/fleet/validation/token",
	}))
	if err != nil {
		fmt.Printf("❌ FAIL: %v\n", err)
		return
	}
	fmt.Printf("PASS (Secret retrieved: %s)\n", res.Msg.Value)
}

func validateStorage(ctx context.Context) {
	fmt.Print("Checking Storage Cluster... ")
	client := storagev1connect.NewStorageServiceClient(http.DefaultClient, "http://localhost:8091")
	_, err := client.GetDownloadURL(ctx, connect.NewRequest(&storagev1.GetDownloadURLRequest{
		Bucket: "assurance-validation",
		Name:   "health.txt",
	}))
	if err != nil {
		fmt.Printf("❌ FAIL: %v\n", err)
		return
	}
	fmt.Println("PASS")
}

func validateIntelligence(ctx context.Context) {
	fmt.Print("Checking Intelligence Cluster... ")
	client := intelligencev1connect.NewIntelligenceServiceClient(http.DefaultClient, "http://localhost:8087")
	_, err := client.Predict(ctx, connect.NewRequest(&intv1.PredictRequest{
		Prompt: "Validation Ping",
	}))
	if err != nil {
		fmt.Printf("❌ FAIL: %v\n", err)
		return
	}
	fmt.Println("PASS")
}

func validateFinOps(ctx context.Context) {
	fmt.Print("Checking FinOps Cluster... ")
	client := finopsv1connect.NewFinOpsServiceClient(http.DefaultClient, "http://localhost:8098")
	_, err := client.ValidateBudget(ctx, connect.NewRequest(&finopsv1.ValidateBudgetRequest{
		ProjectId:       "fleet-assurance",
		RequestedAmount: 50.0,
	}))
	if err != nil {
		fmt.Printf("❌ FAIL: %v\n", err)
		return
	}
	fmt.Println("PASS")
}

func validateCompute(ctx context.Context) {
	fmt.Print("Checking Compute Cluster... ")
	client := computev1connect.NewComputeServiceClient(http.DefaultClient, "http://localhost:8095")
	_, err := client.TriggerFunction(ctx, connect.NewRequest(&computev1.TriggerFunctionRequest{
		FunctionName: "health-check",
	}))
	if err != nil {
		fmt.Printf("❌ FAIL: %v\n", err)
		return
	}
	fmt.Println("PASS")
}

func validateData(ctx context.Context) {
	fmt.Print("Checking Data Cluster... ")
	client := datav1connect.NewDataServiceClient(http.DefaultClient, "http://localhost:8093")
	_, err := client.Query(ctx, connect.NewRequest(&datav1.QueryRequest{
		Collection: "assurance",
		Filter:     "health=true",
	}))
	if err != nil {
		fmt.Printf("❌ FAIL: %v\n", err)
		return
	}
	fmt.Println("PASS")
}

func validateEvents(ctx context.Context) {
	fmt.Print("Checking Events Cluster... ")
	client := eventsv1connect.NewEventsServiceClient(http.DefaultClient, "http://localhost:8094")
	_, err := client.Publish(ctx, connect.NewRequest(&eventsv1.PublishRequest{
		Topic: "assurance-ping",
		Data:  "ping",
	}))
	if err != nil {
		fmt.Printf("❌ FAIL: %v\n", err)
		return
	}
	fmt.Println("PASS")
}
