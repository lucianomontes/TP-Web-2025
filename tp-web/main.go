package main

import (
	"log"
	"net/http"
	"os"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//Leo el html
		data, err := os.ReadFile("index.html")
		if err != nil {
			http.Error(w, "Error al cargar la pagina", http.StatusInternalServerError)
			return
		}
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

	// Iniciar servidor en el puerto 8080
	log.Println("Presentación servida en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
