package main

import (
	"context"
	"database/sql"
	json "encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	datos "tp_web_datos/db/sqlc"
	models "tp_web_logica/api-models"

	_ "github.com/lib/pq"
)

var queries *datos.Queries
var ctx context.Context

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {

	// Configurar rutas
	http.HandleFunc("/games", enableCORS(gamesHandler))
	http.HandleFunc("/wanted_games", enableCORS(wantedHandler))
	http.HandleFunc("/games/", enableCORS(gameHandler))
	http.HandleFunc("/game_state/", enableCORS(stateHandler))

	// Leer variables de entorno
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	// Conexión a la base de datos
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	log.Println("Server connected to database successfully.")
	defer db.Close()
	queries = datos.New(db)
	ctx = context.Background()

	// Iniciar servidor
	log.Println("Server starting on :8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

// Manejador para /games
func gamesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getGames(w, r)
	case http.MethodPost:
		createGame(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Manejador para /wanted_games
func wantedHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getWantedGames(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Manejador para /games/{id}
func gameHandler(w http.ResponseWriter, r *http.Request) {

	// Extraer ID del path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		getGame(w, r, id)
	case http.MethodPut:
		updateGame(w, r, id)
	case http.MethodDelete:
		deleteProduct(w, r, id)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Manejador para /game_state/{id}?state="state"
func stateHandler(w http.ResponseWriter, r *http.Request) {

	// Extraer ID del path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	// 2. Leer el query param "state"
	state := r.URL.Query().Get("state")

	switch r.Method {
	case http.MethodPut:
		updateGameState(w, r, id, state)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /games - Listar todos los juegos
func getGames(w http.ResponseWriter, r *http.Request) {

	games, err := queries.ListGames(ctx)
	if err != nil {
		log.Printf("Error en la capa de datos al listar todos los juegos: %v", err)
		http.Error(w, "Error inesperado", http.StatusInternalServerError)
		return
	}
	log.Printf("Juegos recuperados: %v", len(games))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}

// GET /wanted_games - Listar los juegos deseados
func getWantedGames(w http.ResponseWriter, r *http.Request) {

	games, err := queries.ListWantedGames(ctx)
	if err != nil {
		log.Printf("Error en la capa de datos al listar los juegos deseados: %v", err)
		http.Error(w, "Error inesperado", http.StatusInternalServerError)
		return
	}
	log.Printf("Juegos deseados recuperados: %v", len(games))

	w.Header().Set("Content-Type", "application/json")

	if games == nil {
		games = []datos.ListWantedGamesRow{}
	}
	json.NewEncoder(w).Encode(games)
}

// POST /games - Crear nuevo juego
func createGame(w http.ResponseWriter, r *http.Request) {

	var req models.CreateGameReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Error al intentar decodificar el cuerpo: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convertir fecha (string -> time.Time)
	fechaParsed, err := time.Parse("2006-01-02", req.Fecha)
	if err != nil {
		http.Error(w, "Formato de fecha inválido (usar YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	// Validación de campos vacíos
	if req.Titulo == "" || req.Descripcion == "" || req.Categoria == "" ||
		req.Fecha == "" || req.Estado == "" || req.Imagen == "" {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	// Validar formato de la imagen
	imagenRegex := regexp.MustCompile(`^img/[a-zA-Z0-9_-]+\.(png|jpg|jpeg)$`)
	if !imagenRegex.MatchString(req.Imagen) {
		http.Error(w, "El campo imagen debe tener formato img/nombre.png|jpg|jpeg", http.StatusBadRequest)
		return
	}

	// Crear struct para sqlc
	gameParams := datos.CreateGameParams{
		Titulo:      req.Titulo,
		Descripcion: req.Descripcion,
		Categoria:   req.Categoria,
		Fecha:       fechaParsed,
		Estado:      req.Estado,
		Imagen:      req.Imagen,
	}

	createdGame, err := queries.CreateGame(ctx, gameParams)
	if err != nil {
		log.Printf("Error en la capa de datos al crear el juego: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Juego creado: %+v\n", createdGame)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(createdGame)
}

// GET /games/{id} - Obtener juego específico
func getGame(w http.ResponseWriter, r *http.Request, id int) {
	game, err := queries.GetGame(ctx, int32(id))
	if err != nil {
		log.Printf("Error en la capa de datos al obtener juego con id: {%v}: %v", id, err)
		http.Error(w, "Game Not Found", http.StatusNotFound)
		return
	}
	log.Printf("Juego especifico recuperado: %+v\n", game)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}

// PUT /games/{id} - Actualizar juego
func updateGame(w http.ResponseWriter, r *http.Request, id int) {
	var req models.UpdateGameReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Error al intentar decodificar el cuerpo: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convertir fecha (string -> time.Time)
	fechaParsed, err := time.Parse("2006-01-02", req.Fecha)
	if err != nil {
		http.Error(w, "Formato de fecha inválido (usar YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	// Validación de campos vacíos
	if req.Titulo == "" || req.Descripcion == "" || req.Categoria == "" ||
		req.Fecha == "" || req.Estado == "" || req.Imagen == "" {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	// Validar formato de la imagen
	imagenRegex := regexp.MustCompile(`^img/[a-zA-Z0-9_-]+\.(png|jpg|jpeg)$`)
	if !imagenRegex.MatchString(req.Imagen) {
		http.Error(w, "El campo imagen debe tener formato img/nombre.png|jpg|jpeg", http.StatusBadRequest)
		return
	}

	// Crear struct para sqlc
	updateGameParams := datos.UpdateGameParams{
		ID:          int32(id),
		Titulo:      req.Titulo,
		Descripcion: req.Descripcion,
		Categoria:   req.Categoria,
		Fecha:       fechaParsed,
		Estado:      req.Estado,
		Imagen:      req.Imagen,
	}

	updatedGame, err := queries.UpdateGame(ctx, updateGameParams)

	if err != nil {
		log.Printf("Error en la capa de datos al intentar actualizar el juego con id: {%v}: %v", id, err)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Game Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error inesperado", http.StatusInternalServerError)
		return
	}
	log.Printf("Juego actualizado correctamente: %+v\n", updatedGame)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedGame)
}

// PUT /games/{id}?state="state" - Actualizar estado juego
func updateGameState(w http.ResponseWriter, r *http.Request, id int, state string) {
	var updateGameStateParams datos.UpdateGameStateParams
	updateGameStateParams.ID = int32(id)
	updateGameStateParams.Estado = state

	updatedGame, err := queries.UpdateGameState(ctx, updateGameStateParams)

	if err != nil {
		log.Printf("Error en la capa de datos al intentar actualizar el juego con id: {%v}: %v", id, err)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Game Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error inesperado", http.StatusInternalServerError)
		return
	}
	log.Printf("Estado de juego actualizado correctamente: %+v\n", updatedGame)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedGame)
}

// DELETE /games/{id} - Eliminar juego
func deleteProduct(w http.ResponseWriter, r *http.Request, id int) {
	gameDeleted, err := queries.DeleteGame(ctx, int32(id))

	if err != nil {
		log.Printf("Error en la capa de datos al eliminar el juego con id: {%v}: %v", id, err)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Game Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error inesperado", http.StatusInternalServerError)
		return
	}
	log.Printf("Juego eliminado correctamente: %+v\n", gameDeleted)
	w.WriteHeader(http.StatusNoContent)
}
