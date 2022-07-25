package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func SendOptions(opts string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Allow", opts)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", opts)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Cookie, Authorization, X-Auth-Token")
		w.WriteHeader(http.StatusOK)
	})
}

func SendNotFound(w http.ResponseWriter) {
	w.Header().Set("X-Result-Type", "error")
	data := struct {
		code    int
		message string
	}{http.StatusNotFound, "Not found"}
	SendJson(w, data, http.StatusNotFound)
}

func SendInternalError(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf(ErrInternal, r.Context().Value(KeyRequestId).(string))
	SendJsonError(w, http.StatusInternalServerError, msg)
}

func SendJsonError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("X-Result-Type", "error")
	data := struct {
		Code    int
		Message string
	}{status, msg}
	SendJson(w, data, status)
}

func SendJson(w http.ResponseWriter, data interface{}, code ...int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if len(code) > 0 {
		w.WriteHeader(code[0])
	} else {
		w.WriteHeader(http.StatusOK)
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("write JSON error : %v", err.Error())
	}
}
