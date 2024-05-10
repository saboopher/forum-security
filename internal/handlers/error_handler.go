package handler

import (
	"fmt"
	"forum/internal/exceptions"
	"forum/pkg/cust_encoders"
	"html/template"
	"net/http"
)

func errorHandler(w http.ResponseWriter, r *http.Request) {
	errObj := r.Context().Value("error")
	httpStatus := r.Context().Value("httpStatus")
	fmt.Printf("Handler %v\n", errObj)
	if errObj == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_, ok := errObj.(error)
	if !ok {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	t, parseErr := template.ParseFiles("ui/templates/error.html")
	if parseErr != nil {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
	if httpStatus == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	statusCode, okInt := httpStatus.(int)
	if !okInt {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
	w.WriteHeader(statusCode)
	err := t.Execute(w, errObj)
	if err != nil {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
	}
}
