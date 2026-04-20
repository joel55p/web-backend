package main

import (
	"bufio" // Para leer la request HTTP línea por línea
	"database/sql" // Para interactuar con SQLite
	"log" // Para imprimir logs en la consola
	"net" // Para manejar conexiones TCP
	"strconv" // Para convertir strings a enteros
	"strings" // Para manipular strings (realmente para  parsear la request line y headers)
)

// handle lee la request HTTP cruda, la parsea, y la despacha al handler correcto
//osea recibe la http del cliente, la procesa y devuelve la respuesta adecuada segun la ruta y el método
func handle(conn net.Conn, db *sql.DB) { //va a recibir la conexión TCP y la conexión a la base de datos para poder hacer consultas
	defer conn.Close() //cierra al final para liberar recursos

	reader := bufio.NewReader(conn) //Va a leer la request línea por línea usando un buffer

	//  Leer request line (como "GET /series HTTP/1.1") 
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		log.Print("Error leyendo request:", err)
		return
	}

	parts := strings.Fields(requestLine) // Divide la request line en partes (metodo, ruta, version)
	if len(parts) < 2 { //si no tiene por lo menos 2 partes entonces es una request mal formada
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}

	method := parts[0]
	fullPath := parts[1]

	// Separar path de query string (ej: /series?page=1 -> path=/series, query=page=1)
	pathParts := strings.SplitN(fullPath, "?", 2)
	path := pathParts[0]

	//Leer headers para obtener el Content Length 
	contentLength := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil || line == "\r\n" {
			break
		}
		if strings.HasPrefix(line, "Content-Length:") { //Si el header es Content-Length, extraemos el valor para saber cuánto body leer después
			lengthStr := strings.TrimSpace(strings.TrimPrefix(line, "Content-Length:"))
			contentLength, _ = strconv.Atoi(lengthStr) //Convertimos el string a entero (si falla, contentLength queda en 0 y no se lee body)
		}
	}

	// LEER el  body si hay Content-Length
	body := ""
	if contentLength > 0 {
		bodyBytes := make([]byte, contentLength)
		reader.Read(bodyBytes)
		body = string(bodyBytes)
	}

	// CORS: responder OPTIONS preflight antes de enviar realmente la request 
	// El navegador manda OPTIONS antes de POST/PUT/DELETE cuando el origen es diferente para verificar que el servidor permite esa operación.
	if method == "OPTIONS" {
		conn.Write([]byte(corsHeaders() + "HTTP/1.1 204 No Content\r\n\r\n"))
		return
	}

	// Router: extraer ID si la ruta es /series/:id 
	// Ejemplos: /series -> id=""  |  /series/3 -> id="3" en este caso
	var response string // Aquí es donde se va a decidir qué handler llamar según el método y la ruta
	seriesID := ""

	segments := strings.Split(strings.Trim(path, "/"), "/")
	// segments[0] = "series", segments[1] = id (si es q existe)
	if len(segments) == 2 && segments[0] == "series" {
		seriesID = segments[1] //se le pone el id
	}

	switch { //aqui va a estar la logica para decidir que handler llamar segun el metodo y la ruta
	// GET /series para  listar todas
	case method == "GET" && path == "/series":
		response = handleTodos(db, fullPath)

	// GET /series/:id para obtener una
	case method == "GET" && segments[0] == "series" && seriesID != "":
		response = handleUno(db, seriesID)

	// POST /series para  crear nueva serie
	case method == "POST" && path == "/series":
		response = handleCreate(db, body)

	// PUT /series/:id  para editar serie
	case method == "PUT" && segments[0] == "series" && seriesID != "":
		response = handleUpdate(db, seriesID, body)

	// DELETE /series/:id para eliminar serie
	case method == "DELETE" && segments[0] == "series" && seriesID != "":
		response = handleDelete(db, seriesID)

	default:
		response = jsonResponse(404, `{"error":"ruta no encontrada"}`)
	}

	// Escribir respuesta con headers CORS incluidos
	conn.Write([]byte(corsHeaders() + response)) // Siempre incluimos los headers CORS para permitir que el frontend pueda consumir esta API desde otro origen

	//Manda la respuesta HTTP al cliente, incluyendo los headers CORS para que el navegador no bloquee la peticion.

}

// corsHeaders devuelve los headers necesarios para que el navegador permita fetch() desde otro origen(en README se explica mejor)
func corsHeaders() string {
	return "Access-Control-Allow-Origin: *\r\n" +
		"Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS\r\n" +
		"Access-Control-Allow-Headers: Content-Type\r\n"
}

// jsonResponse construye una respuesta HTTP con status code y body JSON
func jsonResponse(status int, body string) string {
	statusText := map[int]string{
		200: "OK",
		201: "Created",
		204: "No Content",
		400: "Bad Request",
		404: "Not Found",
		500: "Internal Server Error",
	}
	text := statusText[status] 
	if text == "" {
		text = "OK"
	}

	if status == 204 {
		// 204 no tiene body
		return "HTTP/1.1 204 No Content\r\nContent-Type: application/json\r\n\r\n"
	}

	return "HTTP/1.1 " + strconv.Itoa(status) + " " + text + "\r\n" + // Siempre se responde  con Content-Type application/json porque el frontend espera JSON(aun pa errores)
		"Content-Type: application/json\r\n\r\n" +
		body
}