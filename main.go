package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"log/slog"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/spf13/viper"
	"labkoding.my.id/kasir-api/database"
	"labkoding.my.id/kasir-api/external"
	"labkoding.my.id/kasir-api/router"
)

type Config struct {
	Port            string `mapstructure:"PORT"`
	DBType          string `mapstructure:"DB_TYPE"`
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
		DBType:          viper.GetString("DB_TYPE"),
		DBConn:          viper.GetString("DB_CONN"),
		BucketName:      viper.GetString("BUCKET_NAME"),
		AccountID:       viper.GetString("ACCOUNT_ID"),
		AccessKeyID:     viper.GetString("ACCESS_KEY_ID"),
		SecretAccessKey: viper.GetString("SECRET_ACCESS_KEY"),
		PublicEndpoint:  viper.GetString("PUBLIC_ENDPOINT"),
	}

	// Set default DB type to postgres if not specified
	if config.DBType == "" {
		config.DBType = "postgres"
	}

	// Initialize database with type
	dbConfig := database.Config{
		Type:             database.DatabaseType(config.DBType),
		ConnectionString: config.DBConn,
		MaxOpenConns:     25,
		MaxIdleConns:     5,
	}

	db, err := database.InitDB(dbConfig)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Get underlying sql.DB for proper cleanup
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get sql.DB:", err)
	}
	defer sqlDB.Close()

	err = external.InitStorage(config.BucketName, config.AccessKeyID, config.SecretAccessKey, config.AccountID, config.PublicEndpoint)
	if err != nil {
		log.Fatal("Failed to initialize storage:", err)
	}

	// setup slog to write JSON logs to a daily-rotated file
	rl, err := rotatelogs.New(
		"app.%Y-%m-%d.log",
		rotatelogs.WithLinkName("app.log"),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithMaxAge(7*24*time.Hour),
	)
	if err != nil {
		log.Fatal("Failed to create rotatelogs:", err)
	}
	defer rl.Close()

	logger := slog.New(slog.NewJSONHandler(rl, &slog.HandlerOptions{}))
	slog.SetDefault(logger)

	r := chi.NewRouter()
	// use recoverer and custom request logger based on slog
	r.Use(middleware.Recoverer)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			slog.Info("http_request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", ww.Status()),
				slog.String("remote", r.RemoteAddr),
				slog.Duration("duration", time.Since(start)),
			)
		})
	})
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

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
