package main

import (
	"fmt"
	"location-service/service"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	fmt.Println("Start of project")
	rdb := NewRedisClient()
	geo := &service.GeoService{Rdb: rdb}
	h := NewHandler(geo)
	r := chi.NewRouter()
	r.Patch("/update", h.UpdateDriverLocationHandler)
	r.Post("/find", h.FindNearestDriveHandler)
	log.Println("Server started on port 8080")
	http.ListenAndServe(":8080", r)
}
