# Changelog

All notable changes to this project will be documented in this file.
The format is based on [Keep a Changelog], and this project adheres to [Semantic Versioning].

## [Unreleased]

## [v1.0.4] - 2023-11-27
### Added
- `MsgCreateIdentityGISTTransferOp`, `MsgCreateIdentityStateTransferOp` messages handling
- `tx` and `height` columns for the `confirmation` and `vote` tables
- `identity_gist_transfer` and `identity_state_transfer` tables to store new identity transfer types
- Hasura migrations
  - `contract_update` entity relation to the `network`
  - `fee_token_management` entity relation to the `network`
  - `identity_default_transfer` entity relation to the `network`
  - `oracle` entity relation to the `network`
  - `confirmation` entity relation to the `transaction` and `block`
  - `vote` entity relation to the `transaction` and `block`
  - `multisig_proposal_vote` entity relation to the `block`
  - `operation` entity relation to the `identity_gist_transfer` and `identity_state_transfer`

### Changed
- `rarimo-core` dependency updated to `v1.1.0`
- `tokenmanager_params` table replaced with `network` table

## [v1.0.3] - 2023-11-03
### Added
- MsgExec support for modules `auth`, `distribution`, `feegrant`, `gov`, `staking`

## [v1.0.2] - 2023-10-24
### Under the hood changes
- Migrated from GitLab to GitHub

### Under the hood changes

- Initiated project

[Unreleased]: https://github.com/rarimo/bdjuno/compare/v1.0.4...HEAD
[v1.0.4]: https://github.com/rarimo/bdjuno/compare/v1.0.3...v1.0.4
[v1.0.3]: https://github.com/rarimo/bdjuno/compare/v1.0.2...v1.0.3
[v1.0.2]: https://github.com/rarimo/bdjuno/releases/tag/v1.0.2
[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html
