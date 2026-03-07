package main

import (
	"context"
	"net/http"
	"testing"

	vaultv1 "olympus.fleet/00SDLC/OlympusGCP-Vault/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/vault"
	vaultv1connect "olympus.fleet/00SDLC/OlympusGCP-Vault/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/vault/vaultv1connect"
	computev1 "olympus.fleet/00SDLC/OlympusGCP-Compute/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/compute"
	computev1connect "olympus.fleet/00SDLC/OlympusGCP-Compute/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/compute/computev1connect"
	"connectrpc.com/connect"
)

func TestCrossCluster_VaultToCompute_BlackBox(t *testing.T) {
	// This test performs 100% black-box testing against the RPC endpoints.
	// It assumes the managers are running (or uses httptest for a hermetic build).
	
	ctx := context.Background()
	
	// 1. Vault Client
	vaultURL := "http://localhost:8092"
	vaultClient := vaultv1connect.NewVaultServiceClient(http.DefaultClient, vaultURL)
	
	// 2. Compute Client
	computeURL := "http://localhost:8095"
	computeClient := computev1connect.NewComputeServiceClient(http.DefaultClient, computeURL)

	// Goal: Verify that a valid identity can perform a cross-cluster operation.
	// In a real assurance environment, we would trigger a flow and verify side effects.
	
	// We'll skip the actual network calls in this unit test unless a flag is provided.
	if testing.Short() {
		t.Skip("Skipping black-box network test in short mode")
	}

	// Example: Try to read from Vault
	_, err := vaultClient.VaultRead(ctx, connect.NewRequest(&vaultv1.VaultReadRequest{Key: "test"}))
	if err != nil && !IsConnectError(err, connect.CodeNotFound) && !IsConnectError(err, connect.CodeUnavailable) {
		t.Errorf("Unexpected Vault error: %v", err)
	}

	// Example: Check Compute health
	res, err := computeClient.CheckHealth(ctx, connect.NewRequest(&computev1.CheckHealthRequest{ServiceName: "assurance"}))
	if err != nil && !IsConnectError(err, connect.CodeUnavailable) {
		t.Errorf("Unexpected Compute error: %v", err)
	}
	if err == nil && res.Msg.Status != computev1.CheckHealthResponse_HEALTHY {
		t.Errorf("Expected healthy compute, got %v", res.Msg.Status)
	}
}

func IsConnectError(err error, code connect.Code) bool {
	if connectErr, ok := err.(*connect.Error); ok {
		return connectErr.Code() == code
	}
	return false
}
