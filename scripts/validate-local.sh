#!/usr/bin/env bash

set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
tmpdir="$(mktemp -d)"
cleanup() {
  rm -rf "$tmpdir"
}
trap cleanup EXIT

cd "$root_dir"

go test ./...
go vet ./...
go build -o "$tmpdir/supacrawl" ./cmd/supacrawl

"$tmpdir/supacrawl" --help >/dev/null
"$tmpdir/supacrawl" version --json >/dev/null
"$tmpdir/supacrawl" metadata --json >/dev/null

"$tmpdir/supacrawl" --config "$tmpdir/config.toml" init \
  --db "$tmpdir/supacrawl.db" \
  --project-id local-smoke >/dev/null

echo "== status =="
"$tmpdir/supacrawl" --json --config "$tmpdir/config.toml" status

echo "== size =="
"$tmpdir/supacrawl" --json --config "$tmpdir/config.toml" size

echo "== sql =="
"$tmpdir/supacrawl" --json --config "$tmpdir/config.toml" sql \
  "select count(*) as tables from tables"

if [[ -n "${SUPABASE_DB_URL:-}" ]]; then
  echo "== doctor =="
  "$tmpdir/supacrawl" --json --config "$tmpdir/config.toml" doctor

  echo "== metadata sync =="
  "$tmpdir/supacrawl" --json --config "$tmpdir/config.toml" sync --no-progress

  echo "== search =="
  "$tmpdir/supacrawl" --json --config "$tmpdir/config.toml" search table --limit 3

  if [[ "${SUPACRAWL_VALIDATE_FULL:-0}" == "1" ]]; then
    echo "== full sync =="
    "$tmpdir/supacrawl" --json --config "$tmpdir/config.toml" sync \
      --full \
      --no-row-fts \
      --no-progress \
      --batch-size "${SUPACRAWL_BATCH_SIZE:-1000}"
  fi
else
  echo "SUPABASE_DB_URL is not set; skipped live doctor/sync checks"
fi

echo "== ok =="

