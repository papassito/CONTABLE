package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/klik/fcos-kernel/config"
	"github.com/klik/fcos-kernel/pkg/database"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

type ApplicationEngine struct {
	DB   *sql.DB
	UOW  database.UnitOfWork
	Mode string
}

func BootstrapEngine(cfg *config.Config) (*ApplicationEngine, error) {
	var db *sql.DB
	var err error

	switch cfg.DatabaseDriver {
	case "postgres":
		db, err = sql.Open("pgx", cfg.DatabaseURL)
		if err != nil {
			return nil, fmt.Errorf("falla al conectar con PostgreSQL: %w", err)
		}
	case "sqlite":
		db, err = sql.Open("sqlite", cfg.DatabaseURL)
		if err != nil {
			return nil, fmt.Errorf("falla al conectar con SQLite: %w", err)
		}
		if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
			return nil, fmt.Errorf("falla al activar PRAGMA foreign_keys en SQLite: %w", err)
		}
	default:
		return nil, fmt.Errorf("driver de base de datos no soportado: %s", cfg.DatabaseDriver)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("falla en verificación de ping a BD: %w", err)
	}

	uow := database.NewUnitOfWork(db)

	return &ApplicationEngine{
		DB:   db,
		UOW:  uow,
		Mode: cfg.DatabaseDriver,
	}, nil
}

func main() {
	log.Println("Inicializando Fiscal Compliance Operating System (FCOS v2.2)...")

	cfg := config.LoadFromEnv()
	engine, err := BootstrapEngine(cfg)
	if err != nil {
		log.Fatalf("Error crítico durante el arranque del Kernel: %v", err)
		os.Exit(1)
	}
	defer engine.DB.Close()

	log.Printf("FCOS Engine operando satisfactoriamente en modo: %s", engine.Mode)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_ = ctx
}
