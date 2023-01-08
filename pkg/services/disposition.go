package services

import (
	"context"
	"encoding/json"
	"github.com/bdarge/api/pkg/db"
	"github.com/bdarge/api/pkg/disposition"
	"github.com/bdarge/api/pkg/models"
	"github.com/bdarge/api/pkg/util"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net/http"
	"strings"
)

// Server https://stackoverflow.com/a/69480218
type Server struct {
	H db.Handler
	disposition.UnimplementedDispositionServiceServer
}

func (server *Server) CreateDisposition(ctx context.Context, request *disposition.CreateDispositionRequest) (*disposition.CreateDispositionResponse, error) {
	d := &models.Disposition{}
	err := util.Recast(request, d)
	err = server.H.DB.Create(&d).
		Error

	if err != nil {
		return &disposition.CreateDispositionResponse{Status: http.StatusBadRequest, Error: err.Error()}, nil
	} else {
		return &disposition.CreateDispositionResponse{
			Status: http.StatusCreated,
			Id:     int64(d.ID),
		}, nil
	}
}

func (server *Server) GetDisposition(ctx context.Context, request *disposition.GetDispositionRequest) (*disposition.GetDispositionResponse, error) {
	var d models.Disposition
	log.Printf("get disposition with id, %d\n", request.Id)

	err := server.H.DB.Where("id = ?", request.Id).
		First(&d).
		Error

	if err != nil {
		return &disposition.GetDispositionResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	} else {
		data := &disposition.DispositionData{}
		err = util.Recast(d, data)
		return &disposition.GetDispositionResponse{
			Status: http.StatusOK,
			Data:   data,
		}, nil
	}
}

func (server *Server) GetDispositions(ctx context.Context, request *disposition.GetDispositionsRequest) (*disposition.GetDispositionsResponse, error) {
	log.Printf("get all dispositions, filter by requestType if set, request: %v\n", request)

	var dispositions = make([]models.Disposition, 0)

	log.Printf("get order or/and quotes with or without filters; filter=%s; sort porperty=%s\n", request.Search, request.SortProperty)

	if request.Size == 0 {
		request.Size = 10
	}

	if request.SortDirection == "" {
		request.SortDirection = "desc"
	} else {
		request.SortDirection = strings.ToLower(request.SortDirection)
	}

	if request.SortProperty == "" {
		request.SortProperty = "disposition.id"
	} else {
		request.SortProperty = "disposition." + util.ToSnakeCase(request.SortProperty)
	}

	rows, err := server.H.DB.
		Model(&models.Disposition{}).
		Select("disposition.*").
		Limit(int(request.Size)).
		Offset(int(request.Page*request.Size)).
		Where("true = ?", request.RequestType == "").
		Or("RequestType = ?", "%"+request.RequestType+"%").
		Where("true = ?", request.Search == "").
		Or("Description LIKE ?", "%"+request.Search+"%").
		Order(request.SortProperty + " " + request.SortDirection).
		Rows()

	if err != nil {
		return &disposition.GetDispositionsResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	defer rows.Close()

	for rows.Next() {
		var d models.Disposition
		err = server.H.DB.ScanRows(rows, &d)
		if err != nil {
			log.Fatal(err)
		}
		dispositions = append(dispositions, d)
	}

	var total int64

	server.H.DB.Model(&disposition.DispositionData{}).Count(&total)

	page := models.Dispositions{
		Data: dispositions,
		Page: models.Page{
			Page:  int(request.Page),
			Size:  int(request.Size),
			Total: total,
		},
	}

	if err != nil {
		return &disposition.GetDispositionsResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	b, err := json.Marshal(page)

	if err != nil {
		return &disposition.GetDispositionsResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	var result disposition.GetDispositionsResponse
	err = protojson.Unmarshal(b, &result)

	if err != nil {
		return &disposition.GetDispositionsResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	result.Status = http.StatusOK
	return &result, nil
}
