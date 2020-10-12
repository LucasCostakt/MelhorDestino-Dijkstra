package web

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

type httpServer struct {
	http.Handler
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var templates *template.Template

func CreateHttp() *httpServer {
	log.Println("Create new httpServer")
	return new(httpServer)
}

func CreateTemplates() {
	log.Println("Create html Templates")
	templates = template.Must(template.ParseGlob("../../internal/web/templates/*.html"))
}

//Cria as novas rotas
func NewRoutes(h *httpServer) *httpServer {
	log.Println("Init Routes")
	router := http.NewServeMux()
	//criados os endpoint "/", "/consult" e "/save"
	router.Handle("/", http.HandlerFunc(pageHome))
	router.Handle("/save", http.HandlerFunc(save))
	router.Handle("/consult", http.HandlerFunc(pageConsult))

	h.Handler = router

	return h
}

//Inicia o server na porta 5050
func StartAPI(routes *httpServer) {
	log.Println("Start web interface on Port 5050")
	if err := http.ListenAndServe(":5050", routes); err != nil {
		log.Fatal("init server error in StartApi(), ", err)
	}
}

//page home é aonde o usuario vai passar o csv
func pageHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../../internal/web/templates/index.html")
}

//Pagina para consultar a melhor rota
func pageConsult(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//caso o metodo for Get vai exibir o template html sem o conteudo de consulta
	case http.MethodGet:
		err := templates.ExecuteTemplate(w, "consult.html", nil)
		if err != nil {
			log.Println("Cannot Get View pageConsult()", err)
		}
		//caso o metodo for post vai exibir o template html com o conteudo de consulta
	case http.MethodPost:
		//coleta os valores do html
		start := r.FormValue("start")
		final := r.FormValue("final")
		//faz a consulta
		response := RequestConsult(start, final)
		//executa o template passando a melhor rota
		err := templates.ExecuteTemplate(w, "consult.html", response)
		if err != nil {
			log.Println("Cannot Get View pageConsult()", err)
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

//Faz o salvamento do arquivo csv
func save(w http.ResponseWriter, r *http.Request) {
	file, k, err := r.FormFile("file")
	if err != nil {
		log.Println("form file error in save() ", err)
	}
	defer file.Close()
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
	}
	RequestSave(buf, k.Filename)
	//Depois de salvar direciona para a página de consulta
	http.Redirect(w, r, "/consult", 302)
}

func RequestSave(buf *bytes.Buffer, name string) {
	requestSave(buf, &http.Client{}, name)
}

//request passando o arquivo para o salvamento, faz então uma chamada http para a apirest
func requestSave(buf *bytes.Buffer, client httpClient, name string) {
	url := "http://localhost:5000/"
	request, err := NewfileUploadRequest(url, buf, name)
	if err != nil {
		log.Println("Request error in requestSave(), ", err)
	}
	_, err = client.Do(request)
	if err != nil {
		log.Println("Response error in requestSave(), ", err)
	}
}

//Cria o multipart file para fazer o request e retorna *http.Request
func NewfileUploadRequest(url string, buf *bytes.Buffer, name string) (*http.Request, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", name)

	part.Write(buf.Bytes())

	writer.Close()

	request, err := NewRequest(http.MethodPost, url, body.Bytes())
	if err != nil {
		log.Println("Request Error in NewFileUpload(), ", err)
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	return request, err

}

func RequestConsult(start string, final string) string {
	return requestConsult(start, final, &http.Client{})
}

//faz uma chamada http de consulta passando o json de rota e retorna as escalas com o menor valor
func requestConsult(start string, final string, client httpClient) string {
	res := FormatRoute(start, final)
	p := struct {
		Destination string `json:"routes"`
	}{Destination: res}
	js, err := json.Marshal(p)
	if err != nil {
		log.Println("Json marshal error in requestConsult(), ", err)
	}

	request, err := NewRequest(http.MethodPost, "http://localhost:5000/consult", js)
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Println("Request error in requestConsult(), ", err)
	}

	response, err := client.Do(request)
	if err != nil {
		log.Println("Response error in requestConsult(), ", err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Ioutil ReadAll error in requestConsult(), ", err)
	}
	profile := struct {
		Start        string `json:"start"`
		SprintRoutes string `json:"sprintroutes"`
		Cost         int    `json:"cost"`
	}{}

	json.Unmarshal(body, &profile)
	return "Total de custo $" + strconv.Itoa(profile.Cost) + ", e as rotas são: " + profile.Start + " -> " + profile.SprintRoutes
}

//formata os valores vindos do html para fazer o resquest corretamente
func FormatRoute(start string, final string) string {
	s := strings.Trim(start, " ")
	f := strings.Trim(final, " ")
	return strings.ToUpper(s + "-" + f)
}

//Faz os requests
func NewRequest(method string, url string, requestBody []byte) (*http.Request, error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Request error NewRequest() ", err)
		return nil, err
	}
	return request, err
}
