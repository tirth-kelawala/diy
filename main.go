package main

import (
	"github.com/awesomeProject/api"
	"github.com/awesomeProject/factory"
	"log"
)

func main() {
	factory.Init()
	api.HandleRequests()
	log.Println("server started")
}
