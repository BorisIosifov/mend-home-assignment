package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/BorisIosifov/mend-home-assignment/api"
	"github.com/BorisIosifov/mend-home-assignment/storage"
)

var storageString = flag.String("storage", "mysql", "Storage. Valid values are mysql, local_memory")

func main() {
	var (
		strg api.Storage
		err  error
	)

	flag.Parse()

	switch *storageString {
	case "mysql":
		strg, err = storage.PrepareMySQL()
	case "local_memory":
		strg, err = storage.PrepareLocalMemory()
	default:
		log.Fatalf("Unknown storage: %s", *storageString)
		return
	}
	if err != nil {
		log.Fatalf("Error while preparing a storage: %s", err)
	}

	api := api.API{
		Storage: strg,
	}

	err = http.ListenAndServeTLS(":443", "cert/cert.pem", "cert/key.pem", api)
	if err != nil {
		log.Fatalf("ListenAndServe: %s", err)
	}
}
