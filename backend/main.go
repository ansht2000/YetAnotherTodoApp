package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/ansht2000/YetAnotherTodoApp/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	secretKey string
}

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	dbConn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("could not establish connection to database")
	}
	dbQueries := database.New(dbConn)
	secretKey := os.Getenv("SECRET")
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db: dbQueries,
		secretKey: secretKey,
	}
	
	serveMux := http.NewServeMux()
	server := http.Server{
		Addr: ":8080",
		Handler: serveMux,
	}

	serveMux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	serveMux.HandleFunc("POST /api/login", apiCfg.handlerLoginUser)

	server.ListenAndServe()
}