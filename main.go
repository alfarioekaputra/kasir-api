package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
	"labkoding.my.id/kasir-api/database"
	"labkoding.my.id/kasir-api/handler"
	"labkoding.my.id/kasir-api/repositories"
	"labkoding.my.id/kasir-api/services"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func handleError(w http.ResponseWriter, err error) {
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Route("/categories", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			categoryHandler.GetAllCategory(w, r)
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			categoryHandler.CreateCategory(w, r)
		})
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			categoryHandler.GetCategoryByID(w, r)
		})
		r.Put("/update/{id}", func(w http.ResponseWriter, r *http.Request) {
			categoryHandler.UpdateCategory(w, r)
		})
		r.Delete("/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
			categoryHandler.DeleteCategory(w, r)
		})
	})

	r.Route("/products", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			productHandler.GetAllProduct(w, r)
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			productHandler.CreateProduct(w, r)
		})
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			productHandler.GetProductByID(w, r)
		})
		r.Put("/update/{id}", func(w http.ResponseWriter, r *http.Request) {
			productHandler.UpdateProduct(w, r)
		})
		r.Delete("/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
			productHandler.DeleteProduct(w, r)
		})
	})

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running di", addr)

	err = http.ListenAndServe(addr, r)
	if err != nil {
		fmt.Println("gagal running server", err)
	}
}
