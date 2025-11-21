package main

import (
	"context"
	"net/http"
	"time"

	"log"
	db_connect "tp-web/db"
	datos "tp-web/db/sqlc"
	views "tp-web/views"

	"github.com/a-h/templ"
)

var queries *datos.Queries
var ctx context.Context

func main() {

	// Conexión a la base de datos
	db, err := db_connect.InitDb()
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	log.Println("Server connected to database successfully.")
	defer db.Close()

	queries = datos.New(db)
	ctx = context.Background()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		games, err := queries.ListGames(ctx)
		if err != nil {
			log.Printf("Error en la capa de datos al listar todos los juegos: %v", err)
			http.Error(w, "Error inesperado", http.StatusInternalServerError)
			return
		}

		log.Printf("Juegos recuperados: %v", len(games))

		templ.Handler(views.IndexPage("Lista de Juegos", games)).ServeHTTP(w, r)

	})

	// Handler POST para crear un nuevo juego
	http.HandleFunc("/games", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			log.Printf("parse form error: %v", err)
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		description := r.FormValue("description")
		category := r.FormValue("category")
		state := r.FormValue("state")
		image := r.FormValue("image")
		releaseStr := r.FormValue("release_date")

		var releaseDate time.Time
		if releaseStr != "" {
			d, err := time.Parse("2006-01-02", releaseStr)
			if err != nil {
				log.Printf("invalid date: %v", err)
				http.Error(w, "invalid date", http.StatusBadRequest)
				return
			}
			releaseDate = d
		}

		_, err := queries.CreateGame(r.Context(), datos.CreateGameParams{
			Titulo:      title,
			Descripcion: description,
			Categoria:   category,
			Fecha:       releaseDate,
			Estado:      state,
			Imagen:      image,
		})
		if err != nil {
			log.Printf("create game error: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Iniciar servidor en el puerto 8080
	log.Println("Presentación servida en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
