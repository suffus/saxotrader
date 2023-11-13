package main

import (
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
	port := saxotrader.NewPortfolio()
	port.LoadBookings(transF)
	port.GetLatestPositions()

}
