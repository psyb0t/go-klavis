# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.1, 2026-10-17

- `make sec` no longer fails on the generated client's free-form JSON fields. Semgrep's unsafe-deserialization-interface rule is skipped for the generated client only. It still runs on handwritten code, and every other rule still runs on the generated client.

## v0.1.0, 2026-10-17

- Added `klavis`, a REST client generated with oapi-codegen from Klavis AI's OpenAPI 3.1.0 document, API version 0.1.0, at https://api.klavis.ai/openapi.json. It covers all 203 paths and 212 operations: MCP servers and instances, tool calls, Strata, users, OAuth, white labeling, sandboxes and universes.
- Number fields are `float64` instead of oapi-codegen's default `float32`.
- `make generate` downloads the current spec, normalizes description prose emitted as Go comments and regenerates the client. The patch preserves OpenAPI 3.1 constraints and wire data. The spec is not committed.
