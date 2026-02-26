package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	for c := range time.NewTicker(time.Second).C {
		fmt.Println(c.Format(time.RFC3339))
	}
	http.ListenAndServe(":8080", nil)
}
