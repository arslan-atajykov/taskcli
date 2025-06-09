package main

import (
	"log"
	"net/http"
)

func main() {
	r := NewRouter()

	log.Println("server is running on localhost:8888")
	err := http.ListenAndServe(":8888", r)
	if err != nil {
		log.Fatal("Server failed : ", err)
	}
}
