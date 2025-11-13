document.addEventListener("DOMContentLoaded", () => {
  const apiUrl = "http://localhost:8081";

  document.getElementById("gamesList").addEventListener("click", async (event) => {
    // Verificar si el click fue en un botón de borrado
    if (event.target.matches("button.delete-btn")) {
      const gameId = event.target.dataset.gameId;
      try {
        const response = await fetch(`${apiUrl}/games/${gameId}`, {
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json'
          }
        });

        if (!response.ok) {
          throw new Error(`Failed to delete game: ${response.statusText}`);
        }

        // Actualizo la lista después del borrado exitoso
        document.getElementById("fetchGames").click();
        alert("Game deleted successfully");
      } catch (error) {
        console.error('Error deleting game:', error);
        alert(`Error deleting game: ${error.message}`);
      }
    }
  });

  document.getElementById("fetchGames").addEventListener("click", async () => {
    try {
      const response = await fetch(`${apiUrl}/games`);
      if (!response.ok) {
        throw new Error("Failed to fetch games");
      }
      const games = await response.json();
      const gamesList = document.getElementById("gamesList");
      gamesList.innerHTML = ""; 
      games.forEach(game => {
        const li = document.createElement("li");
        const div = document.createElement("div");
        const titulo = document.createElement("h2");
        titulo.textContent = game.titulo;
        const descripcion = document.createElement("p");
        descripcion.style.textAlign = "center";    
        descripcion.textContent = game.descripcion;

        const categoria = document.createElement("p");
        categoria.textContent = `Genero: ${game.categoria}`;

        const fecha = document.createElement("p");
        fecha.textContent = `Fecha salida: ${game.fecha}`;

        const estado = document.createElement("p");
        estado.textContent = `Estado:  ${game.estado}`;

        const imagen = document.createElement("img");
        imagen.src = game.imagen;
        imagen.alt = `Imagen de ${game.titulo}`;

        const eliminarBtn = document.createElement("button");
        eliminarBtn.textContent = "X";
        eliminarBtn.className = "delete-btn"; 
        eliminarBtn.dataset.gameId = game.id; 

        div.appendChild(titulo);
        div.appendChild(descripcion);
        div.appendChild(categoria);
        div.appendChild(fecha);
        div.appendChild(estado);
        div.appendChild(imagen);
        div.appendChild(eliminarBtn);
        div.className = "gameContainer"
        li.appendChild(div);
        li.style.listStyleType = "none";
        li.style.marginBottom = "20px";
        gamesList.appendChild(li);
      });
    } catch (error) {
      console.error("Error fetching games:", error);
    }
  });

  // Create a partir del formulario
  document.getElementById("createGameForm").addEventListener("submit", async (event) => {
    event.preventDefault();

    const imageValue = document.getElementById("gameImage").value;
    if (!imageValue.match(/^img\/[^/]+\.(png|jpg|jpeg)$/)) {
      alert("El campo imagen debe tener formato img/nombre.png|jpg|jpeg");
      return;
    }

    //Los estados validos son: "deseado", "jugando", "completado", "none" 
    const estadoValue = document.getElementById("gameState").value.toLowerCase();
    if (!["deseado", "jugando", "completado", "none"].includes(estadoValue)) {
      alert("Estado debe ser uno de: deseado, jugando, completado, none");
      return;
    }

    const newGame = {
      titulo: document.getElementById("gameTitle").value,
      descripcion: document.getElementById("gameDescription").value,
      categoria: document.getElementById("gameCategory").value,
      fecha: document.getElementById("gameDate").value,
      estado: estadoValue,
      imagen: imageValue,
    };

    try {
      const response = await fetch(`${apiUrl}/games`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(newGame),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText);
      }

      const createdGame = await response.json();
      alert(`Game created: ${createdGame.titulo}`);
      document.getElementById("createGameForm").reset();
      document.getElementById("fetchGames").click();
    } catch (error) {
      console.error("Error creating game:", error);
      alert(error.message);
    }
  });
});

