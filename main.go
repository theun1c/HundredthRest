package main

import "net/http"

func main() {

	// базовая сигнатура хендл функции
	http.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {})
}
