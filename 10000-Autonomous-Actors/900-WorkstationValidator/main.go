package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"connectrpc.com/connect"

	// Cluster Clients
	vaultv1 "OlympusGCP-Vault/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/vault/v1"
	"OlympusGCP-Vault/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/vault/v1/vaultv1connect"

	storagev1 "OlympusGCP-Storage/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/storage/v1"
	"OlympusGCP-Storage/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/storage/v1/storagev1connect"

	intv1 "OlympusGCP-Intelligence/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/intelligence/v1"
	"OlympusGCP-Intelligence/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/intelligence/v1/intelligencev1connect"

	finopsv1 "OlympusGCP-FinOps/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/finops/v1"
	"OlympusGCP-FinOps/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/finops/v1/finopsv1connect"

	computev1 "OlympusGCP-Compute/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/compute/v1"
	"OlympusGCP-Compute/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/compute/v1/computev1connect"

	datav1 "OlympusGCP-Data/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/data/v1"
	"OlympusGCP-Data/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/data/v1/datav1connect"

	eventsv1 "OlympusGCP-Events/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/events/v1"
	"OlympusGCP-Events/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/events/v1/eventsv1connect"
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
		Key: "fleet/validation/token",
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
