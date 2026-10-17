# go-klavis

[![Go Reference](https://pkg.go.dev/badge/github.com/psyb0t/go-klavis.svg)](https://pkg.go.dev/github.com/psyb0t/go-klavis)
[![CI](https://github.com/psyb0t/go-klavis/actions/workflows/pipeline.yml/badge.svg?branch=main)](https://github.com/psyb0t/go-klavis/actions/workflows/pipeline.yml)
[![version](https://raw.githubusercontent.com/psyb0t/go-klavis/badges/version.svg)](https://github.com/psyb0t/go-klavis/tags)
[![license](https://raw.githubusercontent.com/psyb0t/go-klavis/badges/license.svg)](LICENSE)
[![imported by](https://raw.githubusercontent.com/psyb0t/go-klavis/badges/importers.svg)](https://github.com/psyb0t/go-klavis/blob/badges/importers.md)

Unofficial Go client for the [Klavis AI API](https://www.klavis.ai), generated with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) from the [OpenAPI document Klavis AI publishes](https://api.klavis.ai/openapi.json). I am not affiliated with Klavis AI.

The package is `klavis`, in the repository root. Everything in `client.gen.go` is generated. Don't edit it by hand, regenerate it.

## Install

```sh
go get github.com/psyb0t/go-klavis
```

## Usage

Every endpoint in the spec has a typed method on `ClientWithResponses`. Responses decode into `JSON200` and the other status fields.

```go
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-klavis"
)

func main() {
	auth := klavis.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization",
				"Bearer "+os.Getenv("KLAVIS_API_KEY"))

			return nil
		},
	)

	client, err := klavis.NewClientWithResponses(
		"https://api.klavis.ai",
		auth,
	)
	if err != nil {
		panic(err)
	}

	res, err := client.GetAllMcpServersWithResponse(context.Background())
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
```

The same code is in `example_test.go`, so `make test` checks that it compiles.

The client covers all 203 paths and 212 operations: MCP servers and instances, tool calls, Strata, user integrations, OAuth, white labeling, sandboxes and universes. Method names come from the spec's operation IDs, including the longer names on sandbox and Strata endpoints.

## Auth

The spec's `HTTPBearer` scheme takes your Klavis AI API key in `Authorization: Bearer <key>`. The request editor above reads it from `KLAVIS_API_KEY`. Authentication is explicit. Some OAuth endpoints have no security declaration in the spec.

Pick the regional base URL for your account:

| Region | Base URL |
|---|---|
| US | `https://api.klavis.ai` |
| EU | `https://api.eu.klavis.ai` |

Check `JSON200` before using the decoded response. Tool calls and tool lists also return a `Success` field and optional `Error`, so HTTP 200 can still report failure.

## Regenerating

The spec is not committed. `make generate` downloads the current OpenAPI document into the gitignored `api/` directory and regenerates the client:

```sh
make generate
```

`gen.go` declares the steps, so `make generate` is just `go generate ./...`:

- `scripts/fetch-spec.sh` downloads `https://api.klavis.ai/openapi.json` to `api/openapi.json`.
- `scripts/patch-spec.sh` writes `api/openapi.fixed.json`, merging identical schemas the API publishes under module-qualified names (`sandbox__schemas__<module>__<Name>`) into one schema under the plain name. The API changes which copy gets the plain name between responses, so without the merge the generated type names would flip. The patch fails when the copies differ.
- oapi-codegen reads the OpenAPI 3.1 document with `config.yaml`. No schema downconversion is needed.
- The config maps unformatted `number` fields to `float64`, resolves type name collisions and prefixes enum values to avoid duplicate Go identifiers.

`.github/workflows/regen.yml` does the same every day in CI. When the new spec changes the client, it runs `make lint`, `make test` and `make sec`. If a check fails, Claude edits the code and the checks run again, up to three times, and a review step reads the result. A passing change is committed with a CHANGELOG entry and tagged as the next release: a patch version for fixes and additions, a minor version for breaking changes.

## Agent integrations

The [skill](.agents/skills/go-klavis) teaches a coding agent how to use go-klavis when it writes Go. It works in any agent that reads `.agents/skills/`, and installs natively in the clients below.

### Claude Code

```bash
claude plugin marketplace add psyb0t/agents
claude plugin install go-klavis@psyb0t
```

### Codex

```bash
codex plugin marketplace add psyb0t/agents
codex plugin add go-klavis@psyb0t
```

Codex also picks the skill up with no install in any checkout of this repo.

### OpenClaw

The skill is published to ClawHub on every release:

```bash
openclaw skills install @psyb0t/go-klavis
```

## Development

Every command runs in the `Dockerfile.dev` image. The host needs only Docker and Make.

```sh
make help      # list targets
make dep       # go mod tidy + vendor
make lint      # shellcheck, shfmt, go fix, golangci-lint
make test      # go test -race
make sec       # gitleaks, govulncheck, semgrep
make generate  # download the spec and regenerate
```

## License

[MIT](LICENSE). The generated code comes from Klavis AI's published API specification. The specification itself is not redistributed here. It has no license declaration or terms URL.
