package main

import (
	"context"
	"encoding/json"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type Cambio struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type CotacaoDTO struct {
	Bid string `json:"bid"`
}

type Cotacao struct {
	ID  int
	Bid string
}

func autoSchemaGenerate() {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Cotacao{})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	autoSchemaGenerate()
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(response http.ResponseWriter, request *http.Request) {
	cambio, err := getCambio()
	if err != nil {
		log.Println("ocorreu um erro ao realizar a chamada da api de cambio: ", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	cotacao, err := save(cambio)
	if err != nil {
		log.Println("ocorreu um erro ao realizar o insert de uma cotacao: ", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	var cambioDto = CotacaoDTO{Bid: cotacao.Bid}
	json.NewEncoder(response).Encode(cambioDto)

}

func getCambio() (*Cambio, error) {
	client := http.Client{
		Timeout: time.Millisecond * 200,
	}

	resp, err := client.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cambio Cambio
	err = json.NewDecoder(resp.Body).Decode(&cambio)
	if err != nil {
		return nil, err
	}

	return &cambio, nil
}

func save(cambio *Cambio) (*Cotacao, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	cotacao := Cotacao{Bid: cambio.USDBRL.Bid}
	response := db.WithContext(ctx).Create(&cotacao)
	if response.Error != nil {
		return nil, err
	}

	log.Println("Contacao salva com sucesso: ", cotacao)

	return &cotacao, nil
}
