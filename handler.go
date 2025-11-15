package main

import (
	"encoding/json"
	"location-service/service"
	"net/http"
)

type UpdateReq struct {
	DriverID string  `json:"driver_id"`
	Lat      float64 `json:"lat"`
	Long     float64 `json:"long"`
}

type SearchReq struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type Handler struct {
	geo *service.GeoService
}

func NewHandler(geo *service.GeoService) *Handler {
	return &Handler{geo: geo}
}

func (handler *Handler) UpdateDriverLocationHandler(w http.ResponseWriter, r *http.Request) {
	var req UpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := handler.geo.UpdateLocation(r.Context(), req.DriverID, req.Lat, req.Long); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`${"message": "location updated"}`))
}

func (handler *Handler) FindNearestDriveHandler(w http.ResponseWriter, r *http.Request) {
	var req SearchReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := handler.geo.FindNearestDrive(r.Context(), req.Lat, req.Long)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map redis.GeoLocation to a smaller response DTO
	type NearestDriverResp struct {
		DriverID string  `json:"driver_id"`
		Lat      float64 `json:"lat"`
		Long     float64 `json:"long"`
		DistKM   float64 `json:"dist_km"`
	}

	resp := make([]NearestDriverResp, 0, len(result))
	for _, r := range result {
		resp = append(resp, NearestDriverResp{
			DriverID: r.Name,
			Lat:      r.Latitude,
			Long:     r.Longitude,
			DistKM:   r.Dist,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
