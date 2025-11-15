package main

import (
	"fmt"

	"github.com/theun1c/HundredthRest/http"
	"github.com/theun1c/HundredthRest/todo"
)

func main() {
	todoList := todo.NewList()
	httpHandlers := http.NewHTTPHandlers(todoList)
	httpServer := http.NewHTTPServer(httpHandlers)

	if err := httpServer.StartServer(); err != nil {
		fmt.Println("failed to start HTTP server", err)

	}
}
