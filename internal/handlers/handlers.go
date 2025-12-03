package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

func MainHandler(res http.ResponseWriter, req *http.Request) {
	data, err := os.ReadFile("../index.html")
	if err != nil {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Файл index.html не найден"))
		return
	}

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write(data)
}

func UploadHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		res.Write([]byte("Метод не поддерживается"))
		return
	}

	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Ошибка парсинга формы"))
		return
	}

	file, header, err := req.FormFile("myFile")
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Файл не найден"))
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Ошибка чтения файла"))
		return
	}

	processedData, err := service.AutoDetection(string(fileData))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Ошибка преобразования " + err.Error()))
		return
	}

	uploadDir := "../upload"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, 0755)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("Ошибка создания директории: " + err.Error()))
			return
		}
	}

	filename := filepath.Join(uploadDir, header.Filename)
	errCreate := os.WriteFile(filename, []byte(processedData), 0755)
	if errCreate != nil {
		log.Printf("Ошибка сохранения файла: %v", err)
	}

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(processedData))
}
