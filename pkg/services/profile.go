package services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/bdarge/api/out/profile"
	"github.com/bdarge/api/pkg/db"
	"github.com/bdarge/api/pkg/models"
	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/gorm"
	"log"
	"net/http"
)

// ProfileServer https://stackoverflow.com/a/69480218
type ProfileServer struct {
	H db.Handler
	profile.UnimplementedProfileServiceServer
}

func (server *ProfileServer) GetUser(_ context.Context, request *profile.GetUserRequest) (*profile.GetUserResponse, error) {
	var d models.User
	log.Printf("get profile with id, %d\n", request.Id)

	err := server.H.DB.Model(&models.User{}).
		Where("users.id = ?", request.Id).
		First(&d).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &profile.GetUserResponse{Status: http.StatusNotFound, Error: err.Error()}, nil
		}
		return &profile.GetUserResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	log.Printf("profile found: %v", d)

	response, err := mapUser(&d)

	if err != nil {
		log.Printf("Mapping proto type failed: %v", err)
		return &profile.GetUserResponse{
			Status: http.StatusBadRequest, Error: "Mapping proto type failed",
		}, nil
	}

	return &profile.GetUserResponse{
		Status: http.StatusOK,
		Data:   response,
	}, nil
}

func mapUser(d *models.User) (*profile.UserData, error) {
	log.Printf("Marsha to proto type")
	messageInBytes, err := json.Marshal(d)
	if err != nil {
		log.Printf("Marshal Error: %v,", err)
		return nil, err
	}
	log.Printf("raw data:- %s", messageInBytes)
	response := profile.UserData{}

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
