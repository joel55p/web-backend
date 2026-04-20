package main

import (
	"database/sql" // Para interactuar con SQLite
	"log" // Para imprimir logs en la consola
	"net/url" // Para parsear query strings
	"strconv" // Para convertir strings a enteros
	"strings" // Para manipular strings (realmente para parsear la request line y headers)
)

// Este archivo (handleGetAll) lo que realmente va a hacer es manejar GET /series que es listar todas las series 
// Soporta: ?q=nombre  ?sort=name  ?order=asc|desc  ?page=1  ?limit=10
func handleTodos(db *sql.DB, fullPath string) string {

	// Parsear query string
	queryString := ""
	if idx := strings.Index(fullPath, "?"); idx != -1 { // Si hay un "?" en la ruta, entonces lo separamos para obtener solo el query string
		queryString = fullPath[idx+1:] // queryString va a ser todo lo que viene después del "?" en la ruta, que es donde están los params de  busqueda, ordenamiento y paginacion
	}
	params, _ := url.ParseQuery(queryString)

	search  := params.Get("q")  // q es el termino de busqueda para filtrar por nombre (ej: /series?q=got va a buscar series que tengan "got" en el nombre)
	sortBy  := params.Get("sort")  //campo para ordenar 
	order   := params.Get("order") // orden ascendente o descendente (15pts)
	pageStr := params.Get("page")  // numero de pagina para paginacion
	limitStr := params.Get("limit") // cantidad de series a mostrar por pagina para la paginacion

	// Valores default de paginacion (30pts)
	page := 1
	limit := 20
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 { // Si pageStr se puede convertir a entero y es mayor a 0, entonces lo usamos como page. Sino, page queda en 1 por defecto
		page = p
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 { 
		limit = l
	}
	offset := (page - 1) * limit // Calculamos el offset para la consulta SQL (ej: si page=2 y limit=20, entonces offset=20, lo que significa que la consulta va a saltar las primeras 20 series y mostrar las siguientes)

	// Validar sort para evitar SQL injection que es codigo SQL malicioso 
	allowedSort := map[string]bool{ //si se permite ordenar ya sea por nombre, episodio actual, total de episodios o id. Si el sortBy que viene en la query no es uno de estos, entonces se va a ordenar por id por defecto
		"name":            true, 
		"current_episode": true,
		"total_episodes":  true,
		"id":              true,
	}
	if !allowedSort[sortBy] {
		sortBy = "id"
	}
	if order != "asc" && order != "desc" {
		order = "asc" // por defecto
	}

	// Construir query
	query := "SELECT id, name, current_episode, total_episodes, image_url FROM series" //base
	args := []any{} //lista dinamica(slice) osea que no tiene tamaño establecido desde el inicio en este caso de cualquier tipo por el any{}

	if search != "" { // Si el usuario incluyo un term de busqueda, entonces se agrega un WHERE a la consulta para filtrar por nombre usando LIKE. El % antes y despues del termino de busqueda permite encontrar coincidencias que contengan el term en cualquier parte del nombre (antes, después o en medio)
		query += " WHERE name LIKE ?" //el ? es un placeholder que se va a reemplazar por el valor real del term de busqueda en la consulta SQL. Esto ayuda a prevenir SQL injection.
		args = append(args, "%"+search+"%") //args es un slice que va a contener los valores que se van a usar para reemplazar los placeholders (?) en la consulta SQL. En este caso, se agrega el term de busqueda con % antes y despues para el LIKE.
	}

	// Ordenar y paginar
	query += " ORDER BY " + sortBy + " " + order // Agregamos el ORDER BY a la consulta para ordenar por el campo y orden especificados por el usuario (o los valores por defecto si no se especificaron)
	query += " LIMIT ? OFFSET ?" //Agregamos LIMIT y OFFSET para la paginacion.
	args = append(args, limit, offset)  


	//------ EJECUTAR CONSULTA------
	rows, err := db.Query(query, args...) // Ejecutamos la consulta SQL con los argumentos correspondientes para reemplazar los placeholders.

	if err != nil { //si hay error
		log.Print("Error consultando series:", err)
		return jsonResponse(500, `{"error":"error interno del servidor"}`)
	}
	defer rows.Close() //cierre de resultado

	// Slice vacio para que JSON devuelva [] y no null si no hay datos de la series
	series := []Serie{} 
	for rows.Next() { //se itera sobre cada fila del resultado de la consulta
		var s Serie
		err := rows.Scan(&s.ID, &s.Name, &s.CurrentEpisode, &s.TotalEpisodes, &s.ImageURL) 
		if err != nil { //si hay error
			log.Print("Error escaneando fila:", err)
			continue
		}
		series = append(series, s) //se agrega la serie al slice de series que se va a devolver como resultado
	}

	return jsonResponse(200, toJSON(series)) // Se convierte el slice de series a JSON y se devuelve con un status 200 OK
}

