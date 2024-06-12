package main

import (
	"github.com/MaximSvinin/tape/internal/controller/http"
	"github.com/MaximSvinin/tape/pkg/tape"
)

func main() {
	http.Router(tape.NewTapeStorage())
}
