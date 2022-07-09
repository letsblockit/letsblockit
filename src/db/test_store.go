package db

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/gommon/random"
	"github.com/stretchr/testify/require"
)

const (
	schemaNamePrefix   = "test_"
	templateSchemaName = schemaNamePrefix + "template"

	forkSchemaQuery      = `SELECT clone_schema($1,$2);`
	createCloneFuncQuery = `-- From https://wiki.postgresql.org/wiki/Clone_schema
CREATE OR REPLACE FUNCTION clone_schema(source_schema text, dest_schema text) RETURNS void AS
$BODY$
DECLARE
  object text;
  buffer text;
BEGIN
    EXECUTE 'CREATE SCHEMA ' || dest_schema ;

    FOR object IN
        SELECT table_name::text FROM information_schema.tables WHERE table_schema = source_schema
    LOOP
        buffer := dest_schema || '.' || object;
        EXECUTE 'CREATE TABLE ' || buffer || ' (LIKE ' || source_schema || '.' || object || ' INCLUDING CONSTRAINTS INCLUDING INDEXES INCLUDING DEFAULTS)';
        EXECUTE 'INSERT INTO ' || buffer || '(SELECT * FROM ' || source_schema || '.' || object || ')';
    END LOOP;

END;
$BODY$
LANGUAGE plpgsql VOLATILE;`

	// Cannot use prepared statements for these operations, will use fmt.Sprintf
	createSchemaPattern = `CREATE SCHEMA IF NOT EXISTS %s;`
	dropSchemaPattern   = `DROP SCHEMA %s CASCADE;`
)

var initTemplateOnce sync.Once

func GetTestDatabaseURL() string {
	if url := os.Getenv("TEST_DATABASE_URL"); url != "" {
		return url
	}
	return "postgresql:///lbitests"
}

// NewTestStore instantiates an ephemeral PG schema for the test duration
// and ensures it is destroyed on test completion.
func NewTestStore(t *testing.T) Store {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := pgx.Connect(ctx, GetTestDatabaseURL())
	require.NoError(t, err)

	initTemplateOnce.Do(func() { // Create test_template and clone_schema only once
		mustExec(t, conn, createCloneFuncQuery)
		mustExec(t, conn, createSchemaPattern, templateSchemaName)
		require.NoError(t, Migrate(buildConnString(templateSchemaName)))
	})

	for i := 0; i < 30; i++ { // Retry schema fork to be resilient to name collision
		clonedSchema := schemaNamePrefix + random.String(16, random.Alphabetic)
		if _, err = conn.Exec(ctx, forkSchemaQuery, templateSchemaName, clonedSchema); err == nil {
			t.Cleanup(func() {
				mustExec(t, conn, dropSchemaPattern, clonedSchema)
				require.NoError(t, conn.Close(context.Background()))
			})
			store, err := Connect(buildConnString(clonedSchema))
			require.NoError(t, err)
			return store
		}
	}
	require.NoError(t, err)
	return nil
}

// mustExec uses string interpolation to execute schema operations
func mustExec(t *testing.T, p *pgx.Conn, pattern string, args ...interface{}) {
	_, err := p.Exec(context.Background(), fmt.Sprintf(pattern, args...))
	require.NoError(t, err)
}

func buildConnString(schema string) string {
	return fmt.Sprintf("%s?search_path=%s", GetTestDatabaseURL(), schema)
}
