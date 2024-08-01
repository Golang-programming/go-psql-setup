package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-psql-setup/model"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Failed to load environment")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	return db

}

func GetStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Having issue with params: %v", err)
	}

	stock, err := getStock(int64(id))

	if err != nil {
		log.Fatalf("Unable to get stock: %v", err)
	}

	json.NewEncoder(w).Encode(stock)

}
func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Having issue with params: %v", err)
	}

	var stock model.Stock

	err = json.NewDecoder((r.Body)).Decode(&stock)

	if err != nil {
		log.Fatal("Error to decoding body: ", err)
	}

	updatedRows := updateStock(int64(id), stock)

	msg := fmt.Sprintf("Stock successfully updated. Total rows/record affected %v", updatedRows)

	if err != nil {
		log.Fatalf("Unable to update stock: %v", err)
	}

	response := Response{
		Message: msg,
		ID:      int64(id),
	}

	json.NewEncoder(w).Encode(response)
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Having issue with params: %v", err)
	}

	deletedRows := deleteStock(int64(id))

	msg := fmt.Sprintf("Stock successfully deleted. Total rows/record affected %v", deletedRows)

	response := Response{
		Message: msg,
		ID:      int64(id),
	}

	json.NewEncoder(w).Encode(response)
}

func GetAllStocks(w http.ResponseWriter, r *http.Request) {
	stocks, err := getAllStocks()

	if err != nil {
		log.Fatalf("Unable to get all stocks: %v", err)
	}

	json.NewEncoder(w).Encode(stocks)
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock model.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatal("Error to decoding body: ", err)
	}

	insertID := insertStock(stock)

	response := Response{
		Message: "Stock successfully created",
		ID:      insertID,
	}

	json.NewEncoder(w).Encode(response)
}

//------------------------- handler functions ----------------

func insertStock(stock model.Stock) int64 {
	var id int64
	db := createConnection()
	defer db.Close()

	sqlStatement := "INSERT INTO stocks(name, price, compant) VALUES ($1, $2, $3) RETURNING id"
	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)
	if err != nil {
		log.Fatal("Error inserting stock: %v", err)
	}

	return id
}

func getStock(id int64) (model.Stock, error) {
	stock := model.Stock{}
	db := createConnection()
	defer db.Close()

	sqlStatement := "SELECT * FROM stocks WHERE id=$1"

	err := db.QueryRow(sqlStatement, id).Scan(&stock.ID, &stock.Name, &stock.Price, &stock.Company)

	switch err {

	case sql.ErrNoRows:
		fmt.Println("No stock found")
		return stock, nil
	case nil:
		return stock, nil
	default:
		log.Fatal("Unable to execute the query: %v", err)
	}

	return stock, nil
}

func getAllStocks() ([]model.Stock, error) {
	var stocks []model.Stock
	db := createConnection()
	defer db.Close()

	sqlStatement := "SELECT * FROM stocks "

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatal("Unable to execute get stocks query: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var stock model.Stock

		err = rows.Scan(&stock.ID, &stock.Name, &stock.Price, &stock.Company)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		stocks = append(stocks, stock)
	}

	return stocks, err
}

func updateStock(id int64, stock model.Stock) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `UPDATE stocks SET name=$2, price=$3, company=$4 WHERE stockid=$1`

	res, err := db.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Company)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

func deleteStock(id int64) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM stocks WHERE stockid=$1`

	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
