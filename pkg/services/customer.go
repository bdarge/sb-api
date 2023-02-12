package services

import (
	"context"
	"github.com/bdarge/api/out/customer"
	"github.com/bdarge/api/pkg/db"
	"github.com/bdarge/api/pkg/models"
	"github.com/bdarge/api/pkg/util"
	"log"
	"net/http"
)

type CustomerServer struct {
	H db.Handler
	customer.UnimplementedCustomerServiceServer
}

func (server *CustomerServer) CreateCustomer(ctx context.Context, request *customer.CreateCustomerRequest) (*customer.CreateCustomerResponse, error) {
	log.Printf("create customer: %s\n", request.Name)

	c := &models.Customer{}
	err := util.Recast(request, c)
	err = server.H.DB.Create(&c).
		Error

	if err != nil {
		return &customer.CreateCustomerResponse{Status: http.StatusBadRequest, Error: err.Error()}, nil
	}

	return &customer.CreateCustomerResponse{
		Status: http.StatusCreated,
		Id:     c.ID,
	}, nil
}

func (server *CustomerServer) GetCustomer(ctx context.Context, request *customer.GetCustomerRequest) (*customer.GetCustomerResponse, error) {
	var c models.Customer
	log.Printf("get customer with id, %d\n", request.Id)

	err := server.H.DB.Where("id = ?", request.Id).
		First(&c).
		Error

	if err != nil {
		return &customer.GetCustomerResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	data := &customer.CustomerData{}
	err = util.Recast(c, data)
	return &customer.GetCustomerResponse{
		Status: http.StatusOK,
		Data:   data,
	}, nil
}
