package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type viaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type apiCEP struct {
	Code     string `json:"code"`
	State    string `json:"state"`
	City     string `json:"city"`
	District string `json:"district"`
	Address  string `json:"address"`
}

func main() {
	c1 := make(chan viaCEP)
	c2 := make(chan apiCEP)

	// ViaCEP
	go func() {
		for _, cep := range os.Args[1:] {
			req, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao fazer a requisição: %v\n", err)
			}
			defer req.Body.Close()

			res, err := io.ReadAll(req.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao ler a resposta: %v\n", err)
			}

			var data viaCEP
			err = json.Unmarshal(res, &data)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao fazer o unmarshal: %v\n", err)
			}
			c1 <- data
		}
	}()

	// apiCEP
	go func() {
		for _, cep := range os.Args[1:] {
			req, err := http.Get("https://cdn.apicep.com/file/apicep/" + cep + ".json")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao fazer a requisição: %v\n", err)
			}
			defer req.Body.Close()

			res, err := io.ReadAll(req.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao ler a resposta: %v\n", err)
			}

			var data apiCEP
			err = json.Unmarshal(res, &data)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao fazer o unmarshal: %v\n", err)
			}
			c2 <- data
		}
	}()

	select {
	case data := <-c1:
		fmt.Printf("Received from ViaCEP.\n%s", data)

	case data := <-c2:
		fmt.Printf("Received from ApiCEP.\n%s", data)

	case <-time.After(time.Second):
		println("Timeout. Nenhuma chamada da API concluída dentro do prazo")
	}
}
