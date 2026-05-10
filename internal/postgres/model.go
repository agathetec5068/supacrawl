package postgres

import "time"

type Snapshot struct {
	Project            ProjectInfo         `json:"project"`
	Schemas            []Schema            `json:"schemas"`
	Tables             []Table             `json:"tables"`
	Columns            []Column            `json:"columns"`
	Indexes            []Index             `json:"indexes"`
	Constraints        []Constraint        `json:"constraints"`
	Policies           []Policy            `json:"policies"`
	Functions          []Function          `json:"functions"`
	Triggers           []Trigger           `json:"triggers"`
	Extensions         []Extension         `json:"extensions"`
	StorageBuckets     []StorageBucket     `json:"storage_buckets"`
	StorageObjectStats []StorageObjectStat `json:"storage_object_stats"`
}

type ProjectInfo struct {
	ID            string    `json:"id"`
	DatabaseName  string    `json:"database_name"`
	CurrentUser   string    `json:"current_user"`
	ServerVersion string    `json:"server_version"`
	CollectedAt   time.Time `json:"collected_at"`
}

type Schema struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
	Type  string `json:"type"`
}

type Table struct {
	Schema        string `json:"schema"`
	Name          string `json:"name"`
	Kind          string `json:"kind"`
	Owner         string `json:"owner"`
	Comment       string `json:"comment"`
	RLSEnabled    bool   `json:"rls_enabled"`
	RLSForced     bool   `json:"rls_forced"`
	EstimatedRows int64  `json:"estimated_rows"`
}

type Column struct {
	TableSchema string `json:"table_schema"`
	TableName   string `json:"table_name"`
	Name        string `json:"name"`
	Ordinal     int    `json:"ordinal"`
	DataType    string `json:"data_type"`
	IsNullable  bool   `json:"is_nullable"`
	Default     string `json:"default"`
	Comment     string `json:"comment"`
}

type Index struct {
	Schema     string `json:"schema"`
	TableName  string `json:"table_name"`
	Name       string `json:"name"`
	Definition string `json:"definition"`
}

type Constraint struct {
	Schema     string `json:"schema"`
	TableName  string `json:"table_name"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Definition string `json:"definition"`
}

type Policy struct {
	Schema    string `json:"schema"`
	TableName string `json:"table_name"`
	Name      string `json:"name"`
	Command   string `json:"command"`
	Roles     string `json:"roles"`
	Using     string `json:"using"`
	Check     string `json:"check"`
}

type Function struct {
	Schema          string `json:"schema"`
	Name            string `json:"name"`
	IdentityArgs    string `json:"identity_args"`
	Returns         string `json:"returns"`
	Language        string `json:"language"`
	SecurityDefiner bool   `json:"security_definer"`
	Definition      string `json:"definition"`
}

type Trigger struct {
	Schema       string `json:"schema"`
	TableName    string `json:"table_name"`
	Name         string `json:"name"`
	Timing       string `json:"timing"`
	Events       string `json:"events"`
	FunctionName string `json:"function_name"`
	Definition   string `json:"definition"`
}

type Extension struct {
	Name    string `json:"name"`
	Schema  string `json:"schema"`
	Version string `json:"version"`
	Comment string `json:"comment"`
}

type StorageBucket struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Public           bool   `json:"public"`
	FileSizeLimit    string `json:"file_size_limit"`
	AllowedMimeTypes string `json:"allowed_mime_types"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

type StorageObjectStat struct {
	BucketID    string `json:"bucket_id"`
	ObjectCount int64  `json:"object_count"`
	TotalBytes  int64  `json:"total_bytes"`
}

type TableRow struct {
	Schema    string `json:"schema"`
	TableName string `json:"table_name"`
	RowNumber int64  `json:"row_number"`
	JSON      string `json:"json"`
}

type Counts struct {
	Schemas            int `json:"schemas"`
	Tables             int `json:"tables"`
	Columns            int `json:"columns"`
	Indexes            int `json:"indexes"`
	Constraints        int `json:"constraints"`
	Policies           int `json:"policies"`
	Functions          int `json:"functions"`
	Triggers           int `json:"triggers"`
	Extensions         int `json:"extensions"`
	StorageBuckets     int `json:"storage_buckets"`
	StorageObjectStats int `json:"storage_object_stats"`
}

type DataCopyStats struct {
	Tables int   `json:"tables"`
	Rows   int64 `json:"rows"`
}

type DataCopyProgress struct {
	Schema    string
	TableName string
	Rows      int64
	Done      bool
}

func (s Snapshot) Counts() Counts {
	return Counts{
		Schemas:            len(s.Schemas),
		Tables:             len(s.Tables),
		Columns:            len(s.Columns),
		Indexes:            len(s.Indexes),
		Constraints:        len(s.Constraints),
		Policies:           len(s.Policies),
		Functions:          len(s.Functions),
		Triggers:           len(s.Triggers),
		Extensions:         len(s.Extensions),
		StorageBuckets:     len(s.StorageBuckets),
		StorageObjectStats: len(s.StorageObjectStats),
	}
}
