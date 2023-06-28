package main

import (
	"encoding/json"
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type County struct {
	Name string `json:"municipio"`
	Cod  string `json:"cod"`
}

type State struct {
	Name   string   `json:"estado"`
	County []County `json:"municipios"`
}

func NewCounty(name string, cod string) *County {
	return &County{
		Name: name,
		Cod:  cod,
	}
}

func NewState(name string, county []County) *State {
	return &State{
		Name:   name,
		County: county,
	}
}

func getCounty(c *colly.Collector, counties *[]County) {
	c.OnHTML("body > section > article > div.container-codigos > table:nth-child(21) > tbody",
		func(h *colly.HTMLElement) {
			h.ForEach("tr", func(i int, el *colly.HTMLElement) {
				county := NewCounty(
					el.ChildText("td:nth-child(1) > a"),
					el.ChildText("td:nth-child(2)"),
				)
				*counties = append(*counties, *county)
			})
		})
}

func writeFile(filename string, s State) {
	file, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(filename, file, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("File", filename, "created")
}

func writeCSV(filename string, s State) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"municipio", "cod"}
	err = writer.Write(header)
	if err != nil {
		log.Fatal(err)
	}

	for _, county := range s.County {
		row := []string{county.Name, county.Cod}
		err = writer.Write(row)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("CSV file", filename, "created")
}


func main() {
	c := colly.NewCollector()
	counties := make([]County, 0)

	getCounty(c, &counties)

	c.OnScraped(func(r *colly.Response) {
		dir := "../output"
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}

		s := NewState("Rio de Janeiro", counties)
		filepath := dir + "/" + "rio_de_janeiro_cod"
		json := filepath + ".json"
		csv := filepath + ".csv"
		writeFile(json, *s)
		writeCSV(csv, *s)
	})

	c.Visit("https://ibge.gov.br/explica/codigos-dos-municipios.php")
}
