package services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/bdarge/api/out/model"
	"github.com/bdarge/api/out/transaction"
	"github.com/bdarge/api/pkg/db"
	"github.com/bdarge/api/pkg/models"
	"github.com/bdarge/api/pkg/util"
	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"net/http"
	"strings"
)

// Server https://stackoverflow.com/a/69480218
type Server struct {
	H db.Handler
	transaction.UnimplementedTransactionServiceServer
}

func (server *Server) CreateTransaction(_ context.Context, request *transaction.CreateTransactionRequest) (*transaction.CreateTransactionResponse, error) {
	d := &models.Transaction{}
	log.Printf("Create transaction %v", request)
	marshaller := &protojson.MarshalOptions{EmitUnpopulated: false}
	value, err := marshaller.Marshal(request)

	if err != err {
		log.Printf("Failed to create transaction: %v", err)
		return &transaction.CreateTransactionResponse{Status: http.StatusBadRequest, Error: "Failed to create transaction"}, nil
	}

	log.Printf("transaction constructed from a message: %s", value)
	err = json.Unmarshal(value, d)
	log.Printf("transaction model constructed from bytes: %v", d)
	if err != nil {
		log.Printf("Failed to create transaction: %v", err)
		return &transaction.CreateTransactionResponse{Status: http.StatusBadRequest, Error: "Failed to create transaction"}, nil
	}
	err = server.H.DB.Create(&d).
		Error

	if err != nil {
		log.Printf("Failed to create a transaction: %v", err)
		return &transaction.CreateTransactionResponse{
			Status: http.StatusInternalServerError, Error: "Failed to create a transaction",
		}, nil
	}
	return &transaction.CreateTransactionResponse{
		Status: http.StatusCreated,
		Id:     d.ID,
	}, nil
}

func (server *Server) GetTransaction(_ context.Context, request *transaction.GetTransactionRequest) (*transaction.GetTransactionResponse, error) {
	var d models.Transaction
	log.Printf("get transaction with id, %d\n", request.Id)

	err := server.H.DB.Model(&models.Transaction{}).
		Preload("Items").
		Joins("Customer").
		Where("transactions.id = ?", request.Id).
		First(&d).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &transaction.GetTransactionResponse{Status: http.StatusNotFound, Error: err.Error()}, nil
		}
		return &transaction.GetTransactionResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	log.Printf("transaction found: %v", d)

	response, err := mapTransaction(&d)

	if err != nil {
		log.Printf("Mapping proto type failed: %v", err)
		return &transaction.GetTransactionResponse{
			Status: http.StatusBadRequest, Error: "Mapping proto type failed",
		}, nil
	}

	return &transaction.GetTransactionResponse{
		Status: http.StatusOK,
		Data:   response,
	}, nil
}

