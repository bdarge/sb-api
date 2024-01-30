package services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/bdarge/api/out/model"
	"github.com/bdarge/api/out/transactionItem"
	"github.com/bdarge/api/pkg/db"
	"github.com/bdarge/api/pkg/models"
	"github.com/bdarge/api/pkg/util"
	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/gorm"
	// "gorm.io/gorm/clause"
	"log"
	"net/http"
	"strings"
)

// TransactionItemServer transation item server
type TransactionItemServer struct {
	H db.Handler
	transactionItem.UnimplementedTransactionItemServiceServer
}

// CreateTransactionItem creates a new transactionItem item
func (server *TransactionItemServer) CreateTransactionItem(
	_ context.Context, 
	request *transactionItem.CreateTransactionItemRequest) (*transactionItem.CreateTransactionItemResponse, error) {

		d := &models.TransactionItem{}
		log.Printf("Create transactionItem %v", request)
		marshaller := &protojson.MarshalOptions{EmitUnpopulated: false}
		value, err := marshaller.Marshal(request)
	
		if err != err {
			log.Printf("Failed to create transactionItem: %v", err)
			return &transactionItem.CreateTransactionItemResponse{
				Status: http.StatusBadRequest,
				Error: "Failed to create transactionItem"}, nil
			}
	
		log.Printf("transactionItem constructed from a message: %s", value)
		err = json.Unmarshal(value, d)

		log.Printf("transactionItem model constructed from bytes: %v", d)
		if err != nil {
			log.Printf("Failed to create transactionItem: %v", err)
			return &transactionItem.CreateTransactionItemResponse{Status: http.StatusBadRequest, Error: "Failed to create transactionItem"}, nil
		}
		err = server.H.DB.Create(&d).Error
	
		if err != nil {
			log.Printf("Failed to create a transactionItem: %v", err)
			return &transactionItem.CreateTransactionItemResponse{
				Status: http.StatusInternalServerError, Error: "Failed to create a transactionItem",
			}, nil
		}
		return &transactionItem.CreateTransactionItemResponse{
			Status: http.StatusCreated,
			Id:     d.ID,
		}, nil
	}


// GetTransactionItem returns a transactionItem
func (server *TransactionItemServer) GetTransactionItem(_ context.Context, request *transactionItem.GetTransactionItemRequest) (*transactionItem.GetTransactionItemResponse, error) {
	var d models.TransactionItem
	log.Printf("get transactionItem with id, %d\n", request.Id)

	err := server.H.DB.Model(&models.TransactionItem{}).
		Where("id = ?", request.Id).
		First(&d).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &transactionItem.GetTransactionItemResponse{Status: http.StatusNotFound, Error: err.Error()}, nil
		}
		return &transactionItem.GetTransactionItemResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	log.Printf("transactionItem found: %v", d)

	response, err := mapTransactionItem(&d)

	if err != nil {
		log.Printf("Mapping proto type failed: %v", err)
		return &transactionItem.GetTransactionItemResponse{
			Status: http.StatusBadRequest, Error: "Mapping proto type failed",
		}, nil
	}

	return &transactionItem.GetTransactionItemResponse{
		Status: http.StatusOK,
		Data:   response,
	}, nil
}

