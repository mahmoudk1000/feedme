package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/mahmoudk1000/feedme/internal/database"
)

const collectionConcurrency = 10

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env was not found/loaded")
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Print("PORT is not set, using default port :8080")
	}

	intervalEnv := os.Getenv("FETCH_INTERVAL")
	if intervalEnv == "" {
		log.Print("FETCH_INTERVAL is not set, using default 60 seconds")
	}
	interval, err := time.ParseDuration(intervalEnv)
	if err != nil {
		log.Printf("FETCH_INTERVAL is not a valid duration: %v", err)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		DB: dbQueries,
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"link"},
		AllowCredentials: false,
		MaxAge:           300,
	}),
		middleware.Logger,
	)

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReady)
	v1Router.Get("/err", handleErr)
	v1Router.Post("/users", apiCfg.userCreateHandler)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.userGetHandler))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.createFeedHandler))
	v1Router.Get("/feeds", apiCfg.getAllFeedsHandler)
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.getUserFeedFollowes))
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.followFeedHanlder))
	v1Router.Delete(
		"/feed_follows/{feedFollowID}",
		apiCfg.middlewareAuth(apiCfg.deleteFeedFollowHandler),
	)
	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.getUserPostsHandler))

	router.Mount("/v1", v1Router)

	s := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go scrapeWorker(dbQueries, collectionConcurrency, interval)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(s.ListenAndServe())
}
