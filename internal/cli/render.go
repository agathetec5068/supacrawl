package cli

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/davemorin/supacrawl/internal/storage"
	"github.com/davemorin/supacrawl/internal/store"
)

func renderText(w io.Writer, title string, value any) error {
	fmt.Fprintf(w, "%s\n", title)
	switch v := value.(type) {
	case store.Status:
		renderStatus(w, v)
	case store.Report:
		renderReport(w, v)
	case []store.SearchResult:
		renderSearchResults(w, v)
	case store.SQLResult:
		renderSQL(w, v)
	case store.ArchiveSize:
		renderArchiveSize(w, v)
	case storage.DownloadStats:
		renderDownloadStats(w, v)
	default:
		renderAny(w, value)
	}
	return nil
}

func renderStatus(w io.Writer, status store.Status) {
	if status.ProjectID != "" {
		fmt.Fprintf(w, "  project: %s\n", status.ProjectID)
	}
	if status.DatabaseName != "" {
		fmt.Fprintf(w, "  database: %s\n", status.DatabaseName)
	}
	if !status.CollectedAt.IsZero() {
		fmt.Fprintf(w, "  collected: %s\n", status.CollectedAt.Format(time.RFC3339))
	}
	fmt.Fprintf(w, "  schemas: %d\n", status.Schemas)
	fmt.Fprintf(w, "  tables: %d\n", status.Tables)
	fmt.Fprintf(w, "  columns: %d\n", status.Columns)
	fmt.Fprintf(w, "  policies: %d\n", status.Policies)
	fmt.Fprintf(w, "  functions: %d\n", status.Functions)
	fmt.Fprintf(w, "  triggers: %d\n", status.Triggers)
	fmt.Fprintf(w, "  extensions: %d\n", status.Extensions)
	fmt.Fprintf(w, "  storage buckets: %d\n", status.StorageBuckets)
	fmt.Fprintf(w, "  storage objects: %d\n", status.StorageObjects)
	fmt.Fprintf(w, "  data tables: %d\n", status.DataTables)
	fmt.Fprintf(w, "  data rows: %d\n", status.DataRows)
}

func renderReport(w io.Writer, report store.Report) {
	renderStatus(w, report.Status)
	if len(report.SchemaTables) > 0 {
		fmt.Fprintln(w, "\n  schemas")
		for _, row := range report.SchemaTables {
			fmt.Fprintf(w, "    %s: %d tables\n", row.Schema, row.Tables)
		}
	}
	if len(report.PolicyTables) > 0 {
		fmt.Fprintln(w, "\n  policy tables")
		for _, row := range report.PolicyTables {
			fmt.Fprintf(w, "    %s.%s: %d policies\n", row.Schema, row.Table, row.Policies)
		}
	}
}

func renderSearchResults(w io.Writer, results []store.SearchResult) {
	if len(results) == 0 {
		fmt.Fprintln(w, "  no matches")
		return
	}
	for _, result := range results {
		fmt.Fprintf(w, "  [%s] %s\n", result.Kind, result.Title)
		if strings.TrimSpace(result.Snippet) != "" {
			fmt.Fprintf(w, "    %s\n", compactWhitespace(result.Snippet))
		}
	}
}

func renderSQL(w io.Writer, result store.SQLResult) {
	if len(result.Columns) == 0 {
		fmt.Fprintln(w, "  no columns")
		return
	}
	fmt.Fprintf(w, "  %s\n", strings.Join(result.Columns, "\t"))
	for _, row := range result.Rows {
		fmt.Fprintf(w, "  %s\n", strings.Join(row, "\t"))
	}
}

func renderArchiveSize(w io.Writer, size store.ArchiveSize) {
	fmt.Fprintf(w, "  db path: %s\n", size.DBPath)
	fmt.Fprintf(w, "  file bytes: %s\n", humanBytes(size.FileBytes))
	if size.WALBytes > 0 {
		fmt.Fprintf(w, "  wal bytes: %s\n", humanBytes(size.WALBytes))
	}
	fmt.Fprintf(w, "  sqlite bytes: %s\n", humanBytes(size.SQLiteBytes))
	fmt.Fprintf(w, "  row json bytes: %s\n", humanBytes(size.RowJSONBytes))
	if len(size.SourceTables) > 0 {
		fmt.Fprintln(w, "\n  largest source tables")
		for _, row := range size.SourceTables {
			fmt.Fprintf(w, "    %s.%s: %s across %d rows\n", row.Schema, row.Table, humanBytes(row.RowJSONBytes), row.Rows)
		}
	}
}

func renderDownloadStats(w io.Writer, stats storage.DownloadStats) {
	fmt.Fprintf(w, "  objects: %d\n", stats.Objects)
	fmt.Fprintf(w, "  bytes: %s\n", humanBytes(stats.Bytes))
	fmt.Fprintf(w, "  skipped: %d\n", stats.Skipped)
}

func renderAny(w io.Writer, value any) {
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Map {
		iter := rv.MapRange()
		for iter.Next() {
			fmt.Fprintf(w, "  %v: %v\n", iter.Key(), iter.Value())
		}
		return
	}
	fmt.Fprintf(w, "  %v\n", value)
}

func compactWhitespace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func humanBytes(value int64) string {
	const unit = 1024
	if value < unit {
		return strconv.FormatInt(value, 10) + " B"
	}
	div, exp := int64(unit), 0
	for n := value / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(value)/float64(div), "KMGTPE"[exp])
}
