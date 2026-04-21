package main

import (
	"database/sql" // Para interactuar con SQLite
	"log" // Para imprimir logs en la consola
)

// handleDelete maneja DELETE /series/:id es decir que va a elimar series por id
// Elimina la serie y devuelve 204 No Content (sin body ya que es el estandar REST)
func handleDelete(db *sql.DB, id string) string {

	// Ver que existe antes de intentar eliminar
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM series WHERE id = ?", id).Scan(&exists)
	if err != nil || exists == 0 { //si da error
		return jsonResponse(404, `{"error":"serie no encontrada"}`)
	}

	// eliminar serie  de la base de datos
	_, err = db.Exec("DELETE FROM series WHERE id = ?", id)
	if err != nil {
		log.Print("Error eliminando serie:", err) //si da error
		return jsonResponse(500, `{"error":"error al eliminar la serie"}`)
	}

	// 204 No Content: exito pero sin body (es estandar para DELETE)
	return jsonResponse(204, "")
}