package main

import (
	api "github.com/MelhorDestino/internal/apirest"
)

//A função é onde vai fazer somente as chamadas das funções para dar o start na api
//Inicia no endereço localhost:5000
func main() {
	http := api.CreateHttp()
	routes := api.NewRoutes(http)
	api.StartAPI(routes)
}
