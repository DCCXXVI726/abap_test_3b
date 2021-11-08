package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type Valute struct {
	Name	string		`xml:"Name"`
	Value	string		`xml:"Value"`
	Nominal	float64		`xml:"Nominal"`
}

type ValCurs struct {
	Valutes	[]Valute	`xml:"Valute"`
}

func recoveryFunction() {

	if recoveryMessage :=recover(); recoveryMessage != nil {
		fmt.Printf("Не удалось получить стоимость одной норвежской кроны в венгерских форинтах: %s", recoveryMessage)
	}

}

func main() {
	defer recoveryFunction()
	URL := "http://www.cbr.ru/scripts/XML_daily.asp"
	resp, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	
	d := xml.NewDecoder(resp.Body)
	d.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unknown charset %s", charset)
		}
	}

	var valCurs ValCurs

	err = d.Decode(&valCurs)
	if err != nil {
		panic(err)
	}

	valute1 := "Норвежских крон"
	valute2 := "Венгерских форинтов"


	value1, err := findValue(valute1, valCurs.Valutes)
	if err != nil {
		panic(err)
	}

	value2, err := findValue(valute2, valCurs.Valutes)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Cтоимость одной норвежской кроны в венгерских форинтах: %.2f", value1/value2)
}

func findValue(valute string, valutes []Valute) (float64, error) {
	for _, v := range (valutes) {
		if v.Name == valute {
			s := strings.Replace(v.Value, ",", ".", 1)
			tmp, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return 0, err
			}
			return tmp / v.Nominal, nil
		}
	}
	return 0, fmt.Errorf("Can't find %v", valute)
}
