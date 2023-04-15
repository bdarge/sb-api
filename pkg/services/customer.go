package services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/bdarge/api/out/customer"
	"github.com/bdarge/api/out/model"
	"github.com/bdarge/api/pkg/db"
	"github.com/bdarge/api/pkg/models"
	"github.com/bdarge/api/pkg/util"
	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
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
		return &customer.GetCustomerResponse{
			Status: http.StatusInternalServerError, Error: err.Error(),
		}, nil
	}

	response, err := mapCustomer(&c)

	return &customer.GetCustomerResponse{
		Status: http.StatusOK,
		Data:   response,
	}, nil
}

func (server *CustomerServer) GetCustomers(_ context.Context, request *customer.GetCustomersRequest) (*customer.GetCustomersResponse, error) {
	log.Printf("get customers, search filter=%s", request.Search)

	var customers = make([]models.Customer, 0)

	if request.Limit == 0 {
		request.Limit = 10
	}

	if request.SortDirection == "" {
		request.SortDirection = "desc"
	} else {
		request.SortDirection = strings.ToLower(request.SortDirection)
	}

	if request.SortProperty == "" {
		request.SortProperty = "customers.id"
	} else {
		request.SortProperty = "customers." + util.ToSnakeCase(request.SortProperty)
	}

	log.Printf("request: %v", request)

	err := server.H.DB.Model(&models.Customer{}).
		Where("true = ? Or Name LIKE ?", request.Search == "", "%"+request.Search+"%").
		Where("true = ? Or Email LIKE ?", request.Search == "", "%"+request.Search+"%").
		Limit(int(request.Limit)).
		Offset(int(request.Page * request.Limit)).
		Order(request.SortProperty + " " + request.SortDirection).
		Find(&customers).
		Error

	if err != nil {
		return &customer.GetCustomersResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	var total int64

	server.H.DB.Model(&models.Customer{}).
		Count(&total)

	page := models.Customers{
		Data:  customers,
		Limit: request.Limit,
		Page:  request.Page,
		Total: uint32(total),
	}

	if err != nil {
		return &customer.GetCustomersResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	messageInBytes, err := json.Marshal(page)
	if err != nil {
		return &customer.GetCustomersResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}
	log.Printf("customers: %s", messageInBytes)
	var response customer.GetCustomersResponse

	// ignore unknown fields
	unMarshaller := &protojson.UnmarshalOptions{DiscardUnknown: true}
	err = unMarshaller.Unmarshal(messageInBytes, &response)
	if err != nil {
		log.Printf("Error: %v", err)
		return &customer.GetCustomersResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}
	log.Printf("message:- %v", response)
	response.Status = http.StatusOK
	return &response, nil
}

func (server *CustomerServer) UpdateCustomer(_ context.Context, request *customer.UpdateCustomerRequest) (*customer.UpdateCustomerResponse, error) {
	var c models.Customer
	log.Printf("update customer (id = %d)\n", request.Id)

	err := server.H.DB.Model(&models.Customer{}).
		Where("customers.id = ?", request.Id).
		First(&c).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &customer.UpdateCustomerResponse{Status: http.StatusNotFound, Error: err.Error()}, nil
		}
		return &customer.UpdateCustomerResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	value, err := protojson.Marshal(request)
	if err != err {
		return &customer.UpdateCustomerResponse{Status: http.StatusBadRequest, Error: err.Error()}, nil
	}
	log.Printf("update customer constructed from a proto message: %s", value)

	u := &models.Customer{}

	err = json.Unmarshal(value, u)

	log.Printf("update customer request constructed from bytes: %v", u)

	if err != err {
		return &customer.UpdateCustomerResponse{Status: http.StatusBadRequest, Error: err.Error()}, nil
	}

	update := make(map[string]interface{})

	if u.Name != nil && u.Name != c.Name {
		update["name"] = u.Name
	}
	if u.Email != nil && u.Email != c.Email {
		update["email"] = u.Email
	}

	err = server.H.DB.Model(c).Updates(update).Error

	if err != nil {
		log.Printf("Failed to update: %v", err)
		return &customer.UpdateCustomerResponse{
			Status: http.StatusBadRequest, Error: "Failed to update",
		}, nil
	}

	err = server.H.DB.Model(&models.Customer{}).
		Where("customers.id = ?", request.Id).
		First(&c).
		Error

	if err != nil {
		log.Printf("Failed to get updated data: %v", err)
		return &customer.UpdateCustomerResponse{
			Status: http.StatusBadRequest, Error: "Failed to get updated data",
		}, nil
	}

	response, err := mapCustomer(&c)

	if err != nil {
		log.Printf("Mapping to map to proto type failed: %v", err)
		return &customer.UpdateCustomerResponse{
			Status: http.StatusBadRequest, Error: "Failed to map to proto type",
		}, nil
	}

	return &customer.UpdateCustomerResponse{
		Status: http.StatusOK,
		Data:   response,
	}, nil
}

func (server *CustomerServer) DeleteCustomer(_ context.Context, request *customer.DeleteCustomerRequest) (*customer.DeleteCustomerResponse, error) {
	var c models.Customer
	log.Printf("delete customer with id, %d\n", request.Id)

	err := server.H.DB.Model(&models.Customer{}).
		Where("customers.id = ?", request.Id).
		First(&c).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &customer.DeleteCustomerResponse{Status: http.StatusNotFound, Error: err.Error()}, nil
		}
		return &customer.DeleteCustomerResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	err = server.H.DB.Delete(&c).Error

	if err != nil {
		log.Printf("Failed to delete: %v", err)
		return &customer.DeleteCustomerResponse{
			Status: http.StatusBadRequest, Error: "Failed to delete",
		}, nil
	}

	return &customer.DeleteCustomerResponse{
		Status: http.StatusOK,
	}, nil
}

func mapCustomer(d *models.Customer) (*model.CustomerData, error) {
	log.Printf("Marsha to proto type")
	messageInBytes, err := json.Marshal(d)
	if err != nil {
		log.Printf("Marshal Error: %v,", err)
		return nil, err
	}
	log.Printf("raw data:- %s", messageInBytes)
	response := model.CustomerData{}

	// ignore unknown fields
	unMarshaller := &protojson.UnmarshalOptions{DiscardUnknown: true}
	err = unMarshaller.Unmarshal(messageInBytes, &response)
	if err != nil {
		log.Printf("Unmarshal Error: %v,", err)
		return nil, err
	}
	log.Printf("message:- %v", &response)

	return &response, nil
}
