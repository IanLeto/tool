package main

import "net/http"

func main() {
	http.HandlerFunc("/upload", func(writer http.ResponseWriter, request *http.Request) {
		// Parse our multipart form, 10 << 20 specifies a maximum
		// upload of 10 MB files.
		err := request.ParseMultipartForm(10 << 20)
		if err != nil {
			return
		}
		file, handler, err := request.FormFile("")

	})
	http.ListenAndServe(":8008", nil)
}
