package main

import (
	"fmt"
	"net/http"
	"tiny-url/pkg/routes"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	//TODO: handle CORS here
	routes.CreateRoutes(r.PathPrefix("/api").Subrouter())
	fmt.Println("Running server on port 3000")
	err := http.ListenAndServe("localhost:3000", r)
	if err != nil {
		panic(err)
	}
}
