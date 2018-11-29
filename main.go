package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

type Document struct {
	Id   string
	Name string
	Size int64
}

const rootfiles string = "C:/Test/"

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/documents", getDocuments).Methods("GET")
	router.HandleFunc("/documents/{id}", getDocumentsId).Methods("GET")
	router.HandleFunc("/documents", createDocuments).Methods("POST")
	router.HandleFunc("/documents/{id}", deleteDocumentsId).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":9000", router))
}

func GetInfoFile(filepath string) Document {

	var stringmd5 string
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	info, err := file.Stat()

	hash := md5.New()
	_, err = io.Copy(hash, file)
	hashytes := hash.Sum(nil)[:16]
	stringmd5 = hex.EncodeToString(hashytes)
	return Document{Id: stringmd5, Name: info.Name(), Size: info.Size()}
}

func DeleteFile(idFile string) bool {

	var docs = GetFiles()
	for _, num := range docs {
		if num.Id == idFile {
			return DeteleArchive(num.Name)
		}
	}
	return false
}

func DeteleArchive(filename string) bool {

	var path = rootfiles + filename
	var err = os.Remove(path)
	if err != nil {
		log.Fatal(err)
	}
	return true

}
func GetFile(idfile string) Document {

	var listDocuments = GetFiles()
	documentResult := Document{}
	for _, documento := range listDocuments {
		if documento.Id == idfile {
			documentResult = documento
		}
	}
	return documentResult
}

func GetFiles() []Document {

	var docs []Document
	err := filepath.Walk(rootfiles, func(path string, info os.FileInfo, err error) error {
		if path != rootfiles {
			docs = append(docs, GetInfoFile(path))
			fmt.Println(path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return docs

}

func getDocuments(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetFiles())
}

func getDocumentsId(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetFile(vars["id"]))
}

func createDocuments(w http.ResponseWriter, r *http.Request) {

	fileType := r.PostFormValue("type")
	fmt.Println(fileType)
	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Error! you should post a file with the title 'file'"))
		return
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Error! error not filebytes"))
		return
	}

	uploadPath := rootfiles
	newPath := filepath.Join(uploadPath, header.Filename)
	newFile, err := os.Create(newPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Error! could not create path"))
		return
	}
	defer newFile.Close()
	if _, err := newFile.Write(fileBytes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Error! new file error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File created"))

}

func deleteDocumentsId(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	if DeleteFile(vars["id"]) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
