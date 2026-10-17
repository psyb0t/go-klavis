---
name: go-klavis
description: "Use go-klavis when Go code needs Klavis AI MCP servers, tool calls, Strata, user integrations, OAuth or sandboxes, or already imports github.com/psyb0t/go-klavis."
homepage: https://github.com/psyb0t/go-klavis
user-invocable: true
metadata:
  openclaw:
    emoji: "🔌"
    requires:
      bins:
        - go
      env:
        - KLAVIS_API_KEY
permissions:
  filesystem: "Guides edits to the Go source of the current project. The library reads and writes no files."
  shell: "The Go toolchain (`go get`, `go build`, `go test`) to add and compile the dependency."
  network: "`go get` fetches the module from your Go module proxy. At runtime the code you write talks HTTPS to api.klavis.ai or api.eu.klavis.ai. Tool operations can make Klavis contact the requested MCP server."
---

# go-klavis

Unofficial, typed Go client for the Klavis AI API, in package `klavis` at the module root.

Read [references/setup.md](references/setup.md) before adding the dependency or choosing an operation. It covers authentication, regional hosts, operation names, pagination and response fields.

## Security and safety

- Read the Klavis API key from `KLAVIS_API_KEY`. Set `Authorization: Bearer <key>` in a request editor. Never hardcode, commit or log credentials, downstream authentication headers or returned authentication data.
- Choose `https://api.klavis.ai` or `https://api.eu.klavis.ai` for the account's region. Do not send credentials to another host.
- Tool calls can change external services. Creating or deleting instances, changing authentication, user deletion and sandbox execution also act on remote resources. Call them only when the user requested the action.
- Verify MCP server URLs and tool arguments before sending downstream credentials or running a tool.

## When to use

- Go code needs to discover MCP servers, list or call tools, or manage Klavis instances, Strata or user integrations.
- Go code needs Klavis OAuth, sandbox or universe REST endpoints.
- The project already imports `github.com/psyb0t/go-klavis`.

## When not to use

- The code is not Go. This module ships only a Go package.
- You need a client maintained by Klavis AI. This one is unofficial and not affiliated with Klavis AI.
- You want to change the generated client itself. Read `AGENTS.md` in the repository.

## Workflow

1. Check the installed module version and generated signature before editing. Operation IDs determine names, including long sandbox and Strata methods.
2. Read the regional host and key from application configuration. Add `WithRequestEditorFn` for Bearer auth and `WithHTTPClient` for the application's timeout policy.
3. Use `NewClientWithResponses` and the operation's `...WithResponse` method. Inspect transport errors, HTTP status and decoded fields before using the result.
4. For tool lists and calls, check `Success` and `Error` even on HTTP 200. A tool result can also carry `IsError`.
5. Compile the integration and test it against a local HTTP server. Finish when the request method, path, authentication and response handling match the generated API. Do not call the live service merely to verify compilation.

Optional fields are often pointers. Check for nil before dereferencing. The error return covers transport and decoding failures; API errors require checking `StatusCode()`, the status-specific fields and `Body`.
