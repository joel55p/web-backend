package main

import (
	"database/sql"
	"log"
	"net"

	_ "modernc.org/sqlite"
)

func main() {
	// Abrir conexión a SQLite (crea el archivo si no existe)
	db, err := sql.Open("sqlite", "series.db")
	if err != nil {
		log.Fatal("Error abriendo base de datos:", err)
	}
	defer db.Close()

	// Crear tabla si no existe (por si alguien corre el proyecto desde cero)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS series (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		current_episode INTEGER NOT NULL DEFAULT 1,
		total_episodes INTEGER NOT NULL,
		image_url TEXT,
	)`)
	if err != nil {
		log.Fatal("Error creando tabla:", err)
	}

	// Escuchar en el puerto 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error iniciando servidor:", err)
	}
	defer listener.Close()

	log.Print("Backend corriendo en http://localhost:8080")

	// Aceptar conexiones en bucle infinito
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print("Error aceptando conexion:", err)
			continue
		}
		// Cada conexión se maneja en una goroutine separada
		go handle(conn, db)
	}
}