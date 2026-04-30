package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
)

var (
	tlsOnce     sync.Once
	tlsRegistered bool
)

func Connect(databaseURL string) (*sql.DB, error) {
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is empty")
	}

	// If the DSN contains tls=true or tls=tidb, register a custom TLS config
	if strings.Contains(databaseURL, "tls=true") || strings.Contains(databaseURL, "tls=tidb") {
		if err := registerTiDBTLS(); err != nil {
			return nil, fmt.Errorf("register TLS config: %w", err)
		}
		// Replace tls=true with tls=tidb (our registered name)
		databaseURL = strings.Replace(databaseURL, "tls=true", "tls=tidb", 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// registerTiDBTLS registers a TLS config for TiDB Cloud connections.
// It uses the system cert pool, which includes ISRG Root X1 (Let's Encrypt).
func registerTiDBTLS() error {
	var regErr error
	tlsOnce.Do(func() {
		rootCertPool := x509.NewCertPool()

		// Try to load custom CA from TLS_CA_PATH env var, otherwise use system certs
		if caPath := os.Getenv("TLS_CA_PATH"); caPath != "" {
			pem, err := os.ReadFile(caPath)
			if err != nil {
				regErr = fmt.Errorf("read CA file %s: %w", caPath, err)
				return
			}
			if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
				regErr = errors.New("failed to append CA PEM")
				return
			}
		} else {
			// Use system cert pool (includes ISRG Root X1 on most systems)
			systemPool, err := x509.SystemCertPool()
			if err != nil {
				regErr = fmt.Errorf("load system cert pool: %w", err)
				return
			}
			rootCertPool = systemPool
		}

		err := mysql.RegisterTLSConfig("tidb", &tls.Config{
			RootCAs: rootCertPool,
		})
		if err != nil {
			regErr = err
			return
		}
		tlsRegistered = true
	})
	return regErr
}
