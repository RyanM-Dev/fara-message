package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/mhghw/fara-message/api"
	"github.com/mhghw/fara-message/db"
)

// implement this with os args
var port = flag.Int("port", 8080, "Port to run the HTTP server")

func main() {
	db.NewDatabase()
	flag.Parse()
	// err := api.RunWebServer(*port)
	// if err != nil {
	// 	log.Print("failed to start HTTP server:", err)
	// }
	webServer := api.NewWebServer()
	addr := fmt.Sprintf(":%d", *port)
	fmt.Println("address is:", addr)
	err := webServer.Run(addr)
	if err != nil {
		log.Print("failed to start HTTP server:", err)
	}

}
