package handler

import (
	"net/http"

	"github.com/go-chi/render"
)

func JSON_200(w http.ResponseWriter, r *http.Request, v interface{}) {
	render.Status(r, 200)
	render.JSON(w, r, v)
}

func JSON_201(w http.ResponseWriter, r *http.Request) {
	render.Status(r, 201)
	render.JSON(w, r, nil)
}

func JSON_400(w http.ResponseWriter, r *http.Request, v interface{}) {
	render.Status(r, 400)
	render.JSON(w, r, v)
}

func JSON_404(w http.ResponseWriter, r *http.Request, v interface{}) {
	render.Status(r, 404)
	render.JSON(w, r, v)
}
