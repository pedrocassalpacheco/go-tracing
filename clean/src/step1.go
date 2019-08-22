package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
)

func CallStep2(writer http.ResponseWriter, request *http.Request) {

	log.Println("Executing downstream request")
	resp, err := http.Get("http://localhost:8082/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Print the HTTP response status.
	fmt.Println("Response status:", resp.Status)

	// Print the first 5 lines of the response body.
	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan(); i++ {
		fmt.Println(scanner.Text())
		fmt.Fprintf(writer, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func main() {
	log.Println("Server started on: http://localhost:8081")
	http.HandleFunc("/", CallStep2)
	http.ListenAndServe(":8081", nil)
}


