package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	cotacao, err := getCotacao()
	if err != nil {
		panic(err)
	}

	err = fileWrite(cotacao)
	if err != nil {
		panic(err)
	}
}

func getCotacao() (*Cotacao, error) {
	client := http.Client{
		Timeout: time.Millisecond * 300,
	}

	resp, err := client.Get("http://localhost:8080/cotacao")
	if err != nil {
		return nil, err
	}

	var cotacao Cotacao
	err = json.NewDecoder(resp.Body).Decode(&cotacao)
	if err != nil {
		return nil, err
	}

	return &cotacao, nil
}

func fileWrite(cotacao *Cotacao) error {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}

	_, err = file.WriteString("Dólar: " + cotacao.Bid)
	if err != nil {
		return err
	}

	return err
}
