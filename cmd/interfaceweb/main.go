package main

import (
	web "github.com/MelhorDestino/internal/web"
)

//A função maiin é onde vai fazer somente as chamadas das funções para dar o start na interface web
//Inicia no endereço localhost:5050
func main() {
	http := web.CreateHttp()
	routes := web.NewRoutes(http)
	web.CreateTemplates()
	web.StartAPI(routes)
}
