package main

import (
	"fmt"
	"net/http"
)

// 这里的package main 和 helloworld下的package main是不一样的
func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "<h1>Hello World! %s", request.FormValue("name"))
	})
	http.ListenAndServe(":8888", nil)
}
