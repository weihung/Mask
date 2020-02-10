package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
)

func main() {
	defer CloseDB()

	if v := os.Getenv("PORT"); len(v) > 0 {
		port = v
	}

	if v := os.Getenv("API_PATH"); len(v) > 0 {
		path = "/" + v
	}

	log.Println("Mask server Start")
	for i := range os.Args {
		arg := os.Args[i]
		if arg == "-s" {
			file := os.Args[i+1]
			UpdateStore(file)
		}
	}

	go FetchNewData()

	mux := http.NewServeMux()
	mux.HandleFunc(path+"/Mask", ApiMask)
	mux.HandleFunc(path+"/CardStatus", ApiCardStatus)
	log.Println("Mask API server is running")
	handler := cors.Default().Handler(mux)
	log.Println(http.ListenAndServe(":"+port, handler))

	log.Println("Server shut down.")
}
