package main

//Estudo de API Rest em Golang - utilizando o roteador do framework Gorilla Mux: github.com/gorilla/mux
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv" //(biblioteca para converter tipos)

	"github.com/gorilla/mux"
)

//representacao do livro utilizado na api
type Livro struct {
	//os identifadores em json vai permitir colocar letras minusculas
	Id     int    `json:"id"`
	Titulo string `json:"titulo"`
	Autor  string `json:"autor"`
}

//Como nao se usuaremos um banco de dados, construímos um slice de structs de livros
var Livros []Livro = []Livro{
	Livro{
		Id:     1,
		Titulo: "O Guarani",
		Autor:  "José de Alencar",
	},
	Livro{
		Id:     2,
		Titulo: "A Carne",
		Autor:  "Juliano Ribeiro",
	},
	Livro{
		Id:     3,
		Titulo: "Dom Casmurro",
		Autor:  "Machado de Assis",
	},
}

// Ela recebe o w com responseWriter (escritor de argumentos) e r com ponteiro do Request para requisicoes de pedido
//cada http.HandleFunc é uma rota do servidor
func rotaPrincipal(w http.ResponseWriter, r *http.Request) {
	//w - Essa funcao irá dar a resposta ao cliente
	fmt.Fprintf(w, "Bem vindo")
}

//Funcao para acessar a lista de livros "GET"
func listarLivros(w http.ResponseWriter, r *http.Request) {
	//irá formatar o formato json em todas as funções do código
	//w.Header().Add("Content-Type", "application/json")
	//cria lista de livros no formato Json ("encoder" = decodificador)
	encoder := json.NewEncoder(w)
	//método encode que traz a lista de livros
	encoder.Encode(Livros)
	//Pode ser escrito sem criar uma váriavel encoder ex:-> json.NewEncoder(w).Encode(Livros)
}

//Funcao para criar um livro "POST"
func cadastrarLivro(w http.ResponseWriter, r *http.Request) {
	//irá formatar o formato json em todas as funções do código
	//w.Header().Add("Content-Type", "application/json")
	//iouil.ReadAll é uma funcao que recebe como resposta o r.Body em um slice de bytes para permitir criar livro novo
	//Implementa statuscode 201 (StatusCreated)
	w.WriteHeader(http.StatusCreated)
	//tratamento de error na forma de slice de bytes
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) //teste
		return
	}

	//transforma o slice de bites em json. Recebe dois args sendo o primeiro o slice body e o segundo é uma
	//variavel para armazenar o json "&novoLivro"
	var novoLivro Livro
	//unmarshal da biblioteca json irá decodificar o body que está em bytes para json a ser inserido na struct novo livro
	json.Unmarshal(body, &novoLivro)
	novoLivro.Id = len(Livros) + 1
	//adiciona o novo livro ao final da lista (slice)
	Livros = append(Livros, novoLivro)

	encoder := json.NewEncoder(w)
	//método encode que adiciona o novo livro
	encoder.Encode(novoLivro)

}

//método (endpoint) "DELETE"
func excluirLivro(w http.ResponseWriter, r *http.Request) {
	//w.Header().Add("Content-Type", "aplication/json")
	//funcao para quebrar URL em partes de interesse
	//ex: "livros/123" --> ["", "livros","123"]
	vars := mux.Vars(r)

	//variável para converter a string para o tipo int, a partir da biblioteca strconv
	//strconv.Atoi recebe dois parametros sendo o segundo um error, que será substituido por "_"
	id, err := strconv.Atoi(vars["livroId"])
	//condição para o caso de o id passado na busca não existir
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//inicia a variavel em -1 para dizer que a busca não encontrou o indice desejado
	indiceDoLivro := -1

	//loop para localizar o livro no slice de structs (linha 23)
	for indice, livro := range Livros {
		if livro.Id == id {
			//caso encontre o Id desejado irá atribuí-lo a variável indiceDoLivro
			indiceDoLivro = indice
			//caso encontre o livro, break para interromper o loop
			break
		}
	}
	//caso não encontre o livro (indice -1)
	if indiceDoLivro < 0 {
		w.WriteHeader(http.StatusNotFound)
		//return para interromper a ação
		return
	}
	//atualização do indice dos livros da lista excluindo o indice desejado
	//ex: [a, b, c, d] - b ---> [a] + [c, d] => [a, c, d]
	ladoEsquerdo := Livros[0:indiceDoLivro]
	ladoDireito := Livros[indiceDoLivro+1 : len(Livros)]
	//as reticências indicam que é um argumento separado por virgula
	Livros = append(ladoEsquerdo, ladoDireito...)
	//após a exclusão irá indicar que o elemento excluido não pode ser encontrado caso
	//seja requisitado novamente (StatusNoContent)
	w.WriteHeader(http.StatusNoContent)
}

