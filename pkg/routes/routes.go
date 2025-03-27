package api

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/vvjke314/vk-test-03-2025/internal/usecases"
	"github.com/vvjke314/vk-test-03-2025/pkg/handlers"
)

func SetupRoutes(uc *usecases.KeyValueUseCase) *chi.Mux {
	r := chi.NewRouter()

	r.Use(loggingMiddleware)

	r.Route("/kv", func(r chi.Router) {
		logger := log.New(os.Stdout, "KV_HANDLER: ", log.LstdFlags)
		handler := handlers.NewKVHandler(uc, logger)
		r.Post("/", handler.CreateKeyHandler)
		r.Put("/{id}", handler.UpdateKeyHandler)
		r.Get("/{id}", handler.GetKeyHandler)
		r.Delete("/{id}", handler.DeleteKeyHandler)
	})

	return r
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
