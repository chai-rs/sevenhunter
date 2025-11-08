package main

import (
	"fmt"

	"github.com/chai-rs/sevenhunter/cmd/api/config"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	conf := config.Read()
	tm := conf.Auth.New()
	token, err := tm.SignClaims(jwt.RegisteredClaims{
		Issuer: "issuer",
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(token)
}
