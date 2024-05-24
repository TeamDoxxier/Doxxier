package main

import (
	"encoding/json"
	"io"
	"net/http"

	"doxxier.tech/doxxier/compression"
	"doxxier.tech/doxxier/models"
	"doxxier.tech/doxxier/transformers"
)

type CompressionResponse struct {
	Length int
	Data   []byte
}

func main() {
	fileServer := http.FileServer(http.Dir("static"))

	http.Handle("/", fileServer)
	http.Handle("/compress", http.HandlerFunc(compress))
	http.Handle("/convert", http.HandlerFunc(convertImage))
	http.ListenAndServe(":8080", nil)
}

func convertImage(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("image")
	if err != nil {
		w.Write([]byte("Error parsing form"))
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		w.Write([]byte("Error reading file"))
		return
	}

	ctx := &models.DoxxierContext{ // Pass a pointer to ctx
		Content: content,
	}
	err = (&transformers.ImageTransformer{}).Transform(ctx) // Pass ctx as a pointer
	if err != nil {
		w.Write([]byte("Error transforming image"))
		return
	}
	w.Header().Set("content-type", "image/avif")
	w.WriteHeader(http.StatusOK)
	w.Write(ctx.Content)
}

func compress(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Error reading body"))
	}
	defer r.Body.Close()
	compressed := compression.Compress(body, compression.AlgXz)
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
