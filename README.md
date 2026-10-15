# NetAssist

NetAssist is a Go foundation for a self-hosted networking management platform. Its initial CLI provides version, local status, and configuration workflows while the architecture reserves clear ownership for agents, API services, discovery, diagnostics, monitoring, SIEM, and a future networking engine.

## Phase 1 scope

- Cobra CLI: `netassist version`, `netassist status`, `netassist config init`, and `netassist config show`.
- YAML configuration with JSON-oriented structured logging defaults.
- Module seams with no active network scans, device collection, security detection, or traffic tunneling.
- SIEM model boundary for normalized event intake and storage.

## Local use

```bash
go run ./cmd/netassist version
go run ./cmd/netassist config init
go run ./cmd/netassist status
go test ./...
```

## Architecture

| Package | Responsibility |
| --- | --- |
| `internal/cli` | User-facing command interface and local configuration |
| `internal/agent` | Future cross-platform device enrollment and telemetry |
| `internal/api` | Future control-plane transport |
| `internal/discovery` | Authorized asset discovery and enrichment |
| `internal/diagnostics` | User-authorized troubleshooting checks |
| `internal/monitoring` | Availability and performance measurement |
| `internal/siem` | Security event ingestion, normalization, storage and alert seams |
| `internal/engine` | Future independently designed connectivity engine |

## Compatibility

This repository was empty when Phase 1 began. No LabLink code, WireGuard behavior, deployment configuration, domains, or existing branding was changed or removed.

Any future connectivity-engine work must be separately designed, independently security reviewed, and explicitly approved before replacing an existing transport such as WireGuard.
