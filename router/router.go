package router

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"labkoding.my.id/kasir-api/handler"
	"labkoding.my.id/kasir-api/repositories"
	"labkoding.my.id/kasir-api/services"
)

type Router struct {
	db     *sql.DB
	router chi.Router
}

func NewRouter(db *sql.DB, r chi.Router) *Router {
	return &Router{db: db, router: r}
}

func (rt *Router) RegisterCategoryRoutes() {
	categoryRepo := repositories.NewCategoryRepository(rt.db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	rt.router.Route("/categories", func(r chi.Router) {
		r.Get("/", categoryHandler.GetAllCategory)
		r.Post("/", categoryHandler.CreateCategory)
		r.Get("/{id}", categoryHandler.GetCategoryByID)
		r.Put("/{id}", categoryHandler.UpdateCategory)
		r.Delete("/{id}", categoryHandler.DeleteCategory)
	})
}

func (rt *Router) RegisterProductRoutes() {
	productRepo := repositories.NewProductRepository(rt.db)
	productService := services.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	rt.router.Route("/products", func(r chi.Router) {
		r.Get("/", productHandler.GetAllProduct)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/{id}", productHandler.GetProductByID)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})
}

func (rt *Router) RegisterTransactionRoutes() {
	transactionRepo := repositories.NewTransactionRepository(rt.db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	rt.router.Route("/transactions", func(r chi.Router) {
		r.Post("/checkout", transactionHandler.Checkout)
	})
}

func (rt *Router) RegisterReportRoutes() {
	reportRepo := repositories.NewReportRepository(rt.db)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handler.NewReportHandler(reportService)

	rt.router.Route("/report", func(r chi.Router) {
		r.Get("/", reportHandler.RangeReport)
		r.Get("/today", reportHandler.TodayReport)
	})
}

func (rt *Router) RegisterAllRoutes() {
	rt.RegisterCategoryRoutes()
	rt.RegisterProductRoutes()
	rt.RegisterTransactionRoutes()
	rt.RegisterReportRoutes()
}
