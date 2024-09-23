package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var pl, echo = fmt.Println, fmt.Println

type Invoke struct {
	Model  string `json:model`
	Prompt string `json:prompt`
}

func br() {
	fmt.Println("-----------------------------")
}

func zadajPytanie(prompt string) {
	//url := "http://localhost:11434/api/generate"
	//prompt := "jakim modelem językowym jesteś"
	//to jet pytanie testowe do testowania funkcji zapytanie url można odkomentować jeśli ollama jest zainstalowana lokalnie

	invoke := Invoke{
		Model:  llm,
		Prompt: prompt,
	}
	jSon, err := json.Marshal(invoke)
	if err != nil {
		log.Fatalln("nie mogę zrobić json'a", err)
	}

	//fmt.Println(string(jSon), "\n", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jSon))
	if err != nil {
		pl("Error crating request:", err)
		return
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("User-Agent", "curl/8.2.1")
	br()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		pl("Error making request", err)
		return
	}
	defer resp.Body.Close()

	var target map[string]any

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		jazon := scanner.Bytes()
		err := json.Unmarshal(jazon, &target)
		if err != nil {
			pl("We've got jazon error", err)
			return
		}
		fmt.Print(target["response"])
	}
	if err := scanner.Err(); err != nil {
		log.Println("Error reading response:", err)
	}
	pl()
	br()
}
