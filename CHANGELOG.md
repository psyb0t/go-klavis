# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.2.2, 2026-10-17

- Documentation only. The agent skill explains how to look up every endpoint, parameter and model with `go doc` or on pkg.go.dev, and declares its license.
- The skill no longer lists the API credentials under OpenClaw's `requires.env`, which kept OpenClaw from loading it until they were set. It requires only `go`; the credentials are still documented as what your program reads.
- The README points to the full API reference, and the examples in the README and the skill were checked against the current code.

## v0.2.1, 2026-10-17

- Generated type names no longer flip between regenerations. The API publishes the same Google Calendar event model under a plain and a module-qualified name, and which module gets the plain name changes between responses. The spec patch now merges the identical copies, so `GoogleCalendarDataInput`, `GoogleCalendarDataOutput`, `GoogleWorkspaceAtlasDataInput` and `GoogleWorkspaceAtlasDataOutput` all use `GoogleCalendarEvent` for `Events`.
- `SandboxSchemasGoogleCalendarGoogleCalendarEventInput`, `SandboxSchemasGoogleCalendarGoogleCalendarEventOutput`, `SandboxSchemasGoogleWorkspaceAtlasGoogleCalendarEventInput` and `SandboxSchemasGoogleWorkspaceAtlasGoogleCalendarEventOutput` remain as deprecated aliases of `GoogleCalendarEvent`, so code written against v0.1.x and v0.2.0 still compiles.
- Generated doc comments keep the API's own description text. The spec patch no longer rewrites dashes or wording in descriptions.
- `make sec` uses Semgrep 1.181.0. Its pinned requirements now hash every published file, so the development image also builds on arm64.

## v0.2.0, 2026-10-17

- `GoogleWorkspaceAtlasDataInput.Events` and `GoogleWorkspaceAtlasDataOutput.Events` now hold `SandboxSchemasGoogleWorkspaceAtlasGoogleCalendarEventInput` and `SandboxSchemasGoogleWorkspaceAtlasGoogleCalendarEventOutput` instead of `GoogleCalendarEvent`.
- `GoogleCalendarDataInput.Events` and `GoogleCalendarDataOutput.Events` now hold `GoogleCalendarEvent` instead of `SandboxSchemasGoogleCalendarGoogleCalendarEventInput` and `SandboxSchemasGoogleCalendarGoogleCalendarEventOutput`.
- `SandboxSchemasGoogleCalendarGoogleCalendarEventInput` and `SandboxSchemasGoogleCalendarGoogleCalendarEventOutput` are removed. Code that references them must switch to `GoogleCalendarEvent`.
- `SandboxSchemasGoogleWorkspaceAtlasGoogleCalendarEventInput` and `SandboxSchemasGoogleWorkspaceAtlasGoogleCalendarEventOutput` are added. They carry the struct body that `GoogleCalendarEvent` carried for the Google Workspace Atlas sandbox before this release.
- oapi-codegen's name collision resolution assigned the plain `GoogleCalendarEvent` name to a different schema this run. The upstream spec version did not change, so this is a regeneration artifact, not an upstream API update.

## v0.1.1, 2026-10-17

- `make sec` no longer fails on the generated client's free-form JSON fields. Semgrep's unsafe-deserialization-interface rule is skipped for the generated client only. It still runs on handwritten code, and every other rule still runs on the generated client.

## v0.1.0, 2026-10-17

- Added `klavis`, a REST client generated with oapi-codegen from Klavis AI's OpenAPI 3.1.0 document, API version 0.1.0, at https://api.klavis.ai/openapi.json. It covers all 203 paths and 212 operations: MCP servers and instances, tool calls, Strata, users, OAuth, white labeling, sandboxes and universes.
- Number fields are `float64` instead of oapi-codegen's default `float32`.
- `make generate` downloads the current spec, normalizes description prose emitted as Go comments and regenerates the client. The patch preserves OpenAPI 3.1 constraints and wire data. The spec is not committed.
