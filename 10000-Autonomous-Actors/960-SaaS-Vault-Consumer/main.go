package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"OlympusAssurance/90000-Enablement-Labs/P0000-pkg/000-assurance-client"
)

func main() {
	fmt.Println("🚀 SaaS Application: Starting up...")
	
	cwd, _ := os.Getwd()
	// Target the Vault Bridge
	bridgePath := filepath.Join(cwd, "../../../OlympusGCP-Vault/20000-Context-Bridges/900-VaultBridge/VaultBridge.exe")

	client, err := assurance.NewBridgeClient(bridgePath)
	if err != nil {
		log.Fatalf("SaaS Error: Could not connect to Vault Bridge: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("SaaS Logic: Requesting Database Credential from Vault Bridge...")
	
	// We call the 'vault_read' tool provided by the bridge
	dbPass, err := client.CallTool(ctx, "vault_read", map[string]interface{}{
		"key": "production/databases/main/password",
	})
	if err != nil {
		log.Fatalf("SaaS Error: Failed to retrieve DB password: %v", err)
	}

	fmt.Printf("SaaS Success: DB Connection Established using password: [%s]\n", dbPass)
	
	fmt.Println("SaaS Logic: Performing business operation...")
	fmt.Println("✅ SaaS Application Finished Successfully.")
}
