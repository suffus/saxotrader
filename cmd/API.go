package main

import (
	"flag"
	"fmt"

	"github.com/suffus/saxotrader"
)

func main() {
	token := flag.String("token", "", "Token for SaxoTrader API")
	flag.Parse()
	if *token == "" {
		fmt.Println("Please provide a token")
		return
	}
	port := saxotrader.NewSaxoAPICall(*token)
	u, err := port.User()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(u)
	c, err := port.Client()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(c)
	a, err := port.Accounts()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(a)
	// get default account
	var defaultAccount saxotrader.SaxoAccount
	defaultAccount = saxotrader.SaxoAccount{}

	for _, acct := range a.Data {
		if acct.AccountKey == c.DefaultAccountKey {
			defaultAccount = acct
		}
	}
	fmt.Println(defaultAccount)
	//get balance
	b, err := port.Balance(defaultAccount.AccountKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(b)
}
