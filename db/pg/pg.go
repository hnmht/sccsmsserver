package pg

import (
	"database/sql"
	"fmt"
	"sccsmsserver/setting"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var db *sql.DB

// Initialize the database connection
func Init(cfg *setting.PqConfig) (err error) {
	// Step 1: Generate a database connecting string
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.DbName,
	)
	// Step 2: Connect to PostgreSQL Database
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		zap.L().Error("postgresql database Init sql.Open failed:", zap.Error(err))
		return
	}

	// Step 3: Test the database connection
	err = db.Ping()
	if err != nil {
		zap.L().Error("postgresql database Init  db.Ping failed:", zap.Error(err))
		return
	}
	// Step 4: Set the maximum number of open connections to the database.
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	// Step 5: Set the maximum number of connections in the idle connection pool.
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	// Step 6: Check if the database initialization is complete
	ok, err := checkDbInit()
	if err != nil {
		zap.L().Error("postgresql database Init checkDbInit failed:", zap.Error(err))
		return
	}
	// if database initialization isn't complete, perform database initialization
	if !ok {
		_, err = createTable()
		if err != nil {
			zap.L().Error("postgresql database Init createTable failed:", zap.Error(err))
			return
		}
	}

	// Setp 7: Initialize RSA
	_, err = initRsa()
	if err != nil {
		zap.L().Error("postgresql database Init initRsa failed:", zap.Error(err))
		return
	}

	// Step 8: Initialize Current Server public information
	err = ServerPubInfo.Init()
	if err != nil {
		zap.L().Error("postgresql database Init ServerPubInfo.Init failed:", zap.Error(err))
		return
	}

	// Step 9: Initialize Current Server locale list
	err = initSysLocalList()
	if err != nil {
		zap.L().Error("postgresql database Init initSysLocalList failed:", zap.Error(err))
		return
	}

	// Step 12: Upgrade database schema version
	_, err = upgradeDb()
	if err != nil {
		zap.L().Error("postgresql database Init upgradeDb failed", zap.Error(err))
		return
	}

	zap.L().Info("Database connection initialized successfully.")
	return
}

// Close the database connection
func Close() {
	_ = db.Close()
}