//Metodo "PUT"
func modificarLivro(w http.ResponseWriter, r *http.Request) {
	//w.Header().Add("Content-Type", "aplication/json")
	//funcao para quebrar URL em partes de interesse
	//ex: "livros/123" --> ["", "livros","123"]
	vars := mux.Vars(r)

	//variável para converter a string para o tipo int, a partir da biblioteca strconv
	//strconv.Atoi recebe dois parametros sendo o segundo um error, que será substituido por "_"
	id, err := strconv.Atoi(vars["livroId"])

	//request com error de ID: condição para o caso de o id passado na busca não existir irá entregar o BadRequest como resposta
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//Caso o request seja feito com o exito...
	//iremos passar os dados a serem alterados no Body da requisição http
	//a funcão ioutil.ReadAll recebe o body da requisição e o erro.
	body, BodyErr := ioutil.ReadAll(r.Body)

	if BodyErr != nil {
		//error de status 500 caso haja problema ao ler o corpo da requisição
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//nesta etapa transforma o livro em slice de bytes em formato Json
	var livroModificado Livro
	//irá converter o que recebe no corpo para o tipo livroModificado
	errJson := json.Unmarshal(body, &livroModificado)
	//caso haja error na decodificação em Json
	if errJson != nil {
		//bad request caso o usuário digite algo errado na requisição de modificação do livro
		w.WriteHeader(http.StatusBadRequest)
		//o return irá interromper a requisição de modificação
		return
	}

	//caso haja êxito na requisição
	//estabelece que o não há livro localizado com o indice "-1"
	indiceDoLivro := -1
	//loop para localizar o livro no slice de structs (linha 23)
	for indice, livro := range Livros {
		//caso o indice do livro seja igual ao id solicitado
		if livro.Id == id {
			//o indice será encontrado e associado a variável
			indiceDoLivro = indice
			//break para interromper o loop de busca do livro
			break
		}
	}
	//caso não encontre o livro (indice -1)
	if indiceDoLivro < 0 {
		//caso não encontre o livro
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//caso encontre o livro irá modificar os dados do livro
	Livros[indiceDoLivro] = livroModificado

	//funcao da biblioteca json irá apresentar o livro já alterado
	json.NewEncoder(w).Encode(livroModificado)

}

func buscarLivro(w http.ResponseWriter, r *http.Request) {
	//irá formatar o formato json em todas as funções do código
	//w.Header().Add("Content-Type", "application/json")
	//funcao para quebrar URL em partes de interesse
	//ex: "livros/123" --> ["", "livros","123"]
	//partes := strings.Split(r.URL.Path, "/")

	vars := mux.Vars(r)
	//variável para converter a string para o tipo int, a partir da biblioteca strconv
	//strconv.Atoi recebe dois parametros sendo o segundo um error, que será substituido por "_"
	id, _ := strconv.Atoi(vars["livroId"])

	//loop para localizar o livro no slice de structs (linha 23)
	for _, livro := range Livros {
		if livro.Id == id {
			//ao localizar irá codificar o livro em json
			json.NewEncoder(w).Encode(livro)
			return
		}
	}
	//caso não encontre o livro
	w.WriteHeader(http.StatusNotFound)
}

func configurarRotas(roteador *mux.Router) {

	//func http.HandleFunc que recebe "/"(diretorio raiz do server) e a func rotaPrincipal com o chamado do argumento para lidar com a porta.
	//a func http.HandleFunc pode ser utilizada para fazer outras rotas diferentes da principal

	roteador.HandleFunc("/", rotaPrincipal)
	//rota para listar livros-GET
	roteador.HandleFunc("/Livros", listarLivros).Methods("GET")
	//roteador para rotas livros-GET ou POST
	roteador.HandleFunc("/livros/{livroId}", buscarLivro).Methods("GET")
	//rota para cadastrar livros-POST
	roteador.HandleFunc("/livros", cadastrarLivro).Methods("POST")
	//rota para modificar livro-PUT
	roteador.HandleFunc("/livros/{livroId}", modificarLivro).Methods("PUT")
	roteador.HandleFunc("/livros/{livroId}", excluirLivro).Methods("DELETE")
}

//func que recebe com parametro next que recebe um handler e retorna da func middleware um handler
//permite remover todos os content-type das funções e insere o formato json em todo o código
func jsonMiddleware(next http.Handler) http.Handler {
	//retorno de uma função como argumento
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// next -> segue para o próximo middleware onde deverá ser incluído um header em json
		next.ServeHTTP(w, r)
	})
}

//func para encapisular a lógica para subida do servidor na web
func configurarServidor() {
	//configaracao do gorilla mux - criando roteador
	roteador := mux.NewRouter().StrictSlash(true)
	//clean code -> funcão para melhorar o uso de json em cada rota do código
	roteador.Use(jsonMiddleware)
	configurarRotas(roteador)
	fmt.Println("Servidor rodando na porta: 8888")
	//func http.ListenAndServe que recebe a porta do localhost e o servidor mux (nil = DefaultServerMux)
	//Esta func faz rodar o servidor - por causa do "nil" será rodado o multiplexer Default server Mux
	log.Fatal(http.ListenAndServe(":8888", roteador))
}

func main() {

	configurarServidor()

}
