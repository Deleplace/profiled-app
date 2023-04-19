package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	// Register Pprof debug handlers
	_ "net/http/pprof"
)

func main() {
	http.HandleFunc("/", homepage)
	http.HandleFunc("/compute1", compute1)
	http.HandleFunc("/compute2", compute2)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	addr := os.Getenv("ADDR") + ":" + port
	log.Printf("Listening on %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	log.Fatal(err)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
		<a href="/compute1">/compute1</a> <br>
		<a href="/compute2">/compute2</a>
	`)
}

const divisor = 287654321

// compute1 is a handler that does CPU-intensive computations
func compute1(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting computation")
	const limit = 4_000_000_000
	start := time.Now()
	for i := 0; i < limit; i++ {
		n := rand.Int()
		if n%divisor == 0 {
			log.Printf("compute1: found %d\n", n)
			fmt.Fprintf(w, "%d is a multiple of %d", n, divisor)
			fmt.Fprintf(w, "\n\nComputation took %v", time.Since(start))
			return
		}
	}
	log.Printf("compute1: no multiple found\n")
	fmt.Fprintf(w, "Could not find a multiple of %d", divisor)
}

// compute2 is a handler that does CPU-intensive computations
func compute2(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting computation")
	const limit = 4_000_000_000
	start := time.Now()
	rng := rand.New(rand.NewSource(time.Now().UnixMicro()))
	for i := 0; i < limit; i++ {
		n := rng.Int()
		if n%divisor == 0 {
			log.Printf("compute2: found %d\n", n)
			fmt.Fprintf(w, "%d is a multiple of %d", n, divisor)
			fmt.Fprintf(w, "\n\nComputation took %v", time.Since(start))
			return
		}
	}
	log.Printf("compute2: no multiple found\n")
	fmt.Fprintf(w, "Could not find a multiple of %d", divisor)
}
