package main

import (
	"log"
	// "net"
	"net/http"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	http.Handle("/dashboard/", http.StripPrefix("/dashboard/", http.FileServer(http.Dir("./dashboard/"))))
	log.Fatal(http.ListenAndServe(":1233", nil))

}
