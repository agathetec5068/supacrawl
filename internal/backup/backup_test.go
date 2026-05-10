package backup

import (
	"compress/gzip"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/davemorin/supacrawl/internal/postgres"
	"github.com/davemorin/supacrawl/internal/store"
	"github.com/stretchr/testify/require"
)

func TestEncryptedBackupRoundTrip(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	st, err := store.Open(filepath.Join(dir, "supacrawl.db"))
	require.NoError(t, err)
	defer st.Close()

	require.NoError(t, st.BeginDataCopy(ctx, false))
	require.NoError(t, st.PutDataRows(ctx, []postgres.TableRow{
		{Schema: "public", TableName: "companies", RowNumber: 1, JSON: `{"name":"Offline Nexus"}`},
	}, false))
	require.NoError(t, st.FinishDataCopy(ctx, postgres.DataCopyStats{Tables: 1, Rows: 1}, nil))

	identity, err := GenerateIdentity()
	require.NoError(t, err)

	repoPath := filepath.Join(dir, "backup")
	writeResult, err := Writer{RepoPath: repoPath, Recipient: identity.Recipient().String()}.Write(ctx, st)
	require.NoError(t, err)
	require.Equal(t, len(store.ArchiveTableNames), writeResult.Shards)
	require.GreaterOrEqual(t, writeResult.Rows, int64(1))

	manifest, manifestPath, err := ReadManifest(repoPath)
	require.NoError(t, err)
	require.Equal(t, filepath.Join(repoPath, ManifestName), manifestPath)
	require.Len(t, manifest.Shards, len(store.ArchiveTableNames))

	outDir := filepath.Join(dir, "restore")
	pullResult, err := Puller{RepoPath: repoPath, Identity: identity.String()}.Pull(ctx, outDir)
	require.NoError(t, err)
	require.Equal(t, writeResult.Shards, pullResult.Shards)
	require.Equal(t, writeResult.Rows, pullResult.Rows)

	tableRows, err := os.Open(filepath.Join(outDir, "table_rows.jsonl.gz"))
	require.NoError(t, err)
	defer tableRows.Close()
	reader, err := gzip.NewReader(tableRows)
	require.NoError(t, err)
	defer reader.Close()
	restored, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Contains(t, string(restored), "Offline Nexus")

	encrypted, err := os.ReadFile(filepath.Join(repoPath, "shards", "table_rows.jsonl.gz.age"))
	require.NoError(t, err)
	require.NotContains(t, string(encrypted), "Offline Nexus")
}
