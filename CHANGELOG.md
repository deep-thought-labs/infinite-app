# CHANGELOG

## Infinite Improbability Drive

## v0.1.0

Initial identity integration on Cosmos EVM v0.4.0.

- Renamed example chain directory from evmd to infinited
- Added project README, NOTICE, and LICENSE branding for Infinite Drive
- Established Infinite network identity documentation (chain name, token, bech32)

## v0.2.0

Upstream synchronization milestones through the Cosmos EVM v0.6.0 era.

### UPSTREAM INTEGRATION

- Synchronized codebase with cosmos/evm through the v0.6.0 release line.
- Retained Infinite-specific ante handling where upstream removed equivalent logic.

### FEATURES

- Hyperlane integrated as Cosmos SDK modules (x/core, x/warp).
- Bank extension: governance-only MsgSetDenomMetadata for SDK bank denom metadata.

### CONTINUOUS INTEGRATION

- Build and test gates aligned with upstream workflows.
- System test upgrade harness and Docker coverage on macOS where applicable.
- Coverage CI matrix for unit and integration paths.
- Docs-only PR path filters to skip heavy jobs when appropriate.

### IMPROVEMENTS

- local_node.sh PATH and GOPATH for installed binary.
- Fork maintenance documentation for upstream remote and verification.

