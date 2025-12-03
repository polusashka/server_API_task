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
		http.Error(res, "Файл index.html не найден", http.StatusOK)
		return
	}

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write(data)
}

func UploadHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.Error(res, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.Error(res, "Ошибка парсинга формы", http.StatusInternalServerError)
		return
	}

	file, header, err := req.FormFile("myFile")
	if err != nil {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.Error(res, "Файл не найден", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.Error(res, "Ошибка чтения файла", http.StatusInternalServerError)
		return
	}

	processedData, err := service.AutoDetection(string(fileData))
	if err != nil {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.Error(res, "Ошибка перобразования "+err.Error(), http.StatusInternalServerError)
		return
	}

	uploadDir := "../upload"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, 0755)
		if err != nil {
			res.Header().Set("Content-Type", "text/html; charset=utf-8")
			http.Error(res, "Ошибка создания директории: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	filename := filepath.Join(uploadDir, header.Filename)
	//timeStamp := time.Now().Format("2006-01-02_15-04-05")
	//newFile, err := os.Create(filename + fmt.Sprintf("_%s.txt", timeStamp))
	errCreate := os.WriteFile(filename, []byte(processedData), 0755)
	if errCreate != nil {
		log.Printf("Ошибка сохранения файла: %v", err)
	}

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(processedData))
}
