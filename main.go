package main

import (
	"context"
	"log"
	"net/http"
	"os"
	db_connect "tp-web/db"
	datos "tp-web/db/sqlc"
)

var queries *datos.Queries
var ctx context.Context

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//Leo el html
		data, err := os.ReadFile("index.html")
		if err != nil {
			http.Error(w, "Error al cargar la pagina", http.StatusInternalServerError)
			return
		}

		games, err := queries.ListGames(ctx)
		if err != nil {
			log.Printf("Error en la capa de datos al listar todos los juegos: %v", err)
			http.Error(w, "Error inesperado", http.StatusInternalServerError)
			return
		}
		log.Printf("Juegos recuperados: %v", len(games))

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)

	})
	http.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		//Leo el css
		data, err := os.ReadFile("styles.css")
		if err != nil {
			http.Error(w, "Error al cargar los estilos", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/css")
		w.Write(data)

	})
	// Añadir handler para app.js (y otros assets si los hay)
	http.HandleFunc("/app.js", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("app.js")
		if err != nil {
			http.Error(w, "Error al cargar el script", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Write(data)
	})

	// Conexión a la base de datos
	db, err := db_connect.InitDb()
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	log.Println("Server connected to database successfully.")
	defer db.Close()

	queries = datos.New(db)
	ctx = context.Background()

	// Iniciar servidor en el puerto 8080
	log.Println("Presentación servida en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
