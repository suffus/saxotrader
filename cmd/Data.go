package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/suffus/saxotrader"
)

func main() {
	transFilename := os.Args[1]

	transF, err := os.Open(transFilename)
	if err != nil {
		panic(err)
	}
	//	bytesF, _ := io.ReadAll(transF)

	transactions := make([]*saxotrader.Booking, 0)
	r := csv.NewReader(transF)
	r.Read() // get rid of first line

	for {
		booking := saxotrader.Booking{}
		err := saxotrader.Unmarshal(r, &booking)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		transactions = append(transactions, &booking)
	}

	for _, b := range transactions {
		bytes, _ := json.Marshal(b)
		fmt.Println(string(bytes))
	}
}
