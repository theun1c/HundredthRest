package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
)

// не нужно использовать атомик, поскольку секция с возможной
// гонкой данных обрабатывается мьютексом
var money int = 1000 // денег в кармане

// с банком тоже самое
var bank int = 0 // денег в копилке

// создаем мьютекс для критической секции
// в которой 2 разных горутины (хендлера)
// работают с одним атомиком маней. *уже не атомиком. я убрал атомики
//
// при одновременном входе в блоки с обработкой этой переменной
// может случиться гонка данных, поскольку 2 горутины
// не знают друг про друга и работают с 1 переменной
var mtx = sync.Mutex{}

// когда вызывается хендлер - он вызывается в отдельной горутине
// а это означает, что нам необходимо обеспечить контроль данных,
// с которыми будет осуществляться работа.
// мы работаем с глобальной переменной, а значит
// money -= amount может привести к гонке данных
// поэтому следует использовать либо mtx либо atomic
//
// если выбор стоит между атомиками и мьютексами, то лучше всего использовать атомики
func payHandler(w http.ResponseWriter, r *http.Request) {
	httpRequestBoby, err := io.ReadAll(r.Body)

	if err != nil {
		msg := "Failed to read HTTP body: " + err.Error()
		fmt.Println(msg)
		w.Write([]byte(msg))
		return
	}

	httpRequestBobyStr := string(httpRequestBoby)
	amount, err := strconv.Atoi(httpRequestBobyStr)

	if err != nil {
		msg := "Request convertation failed: " + err.Error()
		fmt.Println(msg)
		w.Write([]byte(msg))
	}

	// тот самый блок с возможной гонкой данных на переменной маней
	// обращение из 2х независимых горутин к 1 переменной
	mtx.Lock()
	if money-amount >= 0 {
		money -= amount
		msg := "Succsessfuly payment, now money is: " + strconv.Itoa(money)
		fmt.Println(msg)
		_, err = w.Write([]byte(msg))

		if err != nil {
			fmt.Println("Failed to write HTTP response: ", err)
		}
	}
	mtx.Unlock()
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	httpRequestBody, err := io.ReadAll(r.Body)

	if err != nil {
		msg := "Failed to read HTTP body " + err.Error()
		fmt.Println(msg)
		w.Write([]byte(msg))
		return
	}

	// преобразуем слайс байт в строку
	httpRequestBodyStr := string(httpRequestBody)
	saveAmount, err := strconv.Atoi(httpRequestBodyStr)

	if err != nil {
		msg := "Request convertation failed: " + err.Error()
		fmt.Println(msg)
		w.Write([]byte(msg))
		return
	}

	// тут тоже самое что и функции выше
	// см. описание мьютекса и зачем он тут нужен
	mtx.Lock()
	if saveAmount > 0 && money >= saveAmount {
		bank += saveAmount
		fmt.Println("Now bank is: ", bank)
		money -= saveAmount
		fmt.Println("Now money is: ", money)
		_, err = w.Write([]byte("Succsessfuly save money"))

		msg := "Failed to write HTTP response " + err.Error()
		if err != nil {
			fmt.Println(msg)
			w.Write([]byte(msg))
		}
	}
	mtx.Unlock()
}

func main() {
	http.HandleFunc("/pay", payHandler)
	http.HandleFunc("/save", saveHandler)

	http.ListenAndServe(":8080", nil)
}
