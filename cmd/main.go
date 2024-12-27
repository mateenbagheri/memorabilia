package main

import (
	"github.com/mateenbagheri/memorabilia/server"
)

func main() {
	srv := server.New()
	srv.Start()
}
