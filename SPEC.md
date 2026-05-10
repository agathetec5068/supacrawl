# supacrawl spec

## Goal

Build a local-first archive for Supabase/Postgres metadata and table data. The first version should help agents and maintainers answer questions like:

- What tables exist?
- Which tables have RLS enabled?
- What policies guard a table?
- What functions, triggers, indexes, and constraints exist?
- Which storage buckets are public?
- Where is a column or function mentioned?
- What source rows exist locally for fast analysis?
- Which copied tables dominate local archive size?
- Is the local archive stale?
- Can the local archive be encrypted and restored?

## Non-Goals For MVP

- Do not store Supabase API keys inside the SQLite archive.
- Do not mutate the remote Supabase project.
- Do not depend on the Supabase dashboard.

## Archive Model

SQLite is the canonical local archive. The MVP stores:

- `project_info`
- `schemas`
- `tables`
- `columns`
- `indexes`
- `constraints`
- `policies`
- `functions`
- `triggers`
- `extensions`
- `storage_buckets`
- `storage_object_stats`
- `table_rows`
- `table_row_fts`
- `crawl_runs`
- `data_copy_runs`
- `search_docs`
- `search_fts`

The crawler replaces the metadata snapshot on each `sync`. `sync --data` and `sync --full` also replace the local table-row corpus. This keeps the first version simple and deterministic.

## Commands

- `init`: write local config
- `doctor`: verify archive and remote Postgres connection
- `sync`: crawl metadata into SQLite
- `sync --data` / `sync --full`: crawl metadata plus source table rows into SQLite JSON rows
- `sync --full --no-row-fts`: copy rows without row full-text index
- `size`: summarize archive size and largest copied source tables
- `export`: export one copied source table as JSONL or CSV
- `storage pull`: download Storage blobs listed in copied `storage.objects`
- `status`: print counts and last collection time
- read commands: accept `--sync auto|always|never` and `--stale-after`
- `report`: summarize schemas and policy coverage
- `search`: FTS search over human-readable metadata documents
- `sql`: read-only SQL against the local archive
- `backup keygen`: create an age identity and print the public recipient
- `backup init`: configure local encrypted backups
- `backup push`: write encrypted archive shards and a manifest
- `backup status`: print backup manifest metadata
- `backup pull`: decrypt backup shards into a restore directory

## Future Work

- `publish`, `subscribe`, and `update` for git-backed snapshots
- terminal UI powered by crawlkit-style explorer patterns
- Management API source for edge functions and project settings
- Supabase Management API settings export
- optional Markdown export
- materialized per-table SQLite tables for richer local SQL
