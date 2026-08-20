package main

import (
	"net/http"

	"github.com/cyntraten/metrics-ed/internal/handler"
	"github.com/cyntraten/metrics-ed/internal/repository"
)

func main() {

	savedMetrics := repository.NewMemStorage()
	h := handler.NewHandler(savedMetrics)

	mux := http.NewServeMux()

	mux.HandleFunc(`/update/`, h.UpdateMetrics)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}

}
