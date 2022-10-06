package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"
)

type Assets struct {
	Assets []Asset `json:"commands"`
}

type Asset struct {
	Command             string `json:"command"`
	Label               string `json:"label"`
	Index               int16  `json:"index"`
	Length              int16  `json:"length"`
	Country             string `json:"country"`
	BufferDb            int16  `json:"buffer_db"`
	KeyDb               int16  `json:"key_db"`
	PolygonKey          string `json:"polygon_key"`
	HighFrequencySymbol string `json:"high_frequency_symbol"`
}

func (a Asset) HasFrequencySymbol() bool {
	return a.HighFrequencySymbol != ""
}

func (a Asset) RemoveSpecialChar() string {
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	return re.ReplaceAllString(a.HighFrequencySymbol, "")
}

func (a Asset) GetName() string {
	var name string
	if a.HasFrequencySymbol() {
		highFrequency := a.RemoveSpecialChar()
		name = fmt.Sprintf("zeus-%v-%v-%v-%v-%v",
			strings.Replace(a.Command, "_", "-", -1),
			strings.ToLower(a.Label),
			a.Index,
			a.Length,
			strings.ToLower(highFrequency))
	} else {
		name = fmt.Sprintf("zeus-%v-%v-%v-%v",
			strings.Replace(a.Command, "_", "-", -1),
			strings.ToLower(a.Label),
			a.Index,
			a.Length)
	}
	return name
}

func (a Asset) GenerateSymbolObject() string {
	var symbolName string
	if a.HasFrequencySymbol() {
		symbolName = fmt.Sprintf(
			"'[{\"command\": \"%v\",\"quote_type\": \"%v\",\"index\": %v,\"length\": %v,\"country\": \"%v\",\"high_frequency_symbol\": [\"%v\"]}]'",
			a.Command,
			a.Label,
			a.Index,
			a.Length,
			a.Country,
			strings.Replace(a.HighFrequencySymbol, "\"", "'", -1),
		)
	} else {
		symbolName = fmt.Sprintf(
			"'[{\"command\": \"%v\",\"quote_type\": \"%v\",\"index\": %v,\"length\": %v,\"country\": \"%v\"}]'",
			a.Command,
			a.Label,
			a.Index,
			a.Length,
			a.Country,
		)
	}
	return symbolName
}

func main() {
	// Open json file
	jsonFile, err := os.Open("zeus-to-deploy.json")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	//read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Assets array
	var assets Assets

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &assets)

	// Parsing html file
	t, err := template.ParseFiles("zeus.yaml")

	// standard output
	err = t.Execute(os.Stdout, assets.Assets)
}
