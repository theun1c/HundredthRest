package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/k0kubun/pp"
)

type Payment struct {
	Title    string  `json:"title"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Time     time.Time
}

type PaymentResponse struct {
	Money   float64   `json:"money"`
	History []Payment `json:"history"`
}

var money float64 = 1000
var history = make([]Payment, 0)
var mtx = sync.Mutex{}

func payHandler(w http.ResponseWriter, r *http.Request) {

	fooParam := r.URL.Query().Get("foo")
	booParam := r.URL.Query().Get("boo")

	fmt.Println("foo: ", fooParam)
	fmt.Println("boo: ", booParam)

	var payment Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		fmt.Println("Error while reading body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// нижнюю проверку и анмаршал можно заменить 1 строчкой сверху =)

	// httpRequestBody, err := io.ReadAll(r.Body)

	// if err != nil {
	// 	fmt.Println("Error while reading body: ", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	//
	// if err := json.Unmarshal(httpRequestBody, &payment); err != nil {
	// 	fmt.Println("Error while reading json: ", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	payment.Time = time.Now()

	mtx.Lock()
	if money-float64(payment.Amount) >= 0 {
		money -= float64(payment.Amount)
	} else {
		fmt.Println("Not enougth money, needed: ", payment.Amount, " money now: ", money)
	}

	history = append(history, payment)

	pp.Println(history)
	fmt.Println("money: ", money)

	httpResponse := PaymentResponse{
		Money:   money,
		History: history,
	}

	b, err := json.Marshal(httpResponse)

	if err != nil {
		fmt.Println("Error marshaling struct: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// записываем ответ
	if _, err := w.Write(b); err != nil {
		fmt.Println("Error while write response: ", err)
		return
	}

	mtx.Unlock()

}

func main() {
	http.HandleFunc("/pay", payHandler)

	// err := http.ListenAndServe(":8080", nil)

	// вариант 1
	// if err != nil {
	// 	fmt.Println("Error while starting server: ", err)
	// }

	// вариант 2 - более крутая запись
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error while starting server: ", err)
	}

}
