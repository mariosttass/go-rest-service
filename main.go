package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/mariosttass/go-rest-service/pkg/repository"
)

const (
	dbHost = "DBHOST"
	dbPort = "DBPORT"
	dbUser = "DBUSER"
	dbPass = "DBPASS"
	dbName = "DBNAME"
)

func main() {
	//open connection to DB
	connString := repository.ConnString(dbUser, dbPass, dbHost, dbName, dbPort)
	db, errDB := repository.NewDbConnector(connString)
	if errDB != nil {
		log.Fatal("could not connnect to db", errDB)
	}
	defer db.Stop()

	//passing db instance on the handlers.
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

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig

	fmt.Println("closing")
}

func postRoute(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("response Body:", string(body))

}
