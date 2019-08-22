package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"time"

	// Add trace module
	"./trace"
)

func CallStepSimple(w http.ResponseWriter, r *http.Request) {

	log.Println("Executing downstream request")


	// Add color
	trace.TraceFunctionExecution(r, func() {
		time.Sleep(1 * time.Second)
	}, "Going to sleep")

	httpTrace := trace.TraceHTTPClientGet(r, "http://localhost:8082/", "Downstream call")

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
		fmt.Fprintf(w, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	httpTrace.Finish()
}

func CallStepComplicated(w http.ResponseWriter, r *http.Request) {

	// Need control over the over HTTP headers
	client := http.Client{Timeout: 5 * time.Second}

	// HTTP request
	downstreamRequest, _ := http.NewRequest("GET", "http://localhost:8082/", nil)
	resp, err := trace.GoSensor.TracingHttpRequest("exit", r, downstreamRequest, client)
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
		fmt.Fprintf(w, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

}

func main() {
	log.Println("Server started on: http://localhost:8081")
	trace.SetupTracer("step1")
	log.Printf("Using tracer %T", trace.GoSensor)
	// Added Tracing Handler
	http.HandleFunc("/", trace.GoSensor.TracingHandler("step1", CallStepComplicated))

	http.ListenAndServe(":8081", nil)
}


