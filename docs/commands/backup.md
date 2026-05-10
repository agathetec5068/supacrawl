# backup

`backup` writes encrypted local archive shards from the SQLite database and can
decrypt them back into JSONL gzip files.

Generate an age identity:

```bash
supacrawl backup keygen --out ~/.supacrawl/age.key
```

The command writes the private identity to disk with `0600` permissions and only
prints the public recipient.

Configure a backup repo and recipient:

```bash
supacrawl backup init --recipient age1...
```

You can also use environment variables:

```bash
export SUPACRAWL_AGE_RECIPIENT=age1...
supacrawl backup init --recipient-env SUPACRAWL_AGE_RECIPIENT
```

Write a backup:

```bash
supacrawl backup push
```

This creates `manifest.json` plus encrypted shards under
`~/.supacrawl/backups/shards/` by default.

Inspect a backup:

```bash
supacrawl backup status --json
```

Restore encrypted shards to local files:

```bash
supacrawl backup pull --out ./supacrawl-restore
```

`backup push` is local-only. It does not commit, push, or publish the archive.
