# Olympus Assurance: Historical Record & Learning Log

This document tracks all problems encountered during the iteration of the sovereign fleet and the specific corrections applied. This serves as a "learning" dataset for fleet agents.

| Date | Issue | Component | Correction Applied | Key Learning |
| :--- | :--- | :--- | :--- | :--- |
| 2026-02-18 | Broken Imports in FinOps/Obs | ConnectRPC Generated Code | Qualified imports to physical paths (e.g., `OlympusGCP-FinOps/...`) | Generated ConnectRPC code defaults to short paths; must be manually/programmatically qualified in prefix-free fleets. |
| 2026-02-18 | Brace Mismatch in Fabric | `OlympusGrammar/fabric` | Properly nested `if` statement and restored closing braces. | Automatic regex-based refactors can orphan braces; manual validation is mandatory after global replacements. |
| 2026-02-18 | Incompatible MCP API | `950-SovereignSaaS` | Updated `AddResource` signature to return `([]mcp.ResourceContents, error)`. | MCP SDK version updates frequently change method signatures; refer to module source via `go list -m` when compiler fails. |
| 2026-02-18 | Transitive Resolvability | Root dependencies | Created local "shell" modules for `pretty`, `text`, `go-internal`. | Transitive dependencies of sovereign internal libraries must be locally shadowed to prevent network-based resolution failures. |
| 2026-02-18 | Vault Client Mismatch | `WorkstationValidator` | Renamed `GetSecret` to `VaultRead` to match actual service definition. | Don't assume method names based on domain; always verify against the generated `.connect.go` interface. |
| 2026-02-22 | George Structural Drift | `George` | Reconstituted the workspace from `P0000-pkg` to `10000-Autonomous-Actors` using the **Genesis Grammar**. | Rapid agentic iteration can lead to architectural drift; periodic "Reconstitution Cycles" are mandatory to maintain fleet alignment. |
