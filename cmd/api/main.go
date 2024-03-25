package main

import (
	"fmt"
	"personal_budget_app/internal/server"
)

func main() {

	server_ := server.NewServer()

	err := server_.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
