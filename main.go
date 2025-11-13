package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// с помощью респонс врайтера можно писать сетевые ответы на запросы
	// принимает []byte
	// w.Write()

	str := "hello hello hello"
	// преобразуем строку в слайс байтов
	bytes := []byte(str)
	// передаем и получаем ошибку
	_, err := w.Write(bytes)

	if err != nil {
		fmt.Println("Error", err)
	}
}

// используем атомик на переменной,
// тк на эту переменную возможна гонка данных
// а атомики быстрее мьютексов
var money = atomic.Int64{} // денег в кармане

// с банком тоже самое
var bank = atomic.Int64{} // денег в копилке

// создаем мьютекс для критической секции
// в которой 2 разных горутины (хендлера)
// работают с одним атомиком маней.
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
		fmt.Println("Failed to read HTTP boby", err)
		return
	}

	httpRequestBobyStr := string(httpRequestBoby)
	amount, err := strconv.Atoi(httpRequestBobyStr)

	if err != nil {
		fmt.Println("strconv error", err)
		return
	}

	// тот самый блок с возможной гонкой данных на переменной маней
	// обращение из 2х независимых горутин к 1 переменной
	mtx.Lock()
	if money.Load()-int64(amount) >= 0 {
		// у атомика нет метода минус, поэтому делаем через добавление -числа
		money.Add(int64(-amount))
		fmt.Println("Succsessfuly payment, now money is: ", money.Load())
	} else {
		fmt.Println("No money to pay. monay amount is: ", money.Load())
	}
	mtx.Unlock()
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	httpRequestBody, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Println("Error while reading body", err)
		return
	}

	// преобразуем слайс байт в строку
	httpRequestBodyStr := string(httpRequestBody)
	saveAmount, err := strconv.Atoi(httpRequestBodyStr)

	if err != nil {
		fmt.Println("Failed to convert", err)
		return
	}

	// тут тоже самое что и функции выше
	// см. описание мьютекса и зачем он тут нужен
	mtx.Lock()
	if saveAmount > 0 && money.Load() >= int64(saveAmount) {
		bank.Add(int64(saveAmount))
		fmt.Println("Now bank is: ", bank.Load())
		money.Add(int64(-saveAmount))
		fmt.Println("Now money is: ", money.Load())
	} else {
		fmt.Println("cannot add this amount to bank")
	}
	mtx.Unlock()
}

func main() {

	// базовая сигнатура хендл функции
	//
	//http.HandleFunc("/", handler)

	// кладем значение в переменную money,
	// поскольку при объявлении атомика
	// в него ничего нельзя класть
	money.Add(1000)

	// не нужно добавлять 0 =)
	// bank.Add(0)

	http.HandleFunc("/pay", payHandler)
	http.HandleFunc("/save", saveHandler)

	http.ListenAndServe(":8080", nil)

}
