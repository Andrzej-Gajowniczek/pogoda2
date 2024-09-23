package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var url = "http://ai:11434/api/generate"
var llm = "SpeakLeash/bielik-11b-v2.2-instruct:Q4_K_M"

type WeatherData struct {
	Stacja      string `json:"stacja"`
	Temperatura string `json:"temperatura"`
	Cisnienie   string `json:"cisnienie"`
	Zmierzono   string `json:"godzina_pomiaru"`
	Dnia        string `json:"data_pomiaru"`
	Opady       string `json:"suma_opadu"`
	Wilgotność  string `json:"wilgotnosc_wzgledna"`
	Wiatr       string `json:"predkosc_wiatru"`
	Kierunek    string `json:"kierunek_wiatru"`
}

func main() {
	for {
		response, err := http.Get("https://danepubliczne.imgw.pl/api/data/synop")
		if err != nil {
			fmt.Println("Błąd podczas pobierania danych:", err)
			return
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Błąd podczas odczytu odpowiedzi:", err)
			return
		}

		var data []WeatherData

		if err := json.Unmarshal(body, &data); err != nil {
			fmt.Println("Błąd podczas dekodowania JSON:", err)
			return
		}

		czytelnik := bufio.NewReader(os.Stdin)
		var znaleziono bool
		var miasto string
		//for {
		fmt.Print("\npodaj miasto:")
		miasto, _ = czytelnik.ReadString('\n')
		miasto = miasto[:len(miasto)-1]
		//fmt.Scanf("%v", &miasto)
		//fmt.Println("to jest to miasto:", miasto)
		znaleziono = false
		for _, entry := range data {
			//fmt.Println(entry.Stacja)
			if entry.Stacja == miasto {
				prompt := fmt.Sprintf("%v, temperatura: %v, cisnienie: %v, godzinia pomiaru:%v:00, dnia:%v, opady:%v, względna wilgotność powietrza:%v prędkość wiatru:%v, kierunek wiatru:%v\n",
					miasto,
					entry.Temperatura,
					entry.Cisnienie,
					entry.Zmierzono,
					entry.Dnia,
					entry.Opady,
					entry.Wilgotność,
					entry.Wiatr,
					entry.Kierunek)

				znaleziono = true

				zadajPytanie("jesteś prezenterem pogody i komunikujesz w profesjonallny sposób pogodę dla wybranego miasta:" + prompt)
			}
		}
		if !znaleziono {
			fmt.Println("Brak danych na temat miasta:", miasto)
			fmt.Println("Wybierz z poniższych:")
			for _, nazwa := range data {
				fmt.Print(nazwa.Stacja, ",")
			}

		}
	}
}
