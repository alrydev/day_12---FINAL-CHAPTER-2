package middleware

import (
	"context"
	"encoding/json" //untuk menconvert ke bentuk json
	"fmt"
	"io/ioutil"
	"net/http"
)

func UploadFile(next http.HandlerFunc) http.HandlerFunc {
	// value didapat dari form add projectnya
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, handler, err := r.FormFile("inputImage")
		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode("Error Retrieving the File")
			return
		}
		defer file.Close()

		// untuk message jika yg di upload sudah sukses:
		fmt.Printf("Uploaded File: %+v\n", handler.Filename) // Uploaded File: image.png (nama filenya)

		//untuk melakukan handling nama file image (mengubah nama file yg temporarynya disimpan di upload folder):
		tempFile, err := ioutil.TempFile("uploads", "image-*"+handler.Filename) // contoh image-profil (profil = handler.Filename)
		if err != nil {                                                         // ketika tidak berhasil menambahkan file
			fmt.Println(err)
			fmt.Println("path upload error.")
			json.NewEncoder(w).Encode(err)
			return
		}
		defer tempFile.Close() //jika terjadi apapun, maka akan next semua bagian diatas ke komen

		// untuk baca file:
		fileBytes, err := ioutil.ReadAll(file) // dlm ioutil kita akan membaca semua file yg kita upload di front endnya
		if err != nil {
			fmt.Println(err)
		}

		// Create image temporary file jika ada filemya
		tempFile.Write(fileBytes)

		data := tempFile.Name() // data membawa value tempFile; adalah nama temporary dari folder uploads
		filename := data[8:]    // uploads/[image.jpg] ; image.jpg yg akan disimpan di db, bagian "uploads/" di slice

		// filename diubah jadi dataFile yg akan digunakan ctx nya di main.go
		ctx := context.WithValue(r.Context(), "dataFile", filename)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UploadFile2(next http.HandlerFunc) http.HandlerFunc {
	// value didapat dari form add projectnya
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, handler, err := r.FormFile("inputImage")
		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode("Error Retrieving the File")
			return
		}
		defer file.Close()

		// untuk message jika yg di upload sudah sukses:
		fmt.Printf("Uploaded File: %+v\n", handler.Filename) // Uploaded File: image.png (nama filenya)

		//untuk melakukan handling nama file image (mengubah nama file yg temporarynya disimpan di upload folder):
		tempFile, err := ioutil.TempFile("uploads", "image-*"+handler.Filename) // contoh image-profil (profil = handler.Filename)
		if err != nil {                                                         // ketika tidak berhasil menambahkan file
			fmt.Println(err)
			fmt.Println("path upload error.")
			json.NewEncoder(w).Encode(err)
			return
		}
		defer tempFile.Close() //jika terjadi apapun, maka akan next semua bagian diatas ke komen

		// untuk baca file:
		fileBytes, err := ioutil.ReadAll(file) // dlm ioutil kita akan membaca semua file yg kita upload di front endnya
		if err != nil {
			fmt.Println(err)
		}

		// Create image temporary file jika ada filemya
		tempFile.Write(fileBytes)

		data := tempFile.Name() // data membawa value tempFile; adalah nama temporary dari folder uploads
		filename := data[8:]    // uploads/[image.jpg] ; image.jpg yg akan disimpan di db, bagian "uploads/" di slice

		// filename diubah jadi dataFile yg akan digunakan ctx nya di main.go
		ctx := context.WithValue(r.Context(), "dataFile", filename)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
