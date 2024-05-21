package main

import (
	"encoding/json"
	"io"
	"net/http"

	"doxxier.tech/doxxier/compression"
)

type CompressionResponse struct {
	Length int
	Data   []byte
}

func main() {
	fileServer := http.FileServer(http.Dir("static"))

	http.Handle("/", fileServer)
	http.Handle("/compress", http.HandlerFunc(compress))
	http.ListenAndServe(":8080", nil)
}

func compress(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Error reading body"))
	}
	defer r.Body.Close()
	compressed := compression.Compress(body)
	response := CompressionResponse{
		Length: len(compressed),
		Data:   compressed,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
