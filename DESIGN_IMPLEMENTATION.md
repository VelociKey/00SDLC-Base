# Olympus Assurance: Design & Implementation

This workspace utilizes a decoupled probe-and-reference architecture to validate the fleet.

### Architectural Components:

#### 1. Cluster Validation Engine (`WorkstationValidator`)
*   **Path**: `10000-Autonomous-Actors/900-WorkstationValidator`
*   **Implementation**: A Go-native binary that acts as a fleet consumer.
*   **Strategy**: 
    *   Initiates `ConnectRPC` unary calls directly to local cluster ports (e.g., `:8092` for Vault).
    *   Validates response payloads against expected structural schemas.
    *   Provides color-coded terminal feedback for immediate human/agent assessment.

#### 2. Reference Service (`SovereignSaaS`)
*   **Path**: `10000-Autonomous-Actors/950-SovereignSaaS`
*   **Implementation**: A Model Citizen service fronted by an MCP (Model Context Protocol) server.
*   **Capabilities**:
    *   **Tools**: Exposes operational health via `get_service_status`.
    *   **Resources**: Demonstrates `sovereign://` URI-based data surfacing for configuration.
*   **Standardization**: Uses the shared `Olympus2/P0000-pkg/000-mcp-bridge` for fleet-wide behavioral consistency.

### Fleet Synchronization:
All tools in this workspace are managed by the authoritative `sovereign` CLI tool, ensuring that `go.mod` files always contain local `replace` directives for cross-cluster dependencies.
