package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime/pprof"
	"strconv"
	"time"

	// Register Pprof debug handlers
	_ "net/http/pprof"
)

func main() {
	http.HandleFunc("/", homepage)
	http.HandleFunc("/compute1", compute1)
	http.HandleFunc("/compute2", compute2)
	http.HandleFunc("/customprofile", customprofile)

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

// Custom profiling handler, copied from
// https://cs.opensource.google/go/go/+/refs/tags/go1.20.3:src/net/http/pprof/pprof.go;l=125
func customprofile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	sec, err := strconv.ParseInt(r.FormValue("seconds"), 10, 64)
	if sec <= 0 || err != nil {
		sec = 30
	}

	if sec > 3600 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "profile duration %ds is too large", sec)
		return
	}

	// Set Content Type assuming StartCPUProfile will work,
	// because if it does it starts writing.
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", `attachment; filename="profile"`)
	if err := pprof.StartCPUProfile(w); err != nil {
		// StartCPUProfile failed, so no writes yet.
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not enable CPU profiling: %s", err)
		return
	}
	sleep(r, time.Duration(sec)*time.Second)
	pprof.StopCPUProfile()
}

func sleep(r *http.Request, d time.Duration) {
	select {
	case <-time.After(d):
	case <-r.Context().Done():
	}
}
