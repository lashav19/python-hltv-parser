package main

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Items struct {
	XMLName xml.Name `xml:"items"`
	Items   []Item   `xml:"item"`
}

type Item struct {
	Request Request `xml:"request"`
}

type Request struct {
	Base64 string `xml:"base64,attr"`
	Data   string `xml:",chardata"`
}

func main() {
	// Read XML data from file
	xmlFile, err := os.Open("dump.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var items Items
	xml.Unmarshal(byteValue, &items)

	var tokens []string


	for _, item := range items.Items {
		if item.Request.Base64 == "true" {
			decodedRequest, err := base64.StdEncoding.DecodeString(item.Request.Data)
			if err != nil {
				fmt.Println("Error decoding request:", err)
				continue
			}
			requestStr := string(decodedRequest)
			for _, line := range strings.Split(requestStr, "\n") {
				if strings.HasPrefix(line, "Postman-Token:") {
					tokens = append(tokens, strings.TrimSpace(line))
				}
				}
			}
		}
	}

	// Write tokens to file
	err = ioutil.WriteFile("tokens.txt", []byte(strings.Join(tokens, "\n")), 0644)
	if err != nil {
		fmt.Println("Error writing tokens to file:", err)
	}
