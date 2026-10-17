#!/bin/bash
# Merges module-qualified duplicate schemas into one schema under the plain
# name. The API serves two models with the same class name from different
# modules, and which one keeps the plain name (and which gets a
# sandbox__schemas__<module>__ prefix, split into -Input and -Output) changes
# between responses with no spec version change. That flips the generated Go
# type names. The duplicates are identical, so every reference now points at
# the plain name. The patch fails when the duplicates ever differ.
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

# shellcheck disable=SC2016 # jq variables, not shell expansions
readonly MERGE_FILTER='
def base: sub("^.*__"; "") | sub("-(Input|Output)$"; "");
(.components.schemas | keys | map(select(test("__")) | base) | unique) as $bases
| reduce $bases[] as $name (.;
    [.components.schemas | to_entries[]
      | select((.key | base) == $name and (.key == $name or (.key | test("__"))))]
      as $group
    | if ($group | map(.value) | unique | length) != 1 then
        error("schemas named \($name) differ; review the merge patch")
      else . end
    | ($group | map("#/components/schemas/" + .key)) as $refs
    | .components.schemas |= (
        with_entries(select(.key as $k | $group | map(.key) | index($k) | not))
        + {($name): $group[0].value}
      )
    | walk(
        if type == "object" and (.["$ref"]? | type == "string")
          and (.["$ref"] as $r | $refs | index($r)) then
          .["$ref"] = "#/components/schemas/" + $name
        else . end
      )
  )
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

	jq "$MERGE_FILTER" "$SOURCE_SPEC" >"$tmp"
	chmod 0644 "$tmp"
	mv -- "$tmp" "$OUTPUT_SPEC"
	log INFO "merged module-qualified duplicate schemas"
}

main "$@"
