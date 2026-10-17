#!/bin/bash
# Normalize provider descriptions emitted as Go comments to repository style.
# Replace dash punctuation and the first-person plural sandbox note. Keep
# schema constraints, examples, defaults, enum values and operations unchanged.
set -euo pipefail
trap 'on_error "$?" "$LINENO"' ERR

LOG_FILE="${LOG_FILE:-/tmp/$(basename "$0" .sh).log}"
exec > >(tee -a "$LOG_FILE") 2>&1

SPEC_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/api"
readonly SPEC_DIR
readonly SOURCE_SPEC="${SPEC_DIR}/openapi.json"
readonly OUTPUT_SPEC="${SPEC_DIR}/openapi.fixed.json"

on_error() {
	local exit_code="$1"
	local line="$2"

	printf '{"level":"ERROR","file":"%s","line":%d,"msg":"command failed","exit_code":%d}\n' \
		"${BASH_SOURCE[1]}" "$line" "$exit_code" >&2
}

log() {
	local level="$1"
	shift

	jq -nc \
		--arg time "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" \
		--arg level "$level" \
		--arg file "${BASH_SOURCE[1]##*/}" \
		--arg line "${BASH_LINENO[0]}" \
		--arg func "${FUNCNAME[1]:-main}" \
		--arg msg "$*" \
		'{time:$time,level:$level,file:$file,line:($line|tonumber),func:$func,msg:$msg}' >&2
}

readonly PROSE_FILTER='
def prose:
  gsub("[ \t]*[\u2013\u2014][ \t]*"; ", ")
  | sub("we typically work with"; "typically use");
def normalize:
  if type == "object" then
    with_entries(
      if .key == "example" or .key == "examples" or .key == "default"
        or .key == "enum" or .key == "const" then .
      elif (.key == "description" or .key == "summary")
        and (.value | type == "string") then .value |= prose
      else .value |= normalize end
    )
  elif type == "array" then map(normalize)
  else . end;
normalize
'

main() {
	if [[ ! -f "$SOURCE_SPEC" ]]; then
		log ERROR "source spec missing path=${SOURCE_SPEC}"
		exit 1
	fi

	local tmp
	tmp="$(mktemp "${SPEC_DIR}/.openapi.fixed.XXXXXX")"
	# shellcheck disable=SC2064 # expand tmp before the local leaves scope
	trap "rm -f -- '${tmp}'" EXIT

	jq "$PROSE_FILTER" "$SOURCE_SPEC" >"$tmp"
	chmod 0644 "$tmp"
	mv -- "$tmp" "$OUTPUT_SPEC"
	log INFO "normalized description prose, retained OpenAPI 3.1 schemas"
}

main "$@"
