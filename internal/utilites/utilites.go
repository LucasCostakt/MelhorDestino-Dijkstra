package utilites

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	d "github.com/MelhorDestino/internal/dijkstra"
)

//retira "-" da string separando os dois destinos o de partida e o de chegada
func SplitRoutes(scan string) []string {
	s := strings.Split(scan, "-")
	if len(s) < 2 {
		s = append(s, " ")
	}
	return s
}

//Retorna um array de string contendo todo o conteudo do csv selecionado
func OpenFile(nameFile string) []string {
	var fileContents []string

	file, err := os.Open("../../csv/" + nameFile)
	if err != nil {
		log.Println("os Open error in OpenFile() ", err)
		return []string{"os Open error"}
	}
	in := file
	buf := bufio.NewScanner(in)

	for buf.Scan() {
		fileContents = append(fileContents, buf.Text())
	}

	if err := buf.Err(); err != nil {
		log.Println("reading file error in OpenFile() ", err)
		return []string{" reading file error"}
	}

	defer file.Close()

	return fileContents
}

//Converte o conteudo do arquvido lido em um array de destinos sendo uma struct contendo
//o destino inicial, destino final e o custo
func ConvertFileInDestinations(fileContentes []string) []d.Destination {
	destinations := []d.Destination{}

	for _, dest := range fileContentes {
		s := strings.Split(dest, ",")
		i, _ := strconv.Atoi(s[2])
		dests := d.Destination{
			Start: s[0],
			Final: s[1],
			Cost:  i,
		}
		destinations = append(destinations, dests)
	}

	return destinations
}
