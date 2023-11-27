package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func main() {
	token := flag.String("code", "", "Authorization Code for SaxoTrader API from OAuth2 code workflow")
	client_id := flag.String("client_id", "", "Client ID for SaxoTrader API")
	client_secret := flag.String("client_secret", "", "Client Secret for SaxoTrader API")
	flag.Parse()
	if *token == "" {
		fmt.Println("Please provide a token")
		return

	}
	if *client_id == "" {
		fmt.Println("Please provide a client_id")
		return

	}
	if *client_secret == "" {
		fmt.Println("Please provide a client_secret")
		return

	}
	tokenEndpoint := "https://sim.logonvalidation.net/token"

	client := &http.Client{}

	resp, err := client.PostForm(tokenEndpoint, url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {*token},
		"client_id":     {*client_id},
		"client_secret": {*client_secret},
		"redirect_uri":  {"http://fanjango.com.hk/auth.html"},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp.StatusCode != 201 {
		fmt.Println("Status code not 201, instead got ", resp.StatusCode)
		return
	}
	b, err := io.ReadAll(resp.Body)
	fmt.Println(string(b))

}
