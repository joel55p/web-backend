package main

import (
	"database/sql"
	"log"
)

// handleUpdate maneja PUT /series/:id
//por lo que se va encagar de editar una serie existente. Recibe el ID de la serie a editar, el body con los nuevos datos, luego actualiza en DB y devuelve la serie actualizada.
func handleUpdate(db *sql.DB, id string, body string) string {

	// -----Verificar que la serie existe-----

	var exists int //es int porque solo es para ver si existe o no
	//queryrow solo devuelve una fila y como se usa el COUNT entonces seria o 0 o 1 dependiendo si la serie existe o no, y ese valor se guarda en la variable exists usando el scan.
	err := db.QueryRow("SELECT COUNT(*) FROM series WHERE id = ?", id).Scan(&exists) //ese scan toma el resultado de la consulta SQL y lo guarda en la variable exists. Si la consulta devuelve 0, entonces exist =0, lo que indica que no hay ninguna serie con ese ID. Si devuelve 1, entonces exists =1, lo que indica que si existe una serie con ese ID.
	if err != nil || exists == 0 { //si hay error
		return jsonResponse(404, `{"error":"serie no encontrada"}`)
	}

	// ---se parsesa el body para comprobar datos que vienen del cliente-----
	input, err := parseSerieInput(body) 
	if err != nil {
		return jsonResponse(400, `{"error":"JSON inválido o malformado"}`)
	}

	if errMsg := validarInput(input); errMsg != "" {
		return jsonResponse(400, errMsg)
	}

	_, err = db.Exec( //query para actualizar la serie con los nuevos datos, usando placeholders para evitar inyecciones maliciosas. El ID se incluye al final para indicar que serie se va a actualizar.
		`UPDATE series SET name=?, current_episode=?, total_episodes=?, image_url=?
		 WHERE id=?`,
		input.Name, input.CurrentEpisode, input.TotalEpisodes, input.ImageURL, id,
	)
	if err != nil { //si hay error
		log.Print("Error actualizando serie:", err)
		return jsonResponse(500, `{"error":"error al actualizar la serie"}`) 
	}

	//solo para devolverlos ya actualizados porque en la base de datos ya se actualizo.
	updated := Serie{  //se crea un struct Serie con los datos actualizados para devolverlo en la respuesta. El ID se mantiene igual porque solo se estan actualizando los otros campos.
		Name:           input.Name,
		CurrentEpisode: input.CurrentEpisode,
		TotalEpisodes:  input.TotalEpisodes,
		ImageURL:       input.ImageURL,
	}

	return jsonResponse(200, toJSON(updated))
}