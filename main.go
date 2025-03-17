package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/balayher/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer dbConn.Close()
	dbQueries := database.New(dbConn)

	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		jwtSecret:      jwtSecret,
	}

	ServeMux := http.NewServeMux()
	ServeMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	ServeMux.HandleFunc("GET /api/healthz", handlerReadiness)

	ServeMux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	ServeMux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)

	ServeMux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	ServeMux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	ServeMux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	ServeMux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirps)
	ServeMux.HandleFunc("GET /api/chirps", apiCfg.handlerRetrieveChirps)
	ServeMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerRetrieveSingleChirp)
	ServeMux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteSingleChirp)

	ServeMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	ServeMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	server := &http.Server{
		Handler: ServeMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
