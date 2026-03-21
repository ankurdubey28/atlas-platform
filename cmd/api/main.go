package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file")
	}
	secret := os.Getenv("SECRET")
	fmt.Println(secret)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello World!"))
	})
	http.ListenAndServe(":3000", r)
}
