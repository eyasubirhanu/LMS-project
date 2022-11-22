package main

import (
	"fmt"
	"test/api/config"
	"test/api/route"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	configuration := config.New()
	initialized := route.NewInitializedServer(configuration)
	// Run
	PORT := fmt.Sprintf(":%v", configuration.Get("APP_PORT"))
	teenager(PORT)

	err := initialized.Run(PORT)
	if err != nil {
		panic(err)
	}
}

func teenager(port string) {
	fmt.Print(`
start...
`)
}
