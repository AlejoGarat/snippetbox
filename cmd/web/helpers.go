package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

type Helper struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func (h *Helper) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	h.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (h *Helper) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (h *Helper) notFound(w http.ResponseWriter) {
	h.clientError(w, http.StatusNotFound)
}
