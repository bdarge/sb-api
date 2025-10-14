package services

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/bdarge/api/out/lang"
	"github.com/bdarge/api/pkg/db"
	"github.com/bdarge/api/pkg/models"
	"google.golang.org/protobuf/encoding/protojson"
)

// LangServer language server
type LangServer struct {
	H db.Handler
	lang.UnimplementedLangServiceServer
}

// GetLang returns languages look up
func (server *LangServer) GetLang(ctx context.Context, request *lang.LangGetRequest) (*lang.LangGetResponse, error) {
	log.Printf("Get a list of languages")

	var languages = make([]models.Lang, 0)

	err := server.H.DB.Model(&models.Lang{}).
		Find(&languages).
		Error

	if err != nil {
		return nil, err
	}
	messageInBytes, err := json.Marshal(models.Langs {
		Data: languages,
	})
	if err != nil {
		return &lang.LangGetResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}
	var response lang.LangGetResponse
	// ignore unknown fields
	unMarshaller := &protojson.UnmarshalOptions{DiscardUnknown: true}
	err = unMarshaller.Unmarshal(messageInBytes, &response)
	if err != nil {
		log.Printf("Error: %v", err)
		return &lang.LangGetResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}
	
	return &response, nil
}
