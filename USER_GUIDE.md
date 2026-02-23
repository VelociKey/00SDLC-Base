# Olympus Assurance: User Guide

This guide provides instructions for testing the fleet, finding problems, and reporting them for correction.

### Quick Start: Running the Validator

1.  **Ensure Bridges are Running**: Before validating, the specific cluster bridges (e.g., `VaultBridge`, `StorageBridge`) must be active on their respective ports.
2.  **Execute Validation**:
    ```powershell
    pushd OlympusAssurance\10000-Autonomous-Actors\900-WorkstationValidator
    go run .
    popd
    ```
3.  **Interpret Results**:
    *   **PASS**: The cluster is responsive and the proto-contracts are valid.
    *   **❌ FAIL**: Indicates a connection timeout, port conflict, or structural mismatch in the ConnectRPC response.

### How to Find and Report Problems

1.  **Dependency Errors**: If `go build` fails with "not in go.mod", run `.\sovereign.exe init` from the root. If it persists, report it as a **"Module Synchronization Conflict"**.
2.  **API Mismatches**: If the validator fails with "undefined: v1.File...", the proto-generated code has likely changed. Ask for an **"API Realignment"**.
3.  **MCP Failures**: If `SovereignSaaS` fails to start, check for **"Library Port Incompatibilities"** (e.g., changes in the `mcp-go` SDK).

### Asking for Fixes
When a failure occurs, provide the exact output from the `WorkstationValidator` and ask: *"Correct the cluster misalignment in [Cluster Name] detected by the Assurance Suite."*
