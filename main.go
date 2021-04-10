package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	go func() {
		client := &http.Client{Timeout: 1 * time.Second}

		for {
			time.Sleep(5 * time.Second)

			ids := make([]string, rng.Int31n(200))
			for i := range ids {
				ids[i] = strconv.Itoa(rng.Int() % 100)
			}
			body := bytes.NewBuffer([]byte(fmt.Sprintf(`{"object_ids":[%s]}`, strings.Join(ids, ","))))
			resp, err := client.Post("http://localhost:9090/callback", "application/json", body)

			if err != nil {
				fmt.Println(err)
				continue
			}
			_ = resp.Body.Close()
		}
	}()

	http.HandleFunc("/callback", postRoute)

	go func() { _ = http.ListenAndServe(":9090", nil) }()

	http.HandleFunc("/objects/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(rng.Int63n(4000)+300) * time.Millisecond)

		idRaw := strings.TrimPrefix(r.URL.Path, "/object/")
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		result := []byte(fmt.Sprintf(`{"id":%d,"online":%v}`, id, id%2 == 0))
		fmt.Print("result")
		w.Write(result)
	})
	go func() { _ = http.ListenAndServe(":9010", nil) }()

	//http.HandleFunc("/objects/{id}", otherRoute)
	//
	//go func() { _ =  http.ListenAndServe(":9010", nil) } ()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig

	fmt.Println("closing")
}

type MyResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (mrw *MyResponseWriter) Write(p []byte) (int, error) {
	return mrw.buf.Write(p)
}

func loginmw(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		// Work / inspect body. You may even modify it!

		// And now set a new body, which will simulate the same data we read:
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		// Create a response wrapper:
		mrw := &MyResponseWriter{
			ResponseWriter: w,
			buf:            &bytes.Buffer{},
		}

		// Call next handler, passing the response wrapper:
		handler.ServeHTTP(mrw, r)

		// Now inspect response, and finally send it out:
		// (You can also modify it before sending it out!)
		if _, err := io.Copy(w, mrw.buf); err != nil {
			log.Printf("Failed to send out response: %v", err)
		}
	})
}

func postRoute(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("response Body:", string(body))

	// further handle request
}

func otherRoute(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("other Body:", string(body))
}
