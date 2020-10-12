package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	d "github.com/MelhorDestino/internal/dijkstra"
	conf "github.com/MelhorDestino/internal/utilites"
)

//A função é onde vai fazer as chamadas das funções para dar o start na aplicação e também pegar os valores digitados no teclado
func main() {
	var s string
	//flag vai pegar o nome do arquivo csv contendo as rotas
	flag.Parse()
	flagArg := flag.Arg(0)
	//abre o arquivo csv
	fileContentes := conf.OpenFile(flagArg)
	//converte as rotas do arquivo em um array de structs
	destinations := conf.ConvertFileInDestinations(fileContentes)
	graph := d.Graph{}

	for {
		fmt.Println(`Entre com a rota ou digite "exit" para sair`)
		//coleta o destino de partida e chegada
		_, err := fmt.Scan(&s)
		if err != nil {
			log.Fatal(err)
		}
		//converte tudo em em maiúscula para não ter problema de não achar a rota
		s := strings.ToUpper(s)
		if s == "EXIT" {
			break
		}
		routes := conf.SplitRoutes(s)
		cost, partida, sprintRoutes := graph.CreatBestDestination(routes[0], routes[1], destinations)
		fmt.Printf("Total de custo $%d, e as rotas são: %s -> %s\n", cost, partida, sprintRoutes)
	}

}
