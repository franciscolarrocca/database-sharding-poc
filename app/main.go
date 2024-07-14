package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/serialx/hashring"
)

const (
	API_PORT = "8081"

	DATABASE_PORT_1 = "5431"
	DATABASE_PORT_2 = "5433"
	DATABASE_PORT_3 = "5434"

	API_URI_POST = "/api/products/post"
	API_URI_GET  = "/api/products/get"
)

var (
	databaseClients map[string]*pgx.Conn
	hashRing        *hashring.HashRing
)

func init() {
	hashRing = hashring.New([]string{
		DATABASE_PORT_1,
		DATABASE_PORT_2,
		DATABASE_PORT_3,
	})

	databaseClients = map[string]*pgx.Conn{
		DATABASE_PORT_1: connect(DATABASE_PORT_1),
		DATABASE_PORT_2: connect(DATABASE_PORT_2),
		DATABASE_PORT_3: connect(DATABASE_PORT_3),
	}
}

func main() {
	http.HandleFunc(API_URI_POST, postProductHandler)
	http.HandleFunc(API_URI_GET, getProductHandler)

	log.Printf("Application listening on port %s...", API_PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", API_PORT), nil))
}

func connect(port string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://test:test@localhost:%s/test_db?sslmode=disable", port))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return conn
}

func postProductHandler(w http.ResponseWriter, r *http.Request) {
	productName := r.URL.Query().Get("product_name")
	hash := sha256.Sum256([]byte(productName))
	productCode := hex.EncodeToString(hash[:])[:5]

	databasePort, _ := hashRing.GetNode(productCode)
	databaseClient := databaseClients[databasePort]

	_, err := databaseClient.Exec(context.Background(), "INSERT INTO products(product_name, product_code) VALUES ($1, $2)", productName, productCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"product_name":  productName,
		"product_code":  productCode,
		"database_port": databasePort,
	}

	json.NewEncoder(w).Encode(response)
}

func getProductHandler(w http.ResponseWriter, r *http.Request) {
	productCode := r.URL.Query().Get("product_code")
	databasePort, _ := hashRing.GetNode(productCode)
	databaseClient := databaseClients[databasePort]

	var productName string
	err := databaseClient.QueryRow(context.Background(), "SELECT product_name FROM products WHERE product_code = $1", productCode).Scan(&productName)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"product_name":  productName,
		"product_code":  productCode,
		"database_port": databasePort,
	}

	json.NewEncoder(w).Encode(response)
}
