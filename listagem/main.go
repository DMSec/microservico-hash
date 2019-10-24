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
	pb "microservico-hash/listagem/dmsec"
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
}

type ProdutosDB struct {
	id        int32
	title string
	description  string
	priceincents  int32
}

func birthDate(birthDate time.Time, now time.Time) bool {
	days := now.Day() - birthDate.Day()
	months := now.Month() - birthDate.Month();
	retorno := false

	if (days == 0) && (months == 0) {
		retorno = true
		return retorno
	}else {
		return retorno
	}
}

func GetConnectionDB()(db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root123"
	dbName := "mysql"
	dbHost := "172.17.0.2"
	dbPort := "3306"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName)
	//db, err := sql.Open("mysql", "db_user:password@tcp(localhost:3306)/my_db")
	//db, err := sql.Open(dbDriver,dbUser+":"+dbPass+"@tcp"+"("+dbHost+":"+ dbPort +")@/"+<dbName>)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func getDiscountConnection(host string) (*grpc.ClientConn, error) {
	wd, _ := os.Getwd()
	parentDir := filepath.Dir(wd)
	certFile := filepath.Join(parentDir, "keys", "cert.pem")
	creds, _ := credentials.NewClientTLSFromFile(certFile, "")
	return grpc.Dial(host, grpc.WithTransportCredentials(creds))
}

func setBlackfriday(status bool) {
	db := GetConnectionDB()
	fmt.Print("Entrei em seblack")
	insForm, err := db.Prepare("UPDATE campanhas SET status=? WHERE campanha='Blackfriday'")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(status)
	log.Println("UPDATE: status: ",status )
	defer db.Close()
}

func ReturnOneCliente(id int ) ClientsDB {
	fmt.Print(id)
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
		fmt.Println(cliente.id)
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
	c1 := pb.Cliente{Id: cliente.id, FirstName: cliente.first_name, LastName: cliente.last_name, Birthday: cliente.birthday}
	clientes := map[int]pb.Cliente{
		id: c1,
	}
	found, ok := clientes[id]
	if ok {
		return found, nil
	}else {
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

func getFakeProducts() []*pb.Produto {
	p1 := pb.Produto{Id: 1, Title: "iphone-x", Description: "64GB, black and iOS 12", PriceInCents: 99999}
	p2 := pb.Produto{Id: 2, Title: "notebook-avell-g1511", Description: "Notebook Gamer Intel Core i7", PriceInCents: 150000}
	p3 := pb.Produto{Id: 3, Title: "playstation-4-slim", Description: "1TB Console", PriceInCents: 32999}
	return []*pb.Produto{&p1, &p2, &p3}
}
func getProductsWithDiscountApplied(cliente pb.Cliente, produtos []*pb.Produto) []*pb.Produto {
	host := os.Getenv("DISCOUNT_SERVICE_HOST")
	if len(host) == 0 {
		host = "localhost:11443"
	}
	conn, err := getDiscountConnection(host)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
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
	//produtos := getFakeProducts()
	produtos := getProdutos()
	w.Header().Set("Content-Type", "application/json")
	clienteID := req.Header.Get("X-USER-ID")
	fmt.Print(clienteID)
	if clienteID == "" {
		json.NewEncoder(w).Encode(produtos)
		return
	}
	id, err := strconv.Atoi(clienteID)
	if err != nil {
		http.Error(w, "Cliente ID nao e um numero.", http.StatusBadRequest)
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
	fmt.Print("valor e ",blackfriday)
	if blackfriday == "" {
		http.Error(w, "blackfriday nao existe", http.StatusBadRequest)
		return
	}
	value, err := strconv.Atoi(blackfriday)
	if err != nil {
		http.Error(w, "Cliente ID nao e um numero.", http.StatusBadRequest)
		return
	}

	fmt.Print("abc", value)
	x := value
	newBool := !(x == 0) // returns false
	setBlackfriday(newBool)

	//black, err := setBlackfriday(value)
	//if err != nil {
		//json.NewEncoder(w).Encode("Alterado com sucesso - "+black)
		//return
	//}

}

func main() {
	port := "11080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "It is working.")
	})

	//http.handleFunc("/loadProducts", handleLoadProducts)
	//http.handleFunc("/loadUsers", handleLoadUsers)
	http.HandleFunc("/blackfriday", handleBlackFriday)
	http.HandleFunc("/products", handleGetProducts)
	fmt.Println("Running Listagem em", port)
	http.ListenAndServe(":"+port, nil)
}