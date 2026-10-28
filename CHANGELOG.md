# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.2.4, 2026-10-28

- Updated Semgrep in the development image to 1.182.1 and refreshed its locked dependencies. certifi 2026.10.23 drops the emSign root CAs whose intermediate keys were compromised (GHSA-6q48-6c7w-9gw2).
- Dependabot no longer opens pull requests for `mcp` or the OpenTelemetry packages in `requirements-semgrep.txt`. Semgrep pins their versions itself, so bumping one of them on its own cannot install.
- Updated kin-openapi, which the code generator uses, from 0.142.0 to 0.149.0. The newer version fixes several panics on malformed specs reported in GitHub security advisories.
- The security scan merges its two gitleaks runs into one SARIF run. GitHub code scanning rejects a report with two runs from the same tool, so no scan results had been uploaded.

## v0.2.3, 2026-10-27

- The development image is based on Debian (`golang:1.26.10-trixie`) instead of Alpine. Alpine keeps only the latest build of each package, so the exact pins broke the image build whenever a package was updated. System packages now come from the snapshot.debian.org date the base image was built from, which keeps them fixed for a given base image digest. jq 1.8.2, shellcheck 0.11.0 and shfmt 3.13.1 keep their versions, and Semgrep runs in its own virtual environment.
- Raised the minimum Go version in `go.mod` from 1.26.9 to 1.26.10, which fixes the `crypto/x509` vulnerability GO-2026-6809 found by govulncheck.
- The issue-pull workflow runs at hours derived from the repository name instead of the same four hours as every other repository.

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
