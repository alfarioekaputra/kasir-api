package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
	"labkoding.my.id/kasir-api/database"
	"labkoding.my.id/kasir-api/external"
	"labkoding.my.id/kasir-api/router"
)

type Config struct {
	Port            string `mapstructure:"PORT"`
	DBConn          string `mapstructure:"DB_CONN"`
	BucketName      string `mapstructure:"BUCKET_NAME"`
	AccountID       string `mapstructure:"ACCOUNT_ID"`
	AccessKeyID     string `mapstructure:"ACCESS_KEY_ID"`
	SecretAccessKey string `mapstructure:"SECRET_ACCESS_KEY"`
	PublicEndpoint  string `mapstructure:"PUBLIC_ENDPOINT"`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:            viper.GetString("PORT"),
		DBConn:          viper.GetString("DB_CONN"),
		BucketName:      viper.GetString("BUCKET_NAME"),
		AccountID:       viper.GetString("ACCOUNT_ID"),
		AccessKeyID:     viper.GetString("ACCESS_KEY_ID"),
		SecretAccessKey: viper.GetString("SECRET_ACCESS_KEY"),
		PublicEndpoint:  viper.GetString("PUBLIC_ENDPOINT"),
	}
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	err = external.InitStorage(config.BucketName, config.AccessKeyID, config.SecretAccessKey, config.AccountID, config.PublicEndpoint)
	if err != nil {
		log.Fatal("Failed to initialize storage:", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	appRouter := router.NewRouter(db, r)
	appRouter.RegisterAllRoutes()

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running di", addr)

	err = http.ListenAndServe(addr, r)
	if err != nil {
		fmt.Println("gagal running server", err)
	}
}
