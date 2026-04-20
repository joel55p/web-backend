package main

import (
	"encoding/json" // Para convertir structs a JSON y viceversa
	"strings"
)

// Serie representa una fila de la table
type Serie struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	CurrentEpisode int    `json:"current_episode"`
	TotalEpisodes  int    `json:"total_episodes"`
	ImageURL       string `json:"image_url"`
}

// SerieInput es lo que llega en el body de  POST/PUT (va sin ID)
type SerieInput struct { //datos que el cliente envia 
	Name           string `json:"name"`
	CurrentEpisode int    `json:"current_episode"`
	TotalEpisodes  int    `json:"total_episodes"`
	ImageURL       string `json:"image_url"`
}

// toJSON convierte cualquier valor a JSON string que es lo important pal project
func toJSON(v any) string {
	b, err := json.Marshal(v) //Convierte el valor v a JSON. 
	if err != nil {
		return `{"error":"error serializando JSON"}`
	}
	return string(b) // Devuelve el JSON como string para poder incluirlo en la respuesta HTTP.
}

// parseSerieInput parsea el body JSON a SerieInput
func parseSerieInput(body string) (SerieInput, error) { //recibe el json del cliente como string y devuelve un struct para manejar datos en Go 
	var input SerieInput //para guardar el resultado del parseo
	err := json.Unmarshal([]byte(strings.TrimSpace(body)), &input) //Convierte el JSON del body a un struct SerieInput.
	return input, err // Devuelve el struct y cualquier error que haya ocurrido durante el parseo (ej: JSON mal formado)
}

// validarInput verifica que los datos  obligatorios esten ahi 
// Por lo que devuelve mensaje de error o "" si esta bien
func validarInput(input SerieInput) string {
	if strings.TrimSpace(input.Name) == "" {
		return `{"error":"el dato 'name' es obligatorio"}`
	}
	if input.TotalEpisodes <= 0 {
		return `{"error":"'total_episodes' debe ser > 0"}`
	}
	if input.CurrentEpisode < 0 {
		return `{"error":"'current_episode' no puede ser un valor negativo"}`
	}
	if input.CurrentEpisode > input.TotalEpisodes {
		return `{"error":"'current_episode' no puede ser mayor a 'total_episodes'"}`
	}
	return ""
}