package service

import (
	"context"

	"github.com/seefs001/seng/examples/helloworld/model"
)

type TestService struct {
}

func NewTestService() *TestService {
	return &TestService{}
}

func (t TestService) AddPet(ctx context.Context, pet model.NewPet) (model.ImplResponse, error) {
	return model.ImplResponse{
		Code: 1,
		Body: 100,
	}, nil
}

func (t TestService) DeletePet(ctx context.Context, i int64) (model.ImplResponse, error) {
	panic("implement me")
}

func (t TestService) FindPetByID(ctx context.Context, i int64) (model.ImplResponse, error) {
	panic("implement me")
}

func (t TestService) FindPets(ctx context.Context, strings []string, i int32) (model.ImplResponse, error) {
	panic("implement me")
}
