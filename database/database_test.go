package database

import (
	"blog-go/auth"
	"database/sql"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	if err := Init(":memory:"); err != nil {
		log.Fatal("failed to init test DB:", err)
	}

	code := m.Run()

	if err := Close(); err != nil {
		log.Printf("failed to close test db: %v", err)
	}

	os.Exit(code)
}

func TestBase(t *testing.T) {
	t.Run("check connection", func(t *testing.T) {
		assert.Nil(t, globalDb.Ping())
	})

	t.Run("check tables", func(t *testing.T) {
		requiredTables := []string{"users", "tokens"}

		query := "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?"

		for _, table := range requiredTables {
			var cnt int

			err := queryRow(query, table).Scan(&cnt)
			assert.Nil(t, err)
			assert.Equal(t, 1, cnt, "failed to find table: %s", table)

		}
	})

	t.Run("check indexes", func(t *testing.T) {
		indices := []string{"idx_users_username", "idx_tokens_expires"}

		query := "SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name=?"

		for _, index := range indices {
			var cnt int

			err := queryRow(query, index).Scan(&cnt)
			assert.Nil(t, err)
			assert.Equal(t, 1, cnt)
		}
	})

	t.Run("test crud", func(t *testing.T) {
		err := command(`CREATE TABLE IF NOT EXISTS testtable(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tstring TEXT
		)`, nil)
		assert.Nil(t, err, "failed creating table")

		err = command(`INSERT INTO testtable(tstring) VALUES(?)`, "hello")
		assert.Nil(t, err, "failed inserting")

		var tstring string
		err = queryRow(`SELECT tstring FROM testtable WHERE id=?`, 1).Scan(&tstring)
		assert.Nil(t, err, "failed querying row")
		assert.Equal(t, "hello", tstring)

		err = command(`DELETE FROM testtable WHERE id=?`, 1)
		assert.Nil(t, err, "failed deleting")
	})

	t.Run("transaction rollback", func(t *testing.T) {
		err := command(`CREATE TABLE IF NOT EXISTS test_tx(id INTEGER PRIMARY KEY, name TEXT)`)
		assert.Nil(t, err)

		err = command(`INSERT INTO test_tx(id, name, invalid_column) VALUES(1, 'test', 'fail')`)
		assert.NotNil(t, err, "should fail with invalid column")

		var count int
		err = queryRow(`SELECT COUNT(*) FROM test_tx`).Scan(&count)
		assert.Nil(t, err)
		assert.Equal(t, 0, count, "rollback should keep table empty")
	})

	t.Run("multiple rows query", func(t *testing.T) {
		err := command(`CREATE TABLE IF NOT EXISTS test_multi(id INTEGER, value TEXT)`)
		assert.Nil(t, err)

		testData := []string{"one", "two", "three"}
		for i, val := range testData {
			err := command(`INSERT INTO test_multi(id, value) VALUES(?, ?)`, i+1, val)
			assert.Nil(t, err)
		}

		rows, err := query(`SELECT id, value FROM test_multi ORDER BY id`)
		assert.Nil(t, err)
		defer rows.Close()

		var results []string
		for rows.Next() {
			var id int
			var value string
			err := rows.Scan(&id, &value)
			assert.Nil(t, err)
			results = append(results, value)
		}

		assert.Equal(t, testData, results)
	})

	t.Run("unique constraint", func(t *testing.T) {
		err := command(`INSERT INTO users(username, password) VALUES(?, ?)`, "unique_test", "pass1")
		assert.Nil(t, err)

		err = command(`INSERT INTO users(username, password) VALUES(?, ?)`, "unique_test", "pass2")
		assert.NotNil(t, err, "should fail with duplicate username")
	})

	t.Run("null handling", func(t *testing.T) {
		err := command(`CREATE TABLE IF NOT EXISTS test_null(id INTEGER, optional_field TEXT)`)
		assert.Nil(t, err)

		err = command(`INSERT INTO test_null(id, optional_field) VALUES(?, ?)`, 1, nil)
		assert.Nil(t, err)

		var id int
		var optionalField *string
		err = queryRow(`SELECT id, optional_field FROM test_null WHERE id=1`).Scan(&id, &optionalField)
		assert.Nil(t, err)
		assert.Equal(t, 1, id)
		assert.Nil(t, optionalField)
	})

	t.Run("empty result handling", func(t *testing.T) {
		var name string
		err := queryRow(`SELECT username FROM users WHERE username=?`, "nonexistent").Scan(&name)
		assert.Equal(t, sql.ErrNoRows, err, "should return ErrNoRows for non-existent data")
	})

	t.Run("large data handling", func(t *testing.T) {
		err := command(`CREATE TABLE IF NOT EXISTS test_large(id INTEGER, large_text TEXT)`)
		assert.Nil(t, err)

		largeText := strings.Repeat("A", 1024*1024)

		err = command(`INSERT INTO test_large(id, large_text) VALUES(?, ?)`, 1, largeText)
		assert.Nil(t, err, "should handle large text")

		var result string
		err = queryRow(`SELECT large_text FROM test_large WHERE id=1`).Scan(&result)
		assert.Nil(t, err)
		assert.Equal(t, len(largeText), len(result), "should retrieve same length")
	})
}

func TestAuthRepository(t *testing.T) {
	repo, err := NewAuthRepository()
	assert.Nil(t, err)

	t.Run("save and get user", func(t *testing.T) {
		user := auth.NewUser("testuser", "hashedpass")

		err := repo.SaveUser(user)
		assert.Nil(t, err)

		retrieved, err := repo.GetUser("testuser")
		assert.Nil(t, err)
		assert.Equal(t, "testuser", retrieved.Name)
		assert.Equal(t, "hashedpass", retrieved.Password)
	})

	t.Run("check user exists", func(t *testing.T) {
		exists, err := repo.CheckUserExists("testuser")
		assert.Nil(t, err)
		assert.True(t, exists)

		exists, err = repo.CheckUserExists("nonexistent")
		assert.Nil(t, err)
		assert.False(t, exists)
	})

	t.Run("token operations", func(t *testing.T) {
		token := &auth.Token{
			String:    "test_token_123",
			IssuedAt:  1000,
			ExpiresAt: 2000,
		}

		err := repo.SaveToken(token)
		assert.Nil(t, err)

		err = repo.DeleteToken("test_token_123")
		assert.Nil(t, err)
	})
}
