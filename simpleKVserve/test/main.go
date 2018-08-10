package main

import "net/http"

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte{byte('h')})
	})
	http.ListenAndServe(":1231",nil)
}
