package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type JsonAnswer struct {
	Usdbrl struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

type Cotacao struct {
	Bid string
}

func main() {

	http.HandleFunc("/cotacao", Server)
	http.ListenAndServe(":8080", nil)
}

func Server(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	client := http.Client{}
	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro ao realizar o GET!")
	}

	select {
	case <-time.After(time.Millisecond * 200):
		ctx.Done()
	default:
		answer := resp.Body
		defer answer.Close()
		answer_string, _ := io.ReadAll(answer)

		var answer_to_request JsonAnswer
		err = json.Unmarshal(answer_string, &answer_to_request)
		if err != nil {
			fmt.Println("Erro ao Unmarshal a resposta.")
		}
		date := time.Now().Format("2006-01-02 15:04:05 -07:00:00")
		InsertInDatabase(answer_to_request.Usdbrl.Bid, date)

		cot := Cotacao{answer_to_request.Usdbrl.Bid}
		answer_to_client, _ := json.Marshal(cot)

		w.Header().Add("Content-Type", "application/json")
		w.Write(answer_to_client)

	}
}

func InsertInDatabase(value string, date string) {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)

	defer cancel()

	db := AcessDataBase()
	defer db.Close()

	CreateTable(db)

	stm, err := db.Prepare("INSERT INTO cotacao(valor, date) VALUES(?, ?);")
	if err != nil {
		fmt.Println("Erro ao adicionar dados.")
	}
	defer stm.Close()
	_, err = stm.Exec(value, date)
	if err != nil {
		fmt.Println("Erro ao executar a insercao")
	}
}

func AcessDataBase() *sql.DB {
	db, _ := sql.Open("mysql", "root:root@tcp(localhost:3306)/cotacao")
	return db
}

func CreateTable(db *sql.DB) {
	stm_create_table, err := db.Prepare("CREATE TABLE IF NOT EXISTS cotacao(valor VARCHAR (10) NOT NULL, date VARCHAR (30) NOT NULL);")
	if err != nil {
		fmt.Println("Erro ao criar tabela.")
	}
	defer stm_create_table.Close()
	stm_create_table.Exec()
}
