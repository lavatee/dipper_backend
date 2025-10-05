package main

import (
	"os"
	"os/signal"
	"syscall"

	backend "github.com/lavatee/dipper_backend"
	"github.com/lavatee/dipper_backend/internal/endpoint"
	"github.com/lavatee/dipper_backend/internal/repository"
	"github.com/lavatee/dipper_backend/internal/service"
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
