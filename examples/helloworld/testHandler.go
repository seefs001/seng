package main

import (
	"net/http"

	"github.com/seefs001/seng"
	"github.com/seefs001/seng/examples/helloworld/model"
)

type TestHandler struct {
	service model.DefaultApiServicer
}

func NewTestHandler(service model.DefaultApiServicer) *TestHandler {
	return &TestHandler{service: service}
}

func (t *TestHandler) AddPet(writer http.ResponseWriter, request *http.Request) {
	pets, err := t.service.AddPet(request.Context(), model.NewPet{
		Name: "1",
		Tag:  "1",
	})
	if err != nil {
		seng.RenderInternalServerError(writer, request, err)
	}
	seng.Render(writer, request, http.StatusOK, pets)
}

func (t TestHandler) DeletePet(writer http.ResponseWriter, request *http.Request) {
	panic("implement me")
}

func (t TestHandler) FindPetByID(writer http.ResponseWriter, request *http.Request) {
	panic("implement me")
}

func (t TestHandler) FindPets(writer http.ResponseWriter, request *http.Request) {
	panic("implement me")
}
