package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	backend "github.com/lavatee/dipper_backend"
	"github.com/lavatee/dipper_backend/internal/endpoint"
	"github.com/lavatee/dipper_backend/internal/repository"
	"github.com/lavatee/dipper_backend/internal/service"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})
	if err := InitConfig(); err != nil {
		logrus.Fatalf("Config opening error: %s", err.Error())
	}
	db, err := repository.NewPostgresDB(repository.PostgresConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("DB opening error: %s", err.Error())
	}
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		logrus.Fatalf("Failed to create migrate driver: %s", err.Error())
	}

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	migrationsPath := "file://" + strings.ReplaceAll(filepath.Join(dir, "../schema"), "\\", "/")
	migrations, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		logrus.Fatalf("Failed to create migrate instance: %s", err.Error())
	}
	if err = migrations.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.Fatalf("Migrations error: %s", err.Error())
	}
	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	endp := endpoint.NewEndpoint(services, viper.GetString("botToken"))
	server := &backend.Server{}
	go func() {
		if err := server.Run(viper.GetString("port"), endp.InitRoutes()); err != nil {
			logrus.Fatalf("Server running error: %s", err.Error())
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	if err := server.Shutdown(); err != nil {
		logrus.Fatalf("Shutdown error: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Fatalf("DB closing error: %s", err.Error())
	}
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
