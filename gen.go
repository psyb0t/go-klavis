// Package klavis is a Go client for the Klavis AI API,
// generated from the OpenAPI document Klavis AI publishes.
package klavis

//go:generate bash scripts/fetch-spec.sh
//go:generate bash scripts/patch-spec.sh
//go:generate go tool oapi-codegen --config=config.yaml api/openapi.fixed.json
