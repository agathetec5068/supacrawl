package search

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildFTSQueryQuotesTokens(t *testing.T) {
	require.Equal(t, `"auth" AND "users"`, BuildFTSQuery("auth users"))
	require.Equal(t, `"auth.users"`, BuildFTSQuery("auth.users"))
	require.Equal(t, `"user" AND "table"`, BuildFTSQuery("user-table"))
}
