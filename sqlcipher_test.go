package sqlite3_test

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	sqlite3 "github.com/thinkgos/go-sqlcipher"
)

var (
	testDb  *sql.DB
	testDir = "go-sqlcipher_test"
	tables  = `
CREATE TABLE KeyValueStore (
  KeyEntry   TEXT NOT NULL UNIQUE,
  ValueEntry TEXT NOT NULL
);`
)

func init() {
	// create DB
	key := url.QueryEscape("passphrase")
	tmpdir, err := os.MkdirTemp("", testDir)
	if err != nil {
		panic(err)
	}
	dbname := filepath.Join(tmpdir, "sqlcipher_test")
	dbnameWithDSN := dbname + fmt.Sprintf("?_pragma_key=%s&_pragma_cipher_page_size=4096", key)
	testDb, err = sql.Open("sqlite3", dbnameWithDSN)
	if err != nil {
		panic(err)
	}
	_, err = testDb.Exec(tables)
	if err != nil {
		panic(err)
	}
	testDb.Close()
	// make sure DB is encrypted
	encrypted, err := sqlite3.IsEncrypted(dbname)
	if err != nil {
		panic(err)
	}
	if !encrypted {
		panic(errors.New("go-sqlcipher: DB not encrypted"))
	}
	// open DB for testing
	testDb, err = sql.Open("sqlite3", dbnameWithDSN)
	if err != nil {
		panic(err)
	}
	_, err = testDb.Exec("SELECT count(*) FROM sqlite_master;")
	if err != nil {
		panic(err)
	}
}

var mapping = map[string]string{
	"foo": "one",
	"bar": "two",
	"baz": "three",
}

func TestSQLCipherParallelInsert(t *testing.T) {
	t.Parallel()
	insertValueQuery, err := testDb.Prepare("INSERT INTO KeyValueStore (KeyEntry, ValueEntry) VALUES (?, ?);")
	noError(t, err)
	for key, value := range mapping {
		_, err := insertValueQuery.Exec(key, value)
		noError(t, err)
	}
}

func TestSQLCipherParallelSelect(t *testing.T) {
	t.Parallel()
	getValueQuery, err := testDb.Prepare("SELECT ValueEntry FROM KeyValueStore WHERE KeyEntry=?;")
	if err != nil {
		t.Fatal(err)
	}
	for key, value := range mapping {
		var val string
		err := getValueQuery.QueryRow(key).Scan(&val)
		if err != sql.ErrNoRows {
			if noError(t, err) {
				if value != val {
					t.Fatal("should not be equal")
				}
			}
		}
	}
}

func TestSQLCipherIsEncryptedFalse(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", testDir)
	noError(t, err)
	defer os.RemoveAll(tmpdir)
	dbname := filepath.Join(tmpdir, "unencrypted.sqlite")
	testDb, err := sql.Open("sqlite3", dbname)
	noError(t, err)
	defer testDb.Close()
	_, err = testDb.Exec(tables)
	noError(t, err)
	encrypted, err := sqlite3.IsEncrypted(dbname)
	if noError(t, err) {
		if encrypted {
			t.Fatal("should not be encrypted")
		}
	}
}

func TestSQLCipherIsEncryptedTrue(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", testDir)
	noError(t, err)
	defer os.RemoveAll(tmpdir)
	dbname := filepath.Join(tmpdir, "encrypted.sqlite")
	var key [32]byte
	_, err = io.ReadFull(rand.Reader, key[:])
	noError(t, err)
	dbnameWithDSN := dbname + fmt.Sprintf("?_pragma_key=x'%s'",
		hex.EncodeToString(key[:]))
	testDb, err := sql.Open("sqlite3", dbnameWithDSN)
	noError(t, err)
	defer testDb.Close()
	_, err = testDb.Exec(tables)
	noError(t, err)
	encrypted, err := sqlite3.IsEncrypted(dbname)
	if noError(t, err) {
		if !encrypted {
			t.Fatal("should be encrypted")
		}
	}
}

func TestSQLCipher3DB(t *testing.T) {
	dbname := filepath.Join("testdata", "sqlcipher3.sqlite3")
	dbnameWithDSN := dbname + "?_pragma_key=passphrase&_pragma_cipher_page_size=4096"
	// make sure DB is encrypted
	encrypted, err := sqlite3.IsEncrypted(dbname)
	if err != nil {
		t.Fatal(err)
	}
	if !encrypted {
		t.Fatal("go-sqlcipher: DB not encrypted")
	}
	// open DB for testing
	testDb, err := sql.Open("sqlite3", dbnameWithDSN)
	if err != nil {
		t.Fatal(err)
	}
	// should fail
	_, err = testDb.Exec("SELECT count(*) FROM sqlite_master;")
	if err == nil {
		t.Fatal(errors.New("opening a SQLCipher 3 database with SQLCipher 4 should fail"))
	}
}

func TestSQLCipher4DB(t *testing.T) {
	dbname := filepath.Join("testdata", "sqlcipher4.sqlite3")
	dbnameWithDSN := dbname + "?_pragma_key=passphrase&_pragma_cipher_page_size=4096"
	// make sure DB is encrypted
	encrypted, err := sqlite3.IsEncrypted(dbname)
	if err != nil {
		t.Fatal(err)
	}
	if !encrypted {
		t.Fatal("go-sqlcipher: DB not encrypted")
	}
	// open DB for testing
	testDb, err := sql.Open("sqlite3", dbnameWithDSN)
	if err != nil {
		t.Fatal(err)
	}
	// should succeed
	_, err = testDb.Exec("SELECT count(*) FROM sqlite_master;")
	if err != nil {
		t.Fatal(err)
	}
}

func ExampleIsEncrypted() {
	// create random key
	var key [32]byte
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		log.Fatal(err)
	}
	// set DB name
	dbname := "go-sqlcipher.sqlite"
	dbnameWithDSN := dbname + fmt.Sprintf("?_pragma_key=x'%s'",
		hex.EncodeToString(key[:]))
	// create encrypted DB file
	testDb, err := sql.Open("sqlite3", dbnameWithDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(dbname)
	defer testDb.Close()
	// create table
	_, err = testDb.Exec("CREATE TABLE t(x INTEGER);")
	if err != nil {
		log.Fatal(err)
	}
	// make sure database is encrypted
	encrypted, err := sqlite3.IsEncrypted(dbname)
	if err != nil {
		log.Fatal(err)
	}
	if encrypted {
		fmt.Println("DB is encrypted")
	} else {
		fmt.Println("DB is unencrypted")
	}
	// Output:
	// DB is encrypted
}

func noError(t *testing.T, err error) bool {
	if err != nil {
		t.Fatal(err)
		return false
	}
	return true
}
