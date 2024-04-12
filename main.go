package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/sum", processStuff)

	fmt.Println("starting api")
	go http.ListenAndServe("localhost:8080", nil)
	time.Sleep(time.Second)

	r, err := http.NewRequest(http.MethodPost, "http://localhost:8080/sum?reason=im-a-gopher", bytes.NewReader([]byte("{\"numbers\": [1,2,3,4]}")))
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		panic(err)
	}

	// remember Content-type accept/json
	result, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(result))
}

type Nums struct {
	Numbers []int `json:"numbers"`
}

func processStuff(w http.ResponseWriter, r *http.Request) {
	reason := r.URL.Query().Get("reason")

	var nums Nums
	err := json.NewDecoder(r.Body).Decode(&nums)
	if err != nil {
		fmt.Println(err)
	}

	in, out := make(chan int, 4), make(chan int)

	go sum(in, out)

	for i := 0; i < len(nums.Numbers); i++ {
		in <- nums.Numbers[i]
	}

	for i := 0; i < len(nums.Numbers)/2; i++ {
		result := <-out
		fmt.Printf("the sum result is %v\n", result)
	}

	w.Write([]byte(fmt.Sprintf("hire me because %v\n", reason)))
}

func sum(in chan int, out chan int) {
	for {
		a := <-in
		b := <-in
		fmt.Printf("summing %v with %v\n", a, b)

		out <- a + b
	}
}
