package handler

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/sevigo/hokan/pkg/core"
)

func JSON_200(w http.ResponseWriter, r *http.Request, v interface{}) {
	render.Status(r, 200)
	render.JSON(w, r, v)
}

func JSON_201(w http.ResponseWriter, r *http.Request, msg string) {
	render.Status(r, 201)
	render.JSON(w, r, core.Resp{Code: http.StatusCreated, Msg: msg, Status: "success"})
}

func JSON_400(w http.ResponseWriter, r *http.Request, msg string) {
	render.Status(r, 400)
	render.JSON(w, r, core.Resp{Code: http.StatusBadRequest, Msg: msg, Status: "error"})
}

func JSON_404(w http.ResponseWriter, r *http.Request, msg string) {
	render.Status(r, 404)
	render.JSON(w, r, core.Resp{Code: http.StatusNotFound, Msg: msg, Status: "error"})
}

func JSON_500(w http.ResponseWriter, r *http.Request, msg string) {
	render.Status(r, 500)
	render.JSON(w, r, core.Resp{Code: http.StatusInternalServerError, Msg: msg, Status: "error"})
}
