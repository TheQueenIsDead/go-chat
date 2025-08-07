package main

import (
	"go-chat/pkg"
	"log"
)

func main() {
	chat := pkg.NewServer()
	err := chat.Run()
	if err != nil {
		log.Panic(err)
	}
}
