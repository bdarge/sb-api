package services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/bdarge/api/out/disposition"
	"github.com/bdarge/api/pkg/db"
	"github.com/bdarge/api/pkg/models"
	"github.com/bdarge/api/pkg/util"
	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
)

// Server https://stackoverflow.com/a/69480218
type Server struct {
	H db.Handler
	disposition.UnimplementedDispositionServiceServer
}

func (server *Server) CreateDisposition(_ context.Context, request *disposition.CreateDispositionRequest) (*disposition.CreateDispositionResponse, error) {
	d := &models.Disposition{}
	value, err := protojson.Marshal(request)
	if err != err {
		return &disposition.CreateDispositionResponse{Status: http.StatusBadRequest, Error: err.Error()}, nil
	}
	log.Printf("disposition constructed from a message: %s", value)
	err = json.Unmarshal(value, d)
	log.Printf("disposition model constructed from bytes: %v", d)
	if err != err {
		return &disposition.CreateDispositionResponse{Status: http.StatusBadRequest, Error: err.Error()}, nil
	}
	err = server.H.DB.Create(&d).
		Error

	if err != nil {
		return &disposition.CreateDispositionResponse{Status: http.StatusBadRequest, Error: err.Error()}, nil
	}
	return &disposition.CreateDispositionResponse{
		Status: http.StatusCreated,
		Id:     d.ID,
	}, nil
}

func (server *Server) GetDisposition(_ context.Context, request *disposition.GetDispositionRequest) (*disposition.GetDispositionResponse, error) {
	var d models.Disposition
	log.Printf("get disposition with id, %d\n", request.Id)

	err := server.H.DB.Model(&models.Disposition{}).
		Preload("Items").
		Where("id = ?", request.Id).
		First(&d).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &disposition.GetDispositionResponse{Status: http.StatusNotFound, Error: err.Error()}, nil
		}
		return &disposition.GetDispositionResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	log.Printf("disposition found: %v", d)

	messageInBytes, err := json.Marshal(d)
	if err != nil {
		return &disposition.GetDispositionResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}
	log.Printf("stringfy data in bytes: %s", messageInBytes)
	response := disposition.DispositionData{}

	// ignore unknown fields
	unMarshaller := &protojson.UnmarshalOptions{DiscardUnknown: true}
	err = unMarshaller.Unmarshal(messageInBytes, &response)
	if err != nil {
		log.Printf("Error: %v,", err)
		return &disposition.GetDispositionResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}
	log.Printf("proto message: %v", &response)

	return &disposition.GetDispositionResponse{
		Status: http.StatusOK,
		Data:   &response,
	}, nil
}

func (server *Server) GetDispositions(_ context.Context, request *disposition.GetDispositionsRequest) (*disposition.GetDispositionsResponse, error) {
	log.Printf("get all dispositions, %v", request)

	var dispositions = make([]models.Disposition, 0)

	if request.Limit == 0 {
		request.Limit = 10
	}

	if request.SortDirection == "" {
		request.SortDirection = "desc"
	} else {
		request.SortDirection = strings.ToLower(request.SortDirection)
	}

	if request.SortProperty == "" {
		request.SortProperty = "dispositions.id"
	} else {
		request.SortProperty = "dispositions." + util.ToSnakeCase(request.SortProperty)
	}

	log.Printf("request: %v", request)

	err := server.H.DB.Model(&models.Disposition{}).
		Preload("Items").
		Where("true = ? Or RequestType = ?", request.RequestType == "", request.RequestType).
		Where("true = ? Or Description LIKE ?", request.Search == "", "%"+request.Search+"%").
		Limit(int(request.Limit)).
		Offset(int(request.Page * request.Limit)).
		Order(request.SortProperty + " " + request.SortDirection).
		Find(&dispositions).
		Error

	if err != nil {
		return &disposition.GetDispositionsResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	var total int64

	server.H.DB.Model(&models.Disposition{}).Count(&total)

	page := models.Dispositions{
		Data:  dispositions,
		Limit: request.Limit,
		Page:  request.Page,
		Total: uint32(total),
	}

	if err != nil {
		return &disposition.GetDispositionsResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	messageInBytes, err := json.Marshal(page)
	if err != nil {
		return &disposition.GetDispositionsResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}
	log.Printf("dispositions found: %s", messageInBytes)
	var response disposition.GetDispositionsResponse

	// ignore unknown fields
	unMarshaller := &protojson.UnmarshalOptions{DiscardUnknown: true}
	err = unMarshaller.Unmarshal(messageInBytes, &response)
	if err != nil {
		return &disposition.GetDispositionsResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	response.Status = http.StatusOK
	return &response, nil
}