// GetTransactionItems returns a list of transactionItem in data property
func (server *TransactionItemServer) GetTransactionItems(_ context.Context, request *transactionItem.GetTransactionItemsRequest) (*transactionItem.GetTransactionItemsResponse, error) {
	log.Printf("get all transaction Items, %v", request)

	var transactionItems = make([]models.TransactionItem, 0)

	if request.Limit == 0 {
		request.Limit = 10
	}

	if request.SortDirection == "" {
		request.SortDirection = "desc"
	} else {
		request.SortDirection = strings.ToLower(request.SortDirection)
	}

	if request.SortProperty == "" {
		request.SortProperty = "transaction_items.id"
	} else {
		request.SortProperty = "transaction_items." + util.ToSnakeCase(request.SortProperty)
	}

	log.Printf("request: %v", request)

	err := server.H.DB.Model(&models.TransactionItem{}).
		Where("transaction_id = ?", request.TransactionId).
		Limit(int(request.Limit)).
		Offset(int(request.Page * request.Limit)).
		Order(request.SortProperty + " " + request.SortDirection).
		Find(&transactionItems).
		Error

	if err != nil {
		return &transactionItem.GetTransactionItemsResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	var total int64

	server.H.DB.Model(&models.TransactionItem{}).
		Count(&total)

	page := models.TransactionItems{
		Data:  transactionItems,
		Limit: request.Limit,
		Page:  request.Page,
		Total: uint32(total),
	}

	if err != nil {
		return &transactionItem.GetTransactionItemsResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	messageInBytes, err := json.Marshal(page)
	if err != nil {
		return &transactionItem.GetTransactionItemsResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}
	log.Printf("transactions: %s", messageInBytes)
	var response transactionItem.GetTransactionItemsResponse

	// ignore unknown fields
	unMarshaller := &protojson.UnmarshalOptions{DiscardUnknown: true}
	err = unMarshaller.Unmarshal(messageInBytes, &response)
	if err != nil {
		log.Printf("Failed to get transactions: %v", err)
		return &transactionItem.GetTransactionItemsResponse{Status: http.StatusInternalServerError, Error: "Failed to get transactionItems"}, nil
	}

	log.Printf("message:- %v", &response)
	response.Status = http.StatusOK

	return &response, nil
}

// UpdateTransactionItem updates a transactionItem, and returns it
func (server *TransactionItemServer) UpdateTransactionItem(_ context.Context, request *transactionItem.UpdateTransactionItemRequest) (*transactionItem.UpdateTransactionItemResponse, error) {
	var d models.TransactionItem
	log.Printf("update transactionItem (id = %d)\n", request.Id)

	err := server.H.DB.Model(&models.TransactionItem{}).
		Where("id = ?", request.Id).
		First(&d).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &transactionItem.UpdateTransactionItemResponse{Status: http.StatusNotFound, Error: err.Error()}, nil
		}
		return &transactionItem.UpdateTransactionItemResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	value, err := protojson.Marshal(request)

	if err != err {
		return &transactionItem.UpdateTransactionItemResponse{Status: http.StatusBadRequest, Error: err.Error()}, nil
	}

	log.Printf("update transactionItem constructed from a proto message: %s", value)

	err = server.H.DB.Model(d).Updates(request.Data).Error

	if err != nil {
		log.Printf("Failed to update: %v", err)
		return &transactionItem.UpdateTransactionItemResponse{
			Status: http.StatusBadRequest, Error: "Failed to update",
		}, nil
	}

	err = server.H.DB.Model(&models.TransactionItem{}).
		Where("id = ?", request.Id).
		First(&d).
		Error

	if err != nil {
		log.Printf("Failed to get updated data: %v", err)
		return &transactionItem.UpdateTransactionItemResponse{
			Status: http.StatusBadRequest, Error: "Failed to get updated data",
		}, nil
	}

	response, err := mapTransactionItem(&d)

	if err != nil {
		log.Printf("Mapping to proto type failed: %v", err)
		return &transactionItem.UpdateTransactionItemResponse{
			Status: http.StatusBadRequest, Error: "Error occurred while updating transactionItem",
		}, nil
	}

	return &transactionItem.UpdateTransactionItemResponse{
		Status: http.StatusOK,
		Data:   response,
	}, nil
}

// DeleteTransactionItem deletes a transactionItem
func (server *TransactionItemServer) DeleteTransactionItem(_ context.Context, request *transactionItem.DeleteTransactionItemRequest) (*transactionItem.DeleteTransactionItemResponse, error) {
	var d models.TransactionItem
	log.Printf("delete transactionItem with id: %d\n", request.Id)

	err := server.H.DB.Model(&models.TransactionItem{}).
		Where("id = ?", request.Id).
		First(&d).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &transactionItem.DeleteTransactionItemResponse{Status: http.StatusNotFound, Error: err.Error()}, nil
		}
		log.Printf("Failed to get a transactionItem: %v", err.Error())
		return &transactionItem.DeleteTransactionItemResponse{
			Status: http.StatusInternalServerError, Error: "Failed to get a transactionItem",
		}, nil
	}

	err = server.H.DB.Delete(&d).Error

	if err != nil {
		log.Printf("Failed to delete: %v", err)
		return &transactionItem.DeleteTransactionItemResponse{
			Status: http.StatusBadRequest, Error: "Failed to delete",
		}, nil
	}

	return &transactionItem.DeleteTransactionItemResponse{
		Status: http.StatusOK,
	}, nil
}

func mapTransactionItem(d *models.TransactionItem) (*model.TransactionItem, error) {
	log.Printf("Marsha to proto type")
	messageInBytes, err := json.Marshal(d)
	if err != nil {
		log.Printf("Marshal Error: %v,", err)
		return nil, err
	}
	log.Printf("raw data:- %s", messageInBytes)
	response := model.TransactionItem{}

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

