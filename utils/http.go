package utils

import (
	"encoding/json"
	"net/http"
)

type CommonResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func RespondSuccess(w http.ResponseWriter, r *http.Request, data interface{}) {
	Response(w, r, 0, "success", data)
}

func Response(w http.ResponseWriter, r *http.Request, status int, message string, data interface{}) {
	resp := CommonResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}

	body, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "marshal resp err", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")

	w.Write(body)
}
