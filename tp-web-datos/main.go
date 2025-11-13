package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	sqlc "tp_web_datos/db/sqlc"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "user=userdb password=admin dbname=tpwebdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()
	queries := sqlc.New(db)
	ctx := context.Background()

	// Create
	fecha, err := time.Parse(time.DateOnly, "2025-08-20")
	if err != nil {
		panic(err)
	}

	createdGame, err := queries.CreateGame(ctx,
		sqlc.CreateGameParams{
			Titulo:      "BT6",
			Descripcion: "Juego de disparos en primera persona.",
			Categoria:   "Accion",
			Fecha:       fecha,
			Estado:      "none",
			Imagen:      "/img/btf6.png",
		})

	if err != nil {
		log.Fatalf("failed to create game: %v", err)
	}
	fmt.Printf("Created game: %+v\n", createdGame)

	// Read one game
	game, err := queries.GetGame(ctx, createdGame.ID)
	if err != nil {
		log.Fatalf("failed to get game: %v", err)
	}
	fmt.Printf("Retrieved game: %+v\n", game)

	// List all games
	games, err := queries.ListGames(ctx)
	if err != nil {
		log.Fatalf("failed to list games: %v", err)
	}
	fmt.Printf("All games: %+v\n", games)

	// Update game
	_, err = queries.UpdateGame(ctx, sqlc.UpdateGameParams{
		ID:          createdGame.ID,
		Titulo:      "BT5",
		Descripcion: "Juego de disparos en primera persona.",
		Categoria:   "Accion",
		Fecha:       fecha,
		Estado:      "none",
		Imagen:      "/img/btf5.png",
	})

	if err != nil {
		log.Fatalf("failed to update game: %v", err)
	}

	fmt.Println("Game updated successfully")

	// Update game state
	_, err = queries.UpdateGameState(ctx, sqlc.UpdateGameStateParams{
		ID:     createdGame.ID,
		Estado: "deseado",
	})

	if err != nil {
		log.Fatalf("failed to update game: %v", err)
	}

	fmt.Println("Game updated successfully")

	// List wanted games
	wantedGames, err := queries.ListWantedGames(ctx)
	if err != nil {
		log.Fatalf("failed to get wanted game: %v", err)
	}

	fmt.Printf("Wanted games: %+v\n", wantedGames)

	// Delete game
	_, err = queries.DeleteGame(ctx, createdGame.ID)
	if err != nil {
		log.Fatalf("failed to delete game: %v", err)
	}

	fmt.Println("Game deleted successfully")

	_, err = queries.GetGame(ctx, createdGame.ID)
	if err == sql.ErrNoRows {
		fmt.Println("Game not found after deletion")
	} else if err != nil {
		log.Fatalf("failed to get game after deletion: %v", err)
	}
}
