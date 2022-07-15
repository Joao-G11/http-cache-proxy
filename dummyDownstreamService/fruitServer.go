package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
)

func getFruitRequest(w http.ResponseWriter, r *http.Request) {

	fruitIdx := rand.Intn(3)
	w.Write([]byte(fruits[fruitIdx]))

}

var fruits = []string{"Pear", "Banana", "Apple"}

func main() {

	if len(os.Args) != 3 {
		fmt.Println("usage: go run <address> <port>")
		os.Exit(1)
	}

	http.HandleFunc("/", getFruitRequest)

	err := http.ListenAndServe(fmt.Sprintf("%s:%s", os.Args[1], os.Args[2]), nil)
	if err != nil {
		fmt.Println(err)
	}

}
