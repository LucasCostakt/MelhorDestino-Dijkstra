package apirest

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	d "github.com/MelhorDestino/internal/dijkstra"
	conf "github.com/MelhorDestino/internal/utilites"
)

type httpServer struct {
	http.Handler
}

func CreateHttp() *httpServer {
	log.Println("Create new httpServer")
	return new(httpServer)
}

//Cria as novas rotas
func NewRoutes(h *httpServer) *httpServer {
	log.Println("Init Routes")
	router := http.NewServeMux()
	//criados os endpoint "/" e "/consult"
	router.Handle("/", http.HandlerFunc(csvOperator))
	router.Handle("/consult", http.HandlerFunc(consult))

	h.Handler = router

	return h
}

//Inicia o server na porta 5000
func StartAPI(routes *httpServer) {
	log.Println("Start API on Port 5000")
	if err := http.ListenAndServe(":5000", routes); err != nil {
		log.Fatal("init server error in StartApi(), ", err)
	}
}

//esta função se encarrega de receber os arquivos enviados pela requisição http e retorna uma resposta do que foi feito
func csvOperator(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		//Lê o file vindo do request
		file, handler, err := r.FormFile("file")
		if err != nil {
			log.Println("FormFile error in csvOperator(), ", err)
		}
		//transforma o file em buffer de bytes
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			log.Println("Copy buffer error in csvOperator(), ", err)
		}
		err1 := saveFile(handler.Filename, buf)
		//esse save vai ser o csv onde será consultado toda vez que tiver uma entrado no endpoint "/consult"
		err2 := saveFile("input-my-app.csv", buf)
		//se houver algum erro no salvamento irá retornar que não foi efetuado com sucesso
		if err1 != nil || err2 != nil {
			_, err = io.WriteString(w, "File Uploaded error\n")
			if err != nil {
				log.Println(" write http.ResponseWriter error in csvOperator(), ", err)
			}
		} else {
			_, err = io.WriteString(w, "File Uploaded successfully\n")
			if err != nil {
				log.Println(" write http.ResponseWriter error in csvOperator(), ", err)
			}

		}

		file.Close()
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

//A consulta recebe uma chamad http com um json:
//
//{
//	"routes":"GRU-CDG"
//}
//
//o formato necessáriamente tem que ser no formato "destinoInicial-DestinoFinal"
//retorna o json com a estrutura:
//
//{
//	"strat": "destinoInicial",
//	"sprintRoutes": "escala1 -> escala2 -> ... -> destinoFinal ",
//	"cost": CustoDaViagem
//
//}
func consult(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		p := struct {
			Destination string `json:"routes"`
		}{}
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fileContentes := conf.OpenFile("input-my-app.csv")
		destinations := conf.ConvertFileInDestinations(fileContentes)
		graph := d.Graph{}
		routes := conf.SplitRoutes(strings.ToUpper(p.Destination))
		cost, start, sprintRoutes := graph.CreatBestDestination(routes[0], routes[1], destinations)
		//
		profile := struct {
			Start        string `json:"start"`
			SprintRoutes string `json:"sprintroutes"`
			Cost         int    `json:"cost"`
		}{
			Start:        start,
			SprintRoutes: sprintRoutes,
			Cost:         cost,
		}
		json.NewEncoder(w).Encode(profile)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

//Salva o csv na pasta MelhorDestino/csv onde vai ser a pasta padrão dos arquivos
func saveFile(name string, file *bytes.Buffer) error {
	// Abre o arquivo usando o parâmetro O_RDWR que permite ler e escrever, O_CREATE cria o arquivo se ele não existir,
	// O_TRUNC torna possível reescrever se o arquivo já existe
	f, err := os.OpenFile("../../csv/"+name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Println("open file error in  csvOperator(), ", err)
	}
	//escre o conteúdo do arquivo passado na requisição http para o arquivo salvo na pasta csv
	_, err = f.Write(file.Bytes())
	if err != nil {
		log.Println("write file error in  csvOperator(), ", err)
		return err
	}
	f.Close()
	return nil
}
