package main

import (
	"github.com/vinicius-maker/observabilidade-open-telemetry/internal/controller"
	"net/http"
)

func main() {
	http.HandleFunc("/", controller.Handle)

	http.ListenAndServe(":8081", nil)
}
