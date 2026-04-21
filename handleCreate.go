package main

import (
	"database/sql" // Para interactuar con SQLite
	"log" // Para imprimir logs en la consola
)

// handleCreate maneja POST /series por lo que se va a encargar de crear una serie nueva.
//Este archivo recibe JSON, luego inserta en DB, devuelve 201 con la serie creada
func handleCreate(db *sql.DB, body string) string { //Recibe la conexion a la base de datos para hacer la insercion y el body de la request para parsear el JSON

	input, err := parseSerieInput(body) //se parsea a SerieInput para validar y luego insertar en la base de datos
	if err != nil {
		return jsonResponse(400, `{"error":"JSON inválido o malformado"}`)
	}

	if errMsg := validarInput(input); errMsg != "" { //Si el input no es valido, se devuelve un error 400 con el mensaje 
		return jsonResponse(400, errMsg)
	}

	result, err := db.Exec( //query a ejecutar para insertar la nueva serie, y evitar las inyecciones maliciosas usando placeholders (?)
		`INSERT INTO series (name, current_episode, total_episodes, image_url)
		 VALUES (?, ?, ?, ?)`,
		input.Name, input.CurrentEpisode, input.TotalEpisodes, input.ImageURL,
	) //db.Exec se usa para ejecutar consultas que no devuelven filas (como INSERT, UPDATE, DELETE).
	//si el user no agrega imagen el valor de image_url va a ser "" lo cual esta bien porque en la tabla se permite NULL y "" es un string vacio que no rompe nada.
	if err != nil {
		log.Print("Error insertando serie:", err)
		return jsonResponse(500, `{"error":"error al guardar la serie"}`)
	}

	newID, _ := result.LastInsertId() //LastInsertId devuelve el ID del nuevo registro insertado, lo cual es servible  para devolver la serie creada con su ID en la respuesta.

	created := Serie{
		ID:             int(newID),
		Name:           input.Name,
		CurrentEpisode: input.CurrentEpisode,
		TotalEpisodes:  input.TotalEpisodes,
		ImageURL:       input.ImageURL,
	}

	// 201 Created es el código correcto al crear un recurso nuevo
	return jsonResponse(201, toJSON(created))
}