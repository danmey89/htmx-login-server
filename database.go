package htmxauthentication

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-yaml/yaml"
	_ "github.com/lib/pq"
)

type dbParams struct {
	DbName   string `yaml:"dbName"`
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SslMode  string `yaml:"sslmode"`
}

var config = dbParams{
	DbName: "postgres",
	Host: "localhost",
	User: "postgres",
	Password: "test123",
	SslMode: "disable",
}

func connectDB() (*sql.DB, error) {

	rr, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	var config dbParams

	if err := yaml.Unmarshal(rr, &config); err != nil {
		return nil, err
	}

	conn := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=%s",
		config.Host, config.DbName, config.User, config.Password, config.SslMode)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	return db, nil
}