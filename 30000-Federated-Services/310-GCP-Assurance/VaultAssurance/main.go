package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"

	"Olympus2/90000-Enablement-Labs/P0000-pkg/000-forge-context"
	vaultv1 "OlympusGCP-Vault/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/vault/v1"
	"OlympusGCP-Vault/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/vault/v1/vaultv1connect"
)

func main() {
	fmt.Println("🛡️  SaaS Component: Vault Assurance Logic")
	fmt.Println("----------------------------------------")

	client := vaultv1connect.NewVaultServiceClient(
		http.DefaultClient,
		forgecontext.GetVaultURL(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Scenario: Multi-layered Secret Lifecycle
	fmt.Println("Layer 1: Securely writing application credential...")
	_, err := client.VaultWrite(ctx, connect.NewRequest(&vaultv1.VaultWriteRequest{
		Key:   "assurance/saas/db-pass",
		Value: "super-secure-assurance-token",
	}))
	if err != nil {
		log.Fatalf("FAIL: Could not write secret: %v", err)
	}

	fmt.Println("Layer 2: Verifying credential integrity...")
	resp, err := client.VaultRead(ctx, connect.NewRequest(&vaultv1.VaultReadRequest{
		Key: "assurance/saas/db-pass",
	}))
	if err != nil {
		log.Fatalf("FAIL: Could not read secret: %v", err)
	}
	if resp.Msg.Value != "super-secure-assurance-token" {
		log.Fatalf("FAIL: Secret mismatch. Got %s", resp.Msg.Value)
	}

	fmt.Println("Layer 3: Auditing IAM access for 'saas-worker'...")
	iam, err := client.TestIAMPolicy(ctx, connect.NewRequest(&vaultv1.TestIAMPolicyRequest{
		Identity: "saas-worker",
		Action:   "vault.secrets.get",
		Resource: "assurance/saas/db-pass",
	}))
	if err != nil {
		log.Fatalf("FAIL: IAM check failed: %v", err)
	}
	fmt.Printf("  Access Allowed: %t (Reason: %s)\n", iam.Msg.Allowed, iam.Msg.Reason)

	fmt.Println("\n✅ Vault SaaS Assurance PASSED")
}
