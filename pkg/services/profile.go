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
	var u models.User
	log.Printf("get profile with id, %d\n", request.Id)

	err := server.H.DB.Model(&models.User{}).
		Preload("Roles").
		Preload("Business.Address").
		Joins("Business").
		Joins("Address").
		Where("users.id = ?", request.Id).
		First(&u).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &profile.GetUserResponse{Status: http.StatusNotFound, Error: err.Error()}, nil
		}
		return &profile.GetUserResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	log.Printf("profile found: %v", u)

	response, err := mapUser(&u)

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

func (server *ProfileServer) UpdateUser(_ context.Context, request *profile.UpdateUserRequest) (*profile.UpdateUserResponse, error) {
	var user models.User
	log.Printf("update user (id = %d)\n", request.Id)

	err := server.H.DB.Model(&models.User{}).
		Where("users.id = ?", request.Id).
		First(&user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &profile.UpdateUserResponse{Status: http.StatusNotFound, Error: err.Error()}, nil
		}
		return &profile.UpdateUserResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	value, err := protojson.Marshal(request)
	if err != nil {
		return &profile.UpdateUserResponse{Status: http.StatusBadRequest, Error: err.Error()}, nil
	}
	log.Printf("update profile constructed from a proto message: %s", value)

	if err != nil {
		return &profile.UpdateUserResponse{Status: http.StatusBadRequest, Error: err.Error()}, nil
	}

	err = server.H.DB.Model(user).Updates(request.Data).Error

	if err != nil {
		log.Printf("Failed to update: %v", err)
		return &profile.UpdateUserResponse{
			Status: http.StatusBadRequest, Error: "Failed to update",
		}, nil
	}

	err = server.H.DB.Model(&models.User{}).
		Preload("Roles").
		Preload("Business.Address").
		Joins("Business").
		Joins("Address").
		Where("users.id = ?", request.Id).
		First(&user).
		Error

	if err != nil {
		log.Printf("Failed to get updated data: %v", err)
		return &profile.UpdateUserResponse{
			Status: http.StatusBadRequest, Error: "Failed to get updated data",
		}, nil
	}

	response, err := mapUser(&user)

	if err != nil {
		log.Printf("Mapping to map to proto type failed: %v", err)
		return &profile.UpdateUserResponse{
			Status: http.StatusBadRequest, Error: "Failed to map to proto type",
		}, nil
	}

	return &profile.UpdateUserResponse{
		Status: http.StatusOK,
		Data:   response,
	}, nil
}

func (server *ProfileServer) GetBusiness(_ context.Context, request *profile.GetBusinessRequest) (*profile.GetBusinessResponse, error) {
	var u models.Business
	log.Printf("get business profile with id, %d\n", request.Id)

	err := server.H.DB.Model(&models.Business{}).
		Joins("Address").
		Where("businesses.id = ?", request.Id).
		First(&u).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &profile.GetBusinessResponse{Status: http.StatusNotFound, Error: err.Error()}, nil
		}
		return &profile.GetBusinessResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	log.Printf("profile found: %v", u)

	response, err := mapBusiness(&u)

	if err != nil {
		log.Printf("Mapping proto type failed: %v", err)
		return &profile.GetBusinessResponse{
			Status: http.StatusBadRequest, Error: "Mapping proto type failed",
		}, nil
	}

	return &profile.GetBusinessResponse{
		Status: http.StatusOK,
		Data:   response,
	}, nil
}

func mapBusiness(d *models.Business) (*profile.BusinessData, error) {
	log.Printf("Marsha to proto type")
	messageInBytes, err := json.Marshal(d)
	if err != nil {
		log.Printf("Marshal Error: %v,", err)
		return nil, err
	}
	log.Printf("raw data:- %s", messageInBytes)
	response := profile.BusinessData{}

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
