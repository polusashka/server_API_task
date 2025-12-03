package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

func MainHandler(res http.ResponseWriter, req *http.Request) {
	data, err := os.ReadFile("../index.html")
	if err != nil {
		http.Error(res, "Файл index.html не найден", http.StatusOK)
		return
	}

	res.Header().Set("Content-type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write(data)
}

func UploadHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(res, "Ошибка парсинга формы", http.StatusInternalServerError)
		return
	}

	file, header, err := req.FormFile("myFile")
	if err != nil {
		http.Error(res, "Файл не найден", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(res, "Ошибка чтения файла", http.StatusInternalServerError)
		return
	}

	processedData, err := service.AutoDetection(string(fileData))
	if err != nil {
		http.Error(res, "Ошибка перобразования"+err.Error(), http.StatusInternalServerError)
		return
	}
	filename := filepath.Join("../upload", header.Filename)
	//timeStamp := time.Now().Format("2006-01-02_15-04-05")
	//newFile, err := os.Create(filename + fmt.Sprintf("_%s.txt", timeStamp))
	errCreate := os.WriteFile(filename, []byte(processedData), 0755)
	if errCreate != nil {
		http.Error(res, "Ошибка создания локального файла"+errCreate.Error(), http.StatusInternalServerError)
		return
	}
	//defer newFile.Close()

	//os.WriteFile(filename, []byte(processedData), 0755)

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(processedData))
}
