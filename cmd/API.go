package main

import (
	"flag"
	"fmt"

	"github.com/suffus/saxotrader"
)

func main() {
	token := flag.String("token", "", "Token for SaxoTrader API")
	symbol := flag.String("symbol", "", "Symbol for SaxoTrader API")
	exchange := flag.String("exchange", "", "Exchange for SaxoTrader API")
	assetType := flag.String("assetType", "", "Asset Type for SaxoTrader API")
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
	a, err := port.Accounts()
	if err != nil {
		fmt.Println(err)
		return
	}
	// get default account
	var defaultAccount saxotrader.SaxoAccount
	defaultAccount = saxotrader.SaxoAccount{}
	fmt.Println("there are ", len(a.Data), " accounts")
	for _, acct := range a.Data {
		if acct.AccountKey == c.DefaultAccountKey {
			defaultAccount = acct
		}
	}
	fmt.Println(defaultAccount.AccountName, defaultAccount.AccountKey, defaultAccount.AccountType, defaultAccount.AccountId)
	//get balance
	b, err := port.Balance()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Cash available is ", b.CashBalance)
	pos, err := port.NetPositions(saxotrader.SaxoInstruction{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("there are ", len(pos), " positions")
	for i := range pos {
		details, err := port.InstrumentDetails(saxotrader.SaxoInstruction{Uic: pos[i].NetPositionBase.Uic})
		if err != nil {
			fmt.Println(err)
			//return
		}
		sym := details[0].Symbol
		fmt.Println(pos[i].NetPositionBase.Uic, sym, pos[i].NetPositionBase.Amount, pos[i].NetPositionBase.ValueDate, pos[i].NetPositionView.CurrentPrice)
	}
	instr, err := port.Instruments(saxotrader.SaxoInstruction{AssetTypes: []string{*assetType}, Keywords: *symbol, ExchangeId: *exchange})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(instr)
	for i := range instr {
		fmt.Println(instr[i].Symbol, instr[i].Description, instr[i].AssetType, instr[i].ExchangeId, instr[i].Identifier, instr[i].PrimaryListing)
		details, err := port.InstrumentDetails(saxotrader.SaxoInstruction{Uic: instr[i].Identifier})
		if err != nil {
			fmt.Println(err)
			//return
		} else {
			fmt.Println(details)
		}
	}
	fmt.Println(len(instr))
}
