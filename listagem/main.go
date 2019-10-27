package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	pb "github.com/DMSec/microservico-hash/listagem/dmsec"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type ClientsDB struct {
	id        int32
	first_name string
	last_name  string
	birthday  string
}

type Campanhas struct {
	id        int32
	campanha string
	status  bool
	pct int32
}

type ProdutosDB struct {
	id        int32
	title string
	description  string
	priceincents  int32
}

func GetConnectionDB()(db *sql.DB) {
	dbDriver := "mysql"

	dbUser := os.Getenv("MYSQL_USER")
	if len(dbUser) == 0{
		dbUser = "root"
	}

	dbPass := os.Getenv("MYSQL_PASSWORD")
	if len(dbPass) == 0{
		dbPass = "root123"
	}

	dbName := os.Getenv("MYSQL_DBNAME")
	if len(dbName) == 0 {
		dbName = "mysql"
	}

	dbHost := os.Getenv("MYSQL_HOST")
	if len(dbHost) == 0 {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("MYSQL_PORT")
	if len(dbPort) == 0{
		dbPort = "3306"
	}

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func getDescontoConnection(host string) (*grpc.ClientConn, error) {
	wd, _ := os.Getwd()
	parentDir := filepath.Dir(wd)
	certFile := filepath.Join(parentDir, "keys", "cert.pem")
	creds, _ := credentials.NewClientTLSFromFile(certFile, "")
	return grpc.Dial(host, grpc.WithTransportCredentials(creds))
}

func setBlackfriday(status bool, pct int32) {
	db := GetConnectionDB()
	insForm, err := db.Prepare("UPDATE campanhas SET status=?, pct=? WHERE campanha='Blackfriday'")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(status, pct)
	if status {
		log.Println("Blackfriday ativado!")
	}else{
		log.Println("Blackfriday desativado!")
	}
	defer db.Close()
}

func ReturnOneCliente(id int ) ClientsDB {
	db := GetConnectionDB()
	selDB, err := db.Query("SELECT * FROM clientes WHERE id=?", int32(id))
	if err != nil {
		panic(err.Error())
	}
	cliente := ClientsDB{}
	for selDB.Next() {
		var idcli int
		var first_name, last_name, birthday string
		err = selDB.Scan(&idcli, &first_name, &last_name, &birthday)
		if err != nil {
			panic(err.Error())
		}
		cliente.id = int32(idcli)
		cliente.first_name = first_name
		cliente.last_name = last_name
		cliente.birthday = birthday
	}
	defer db.Close()
	return cliente
}

func findClienteByID(id int) (pb.Cliente, error) {
	cliente := ReturnOneCliente(id)
	cli := pb.Cliente{Id: cliente.id, FirstName: cliente.first_name, LastName: cliente.last_name, Birthday: cliente.birthday}
	clientes := map[int]pb.Cliente{
		id: cli,
	}
	found, ok := clientes[id]
	if ok {
		log.Println("Cliente:", found.Id , "-", found.FirstName ,"encontrado.")
		return found, nil
	}else {
		log.Println("Cliente nao encontrado")
		return found, errors.New("Cliente nao encontrado.")
	}
}

func getProdutos() []*pb.Produto {
	db := GetConnectionDB()
	selDB, err := db.Query("SELECT * FROM produtos ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	res := []*pb.Produto{}
	for selDB.Next() {
		var id, priceincents int
		var title, description string
		err = selDB.Scan(&id, &title, &description, &priceincents)
		if err != nil {
			panic(err.Error())
		}
		produtos := pb.Produto{Id:int32(id),Title:title, Description:description, PriceInCents:int32(priceincents)}

		res = append(res, &produtos)
	}

	defer db.Close()
	return res
}

func getProductsWithDiscountApplied(cliente pb.Cliente, produtos []*pb.Produto) []*pb.Produto {
	host := os.Getenv("DISCOUNT_SERVICE_HOST")
	if len(host) == 0 {
		host = "localhost:11443"
	}
	conn, err := getDescontoConnection(host)
	if err != nil {
		log.Fatalf("Nao é possivel conectar no servico: %v", err)
	}
	defer conn.Close()
	c := pb.NewDescontoClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	produtosComDescontoAplicado := make([]*pb.Produto, 0)
	for _, produto := range produtos {
		r, err := c.AplicarDesconto(ctx, &pb.DescontoRequisicao{Cliente: &cliente, Produto: produto})
		if err == nil {
			produtosComDescontoAplicado = append(produtosComDescontoAplicado, r.GetProduto())
		} else {
			log.Println("Falha para aplicar desconto.", err)
		}
	}
	if len(produtosComDescontoAplicado) > 0 {
		return produtosComDescontoAplicado
	}
	return produtos
}
func handleGetProducts(w http.ResponseWriter, req *http.Request) {
	produtos := getProdutos()
	w.Header().Set("Content-Type", "application/json")
	clienteID := req.Header.Get("X-USER-ID")
	if clienteID == "" {
		json.NewEncoder(w).Encode(produtos)
		return
	}
	id, err := strconv.Atoi(clienteID)
	if err != nil {
		log.Println("Cliente ID incorreto. ", err)
		http.Error(w, "Cliente ID incorreto", http.StatusBadRequest)
		return
	}

	cliente, err := findClienteByID(id)
	if err != nil {
		json.NewEncoder(w).Encode(produtos)
		return
	}
	produtosComDescontoAplicado := getProductsWithDiscountApplied(cliente, produtos)
	json.NewEncoder(w).Encode(produtosComDescontoAplicado)
}

func handleBlackFriday(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	blackfriday := req.Header.Get("blackfriday")
	pct_blackfriday := req.Header.Get("pct")
	if blackfriday == "" {
		log.Println("Requisição: Ativação da Campanha Blackfriday está sem valor")
		http.Error(w, "Requisição: Ativação da Campanha Blackfriday está sem valor", http.StatusBadRequest)
		return
	}
	if pct_blackfriday == "" {
		log.Println("Requisição: Pct da Campanha Blackfriday está sem valor")
		http.Error(w, "Requisição: Pct da Campanha Blackfriday está sem valor", http.StatusBadRequest)
		return
	}

	ativar, err := strconv.Atoi(blackfriday)
	if err != nil {
		log.Println("Requisição: Ativação da Campanha Blackfriday com valor nulo")
		http.Error(w, "Requisição: Ativação da Campanha Blackfriday com valor nulo", http.StatusBadRequest)
		return
	}

	pct, err := strconv.Atoi(pct_blackfriday)
	if err != nil {
		log.Println("Requisição: Pct da Campanha Blackfriday com valor nulo")
		http.Error(w, "Requisição: Pct da Campanha Blackfriday com valor nulo", http.StatusBadRequest)
		return
	}

	x := ativar
	newBool := !(x == 0) // returns false
	setBlackfriday(newBool, int32(pct))

}

func main() {
	port := "11080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "It is working.")
	})

	http.HandleFunc("/blackfriday", handleBlackFriday)
	http.HandleFunc("/products", handleGetProducts)
	log.Println("Iniciado o serviço de Listagem na porta: ", port)
	http.ListenAndServe(":"+port, nil)
}