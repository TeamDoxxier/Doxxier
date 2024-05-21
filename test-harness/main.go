package main

import (
	"net/http"
)

func main() {
	fileServer := http.FileServer(http.Dir("static"))

	http.Handle("/", fileServer)
	http.Handle("/compress", http.HandlerFunc(compress))
	http.ListenAndServe(":8080", nil)
}

func compress(w http.ResponseWriter, r *http.Request) {
	compressed := doxxier.compress(r.Body)
	if len(compressed) == 0 {
		t.Error("Compressed data is empty")
	}
	w.Write(compressed.len)
}
