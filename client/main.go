package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Valor struct {
	Bid string
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		fmt.Println("Erro na requisicao")
	}
	res, _ := client.Do(req)

	answer, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	var v Valor
	err = json.Unmarshal(answer, &v)
	file_answer := fmt.Sprintln("Dolar:", v.Bid)

	arq, err := os.Open("cotacao.txt")
	defer arq.Close()

	if err != nil {
		arq, _ = os.Create("cotacao.txt")
		arq.WriteString(file_answer)
	}

	arq.WriteString(file_answer)
	fmt.Println(string(answer))
}
