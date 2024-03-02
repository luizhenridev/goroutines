package main

import (
	"cep/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/gommon/log"
)

func main() {
	cep := "60711520"

	c1 := make(chan model.Endereco)
	c2 := make(chan model.EnderecoDetalhado)

	go firstAPI(cep, c1)
	go secondAPI(cep, c2)

	select {
	case cep1 := <-c1:
		println(fmt.Printf("Essas são as informações do seu CEP de acordo com a Brasil API: %+v\n", cep1))
	case cep2 := <-c2:
		println(fmt.Printf("Essas são as informações do seu CEP de acordo com a Via Cep API: %+v\n", cep2))
	case <-time.After(time.Second * 1):
		println("timeout")
	}

}

func firstAPI(cep string, ch chan<- model.Endereco) {

	req, err := http.NewRequest("GET", fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep), nil)
	if err != nil {
		log.Error(fmt.Println("An error happend while was crating the request"))
		return
	}

	var jsonPayload model.Endereco
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(fmt.Println("An error happend while was crating the request"))
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error(fmt.Println("Failed to read response body"))
		return
	}

	err = json.Unmarshal(body, &jsonPayload)
	if err != nil {
		log.Error(fmt.Println("Error to unmarshal response body"))
	}

	ch <- jsonPayload

}

func secondAPI(cep string, ch chan<- model.EnderecoDetalhado) {

	req, err := http.NewRequest("GET", fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep), nil)
	if err != nil {
		log.Error(fmt.Println("An error happend while was crating the request"))
		return
	}

	var jsonPayload model.EnderecoDetalhado
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(fmt.Println("An error happend while was crating the request"))
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error(fmt.Println("Failed to read response body"))
		return
	}

	err = json.Unmarshal(body, &jsonPayload)
	if err != nil {
		log.Error(fmt.Println("Error to unmarshal response body"))
	}

	ch <- jsonPayload

}
