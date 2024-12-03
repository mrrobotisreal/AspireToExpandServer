package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func HandleUploadProfileImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // This limits images to 10MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	savePath := filepath.Join("uploads", "profileImages", handler.Filename)
	dst, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	// defer dst.Close() // <- close doesn't exist?

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Unable to copy/save file", http.StatusInternalServerError)
		return
	}

	imageURL := fmt.Sprintf("https://%s:8888/uploads/profileImages/%s", os.Getenv("IP_ADDRESS"), handler.Filename)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"imageURL": "%s"}`, imageURL)
}
