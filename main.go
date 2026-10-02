package main

import (
	"log"

	"github.com/MaxMoskalenko/tx-technical-challenge/cmd"
)

func main() {
	if err := cmd.Start(); err != nil {
		log.Fatal("error:", err.Error())
	}
}