func (server *Server) GetTransactions(_ context.Context, request *transaction.GetTransactionsRequest) (*transaction.GetTransactionsResponse, error) {
	log.Printf("get all transactions, %v", request)

	var transactions = make([]models.Transaction, 0)

	if request.Limit == 0 {
		request.Limit = 10
	}

	if request.SortDirection == "" {
		request.SortDirection = "desc"
	} else {
		request.SortDirection = strings.ToLower(request.SortDirection)
	}

	if request.SortProperty == "" {
		request.SortProperty = "transactions.id"
	} else {
		request.SortProperty = "transactions." + util.ToSnakeCase(request.SortProperty)
	}

	log.Printf("request: %v", request)

	err := server.H.DB.Model(&models.Transaction{}).
		Preload(clause.Associations).
		Where("true = ? Or RequestType = ?", request.RequestType == "", request.RequestType).
		Where("true = ? Or Description LIKE ?", request.Search == "", "%"+request.Search+"%").
		Limit(int(request.Limit)).
		Offset(int(request.Page * request.Limit)).
		Order(request.SortProperty + " " + request.SortDirection).
		Find(&transactions).
		Error

	if err != nil {
		return &transaction.GetTransactionsResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	var total int64

	server.H.DB.Model(&models.Transaction{}).
		Where("true = ? Or RequestType = ?", request.RequestType == "", request.RequestType).
		Where("true = ? Or Description LIKE ?", request.Search == "", "%"+request.Search+"%").
		Count(&total)

	page := models.Transactions{
		Data:  transactions,
		Limit: request.Limit,
		Page:  request.Page,
		Total: uint32(total),
	}

	if err != nil {
		return &transaction.GetTransactionsResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	messageInBytes, err := json.Marshal(page)
	if err != nil {
		return &transaction.GetTransactionsResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}
	log.Printf("transactions: %s", messageInBytes)
	var response transaction.GetTransactionsResponse

	// ignore unknown fields
	unMarshaller := &protojson.UnmarshalOptions{DiscardUnknown: true}
	err = unMarshaller.Unmarshal(messageInBytes, &response)
	if err != nil {
		log.Printf("Failed to get transactions: %v", err)
		return &transaction.GetTransactionsResponse{Status: http.StatusInternalServerError, Error: "Failed to get transactions"}, nil
	}
	log.Printf("message:- %v", response)
	response.Status = http.StatusOK
	return &response, nil
}

func (server *Server) UpdateTransaction(_ context.Context, request *transaction.UpdateTransactionRequest) (*transaction.UpdateTransactionResponse, error) {
	var d models.Transaction
	log.Printf("update transaction (id = %d)\n", request.Id)

	err := server.H.DB.Model(&models.Transaction{}).
		Where("transactions.id = ?", request.Id).
		First(&d).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &transaction.UpdateTransactionResponse{Status: http.StatusNotFound, Error: err.Error()}, nil
		}
		return &transaction.UpdateTransactionResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	value, err := protojson.Marshal(request)

	if err != err {
		return &transaction.UpdateTransactionResponse{Status: http.StatusBadRequest, Error: err.Error()}, nil
	}

	log.Printf("update transaction constructed from a proto message: %s", value)

	err = server.H.DB.Model(d).Updates(request.Data).Error

	if err != nil {
		log.Printf("Failed to update: %v", err)
		return &transaction.UpdateTransactionResponse{
			Status: http.StatusBadRequest, Error: "Failed to update",
		}, nil
	}

	err = server.H.DB.Model(&models.Transaction{}).
		Where("transactions.id = ?", request.Id).
		First(&d).
		Error

	if err != nil {
		log.Printf("Failed to get updated data: %v", err)
		return &transaction.UpdateTransactionResponse{
			Status: http.StatusBadRequest, Error: "Failed to get updated data",
		}, nil
	}

	response, err := mapTransaction(&d)

	if err != nil {
		log.Printf("Mapping to proto type failed: %v", err)
		return &transaction.UpdateTransactionResponse{
			Status: http.StatusBadRequest, Error: "Error occurred while updating transaction",
		}, nil
	}

	return &transaction.UpdateTransactionResponse{
		Status: http.StatusOK,
		Data:   response,
	}, nil
}

func (server *Server) DeleteTransaction(_ context.Context, request *transaction.DeleteTransactionRequest) (*transaction.DeleteTransactionResponse, error) {
	var d models.Transaction
	log.Printf("delete transaction with id: %d\n", request.Id)

	err := server.H.DB.Model(&models.Transaction{}).
		Where("transactions.id = ?", request.Id).
		First(&d).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &transaction.DeleteTransactionResponse{Status: http.StatusNotFound, Error: err.Error()}, nil
		}
		log.Printf("Failed to get a transaction: %v", err.Error())
		return &transaction.DeleteTransactionResponse{
			Status: http.StatusInternalServerError, Error: "Failed to get a transaction",
		}, nil
	}

	err = server.H.DB.Delete(&d).Error

	if err != nil {
		log.Printf("Failed to delete: %v", err)
		return &transaction.DeleteTransactionResponse{
			Status: http.StatusBadRequest, Error: "Failed to delete",
		}, nil
	}

	return &transaction.DeleteTransactionResponse{
		Status: http.StatusOK,
	}, nil
}

func mapTransaction(d *models.Transaction) (*model.TransactionData, error) {
	log.Printf("Marsha to proto type")
	messageInBytes, err := json.Marshal(d)
	if err != nil {
		log.Printf("Marshal Error: %v,", err)
		return nil, err
	}
	log.Printf("raw data:- %s", messageInBytes)
	response := model.TransactionData{}

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
