package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Document struct {
	Id   string
	Name string
	Size int64
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/documents", getDocuments).Methods("GET")
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

func GetFiles(Path string) []Document {

	files, err := ioutil.ReadDir(Path)
	if err != nil {
		log.Fatal(err)
	}

	var docs []Document
	for _, f := range files {

		docs = append(docs, GetInfoFile(Path+f.Name()))

	}
	return docs
}

func getDocuments(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetFiles("C:/Test/"))
}
