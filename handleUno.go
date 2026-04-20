package main

import (
	"database/sql" // Para interactuar con SQLite
	"log" // Para imprimir logs en la consola
)


// handleGetOne va a manejar  GET /series/:id 
//Basicamente obtener una serie por ID
func handleUno(db *sql.DB, id string) string { //recibe la conexion a la base de datos y el ID de la serie que se quiere obtener (como string porque viene de la ruta)
	var s Serie
	err := db.QueryRow( //queryRow se usa para consultas que devuelven una sola fila. En este caso, se selecciono la serie que tiene el ID que se quiere.
		"SELECT id, name, current_episode, total_episodes, image_url FROM series WHERE id = ?",
		id,
	).Scan(&s.ID, &s.Name, &s.CurrentEpisode, &s.TotalEpisodes, &s.ImageURL) //el & se usa para pasar la dirección de memoria de los campos del struct s, para que Scan pueda escribir los valores obtenidos de la consulta SQL en esos campos.

	if err == sql.ErrNoRows {
		return jsonResponse(404, `{"error":"serie no encontrada"}`)
	}
	if err != nil {
		log.Print("Error obteniendo serie:", err)
		return jsonResponse(500, `{"error":"error interno del servidor"}`)
	}

	return jsonResponse(200, toJSON(s))
}