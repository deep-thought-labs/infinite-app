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

- [\#758](https://github.com/cosmos/evm/pull/758) Cleanup precompiles abi.json.
- [\#810](https://github.com/cosmos/evm/pull/810) Fix integration test suite to resolve lock contention problem from external app injection
- [\#811](https://github.com/cosmos/evm/pull/811) Use sdk's DefaultBondDenom for default evm denom in genesis.
- [\#823](https://github.com/cosmos/evm/pull/823) Remove authz dependency from test suite and EvmApp interface
- [\#829](https://github.com/cosmos/evm/pull/829) Seperate test app interface
- [\#968](https://github.com/cosmos/evm/pull/968) Use normal gas config in ibc transfer to prevent potential DoS attack
- [\#1029](https://github.com/cosmos/evm/pull/1029) Mark EvmCoinInfo.Decimals field as deprecated
- [\#1079](https://github.com/cosmos/evm/pull/1079) Access Control List is now case-insensitive
- [\#1108](https://github.com/cosmos/evm/pull/1108) Add json-rpc http server request body limit.

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

