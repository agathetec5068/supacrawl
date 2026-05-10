package backup

import (
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"filippo.io/age"
	"github.com/davemorin/supacrawl/internal/store"
)

const ManifestName = "manifest.json"

type Writer struct {
	RepoPath  string
	Recipient string
}

type Puller struct {
	RepoPath string
	Identity string
}

type Manifest struct {
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	Tool      string    `json:"tool"`
	Shards    []Shard   `json:"shards"`
}

type Shard struct {
	Table       string `json:"table"`
	File        string `json:"file"`
	Rows        int64  `json:"rows"`
	PlainSHA256 string `json:"plain_sha256"`
}

type WriteResult struct {
	RepoPath     string    `json:"repo_path"`
	ManifestPath string    `json:"manifest_path"`
	CreatedAt    time.Time `json:"created_at"`
	Shards       int       `json:"shards"`
	Rows         int64     `json:"rows"`
}

type PullResult struct {
	RepoPath string `json:"repo_path"`
	OutDir   string `json:"out_dir"`
	Shards   int    `json:"shards"`
	Rows     int64  `json:"rows"`
}

func GenerateIdentity() (*age.X25519Identity, error) {
	return age.GenerateX25519Identity()
}

func (w Writer) Write(ctx context.Context, st *store.Store) (WriteResult, error) {
	if strings.TrimSpace(w.RepoPath) == "" {
		return WriteResult{}, fmt.Errorf("backup repo path is empty")
	}
	recipient, err := age.ParseX25519Recipient(strings.TrimSpace(w.Recipient))
	if err != nil {
		return WriteResult{}, err
	}
	shardDir := filepath.Join(w.RepoPath, "shards")
	if err := os.MkdirAll(shardDir, 0o755); err != nil {
		return WriteResult{}, err
	}
	manifest := Manifest{
		Version:   1,
		CreatedAt: time.Now().UTC(),
		Tool:      "supacrawl",
	}
	var totalRows int64
	for _, table := range store.ArchiveTableNames {
		select {
		case <-ctx.Done():
			return WriteResult{}, ctx.Err()
		default:
		}
		shard, err := w.writeTable(ctx, st, shardDir, table, recipient)
		if err != nil {
			return WriteResult{}, err
		}
		manifest.Shards = append(manifest.Shards, shard)
		totalRows += shard.Rows
	}
	manifestPath := filepath.Join(w.RepoPath, ManifestName)
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return WriteResult{}, err
	}
	if err := os.WriteFile(manifestPath, append(data, '\n'), 0o644); err != nil {
		return WriteResult{}, err
	}
	return WriteResult{
		RepoPath:     w.RepoPath,
		ManifestPath: manifestPath,
		CreatedAt:    manifest.CreatedAt,
		Shards:       len(manifest.Shards),
		Rows:         totalRows,
	}, nil
}

func (w Writer) writeTable(ctx context.Context, st *store.Store, shardDir string, table string, recipient age.Recipient) (Shard, error) {
	fileName := table + ".jsonl.gz.age"
	path := filepath.Join(shardDir, fileName)
	out, err := os.Create(path)
	if err != nil {
		return Shard{}, err
	}
	defer out.Close()
	encrypted, err := age.Encrypt(out, recipient)
	if err != nil {
		return Shard{}, err
	}
	gz := gzip.NewWriter(encrypted)
	hash := sha256.New()
	writer := io.MultiWriter(gz, hash)
	rows, writeErr := st.WriteTableJSONL(ctx, table, writer)
	closeGzipErr := gz.Close()
	closeAgeErr := encrypted.Close()
	if writeErr != nil {
		return Shard{}, writeErr
	}
	if closeGzipErr != nil {
		return Shard{}, closeGzipErr
	}
	if closeAgeErr != nil {
		return Shard{}, closeAgeErr
	}
	return Shard{
		Table:       table,
		File:        filepath.ToSlash(filepath.Join("shards", fileName)),
		Rows:        rows,
		PlainSHA256: hex.EncodeToString(hash.Sum(nil)),
	}, nil
}

func ReadManifest(repoPath string) (Manifest, string, error) {
	manifestPath := filepath.Join(repoPath, ManifestName)
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return Manifest{}, manifestPath, err
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return Manifest{}, manifestPath, err
	}
	return manifest, manifestPath, nil
}

func (p Puller) Pull(ctx context.Context, outDir string) (PullResult, error) {
	if strings.TrimSpace(outDir) == "" {
		return PullResult{}, fmt.Errorf("output directory is empty")
	}
	identity, err := age.ParseX25519Identity(strings.TrimSpace(p.Identity))
	if err != nil {
		return PullResult{}, err
	}
	manifest, _, err := ReadManifest(p.RepoPath)
	if err != nil {
		return PullResult{}, err
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return PullResult{}, err
	}
	var totalRows int64
	for _, shard := range manifest.Shards {
		select {
		case <-ctx.Done():
			return PullResult{}, ctx.Err()
		default:
		}
		if err := p.pullShard(shard, outDir, identity); err != nil {
			return PullResult{}, err
		}
		totalRows += shard.Rows
	}
	return PullResult{RepoPath: p.RepoPath, OutDir: outDir, Shards: len(manifest.Shards), Rows: totalRows}, nil
}

func (p Puller) pullShard(shard Shard, outDir string, identity age.Identity) error {
	inputPath := filepath.Join(p.RepoPath, filepath.FromSlash(shard.File))
	in, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer in.Close()
	decrypted, err := age.Decrypt(in, identity)
	if err != nil {
		return err
	}
	outPath := filepath.Join(outDir, shard.Table+".jsonl.gz")
	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, decrypted)
	return err
}
