package services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/bdarge/api/out/model"
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
	log.Printf("get user with id, %d\n", request.Id)

	err := server.H.DB.Model(&models.User{}).
		Preload("Roles").
		Preload("Account").
		Joins("Address").
		Where("users.id = ?", request.Id).
		First(&u).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Failed to load user %v", err)
			return &profile.GetUserResponse{Status: http.StatusNotFound, Error: "User not found"}, nil
		}
		log.Printf("Failed to load user %v", err)
		return &profile.GetUserResponse{Status: http.StatusInternalServerError, Error: "Failed to load user"}, nil
	}

	log.Printf("user found: %v", u)

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
	log.Printf("update user (id = %d)\n data=%v", request.Id, request.Data)

	err := server.H.DB.Model(&models.User{}).
		Joins("Address").
		Where("users.id = ?", request.Id).
		First(&user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &profile.UpdateUserResponse{Status: http.StatusNotFound, Error: err.Error()}, nil
		}
		return &profile.UpdateUserResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	user.UserName = request.Data.Username

	if request.Data.Address.Landline != user.Address.Landline {
		user.Address.Landline = request.Data.Address.Landline
	}
	if request.Data.Address.Mobile != user.Address.Mobile {
		user.Address.Mobile = request.Data.Address.Mobile
	}
	if request.Data.Address.Street != user.Address.Street {
		user.Address.Street = request.Data.Address.Street
	}
	if request.Data.Address.City != user.Address.City {
		user.Address.City = request.Data.Address.City
	}
	if request.Data.Address.Country != user.Address.Country {
		user.Address.Country = request.Data.Address.Country
	}
	if request.Data.Address.PostalCode != user.Address.PostalCode {
		user.Address.PostalCode = request.Data.Address.PostalCode
	}

	err = server.H.DB.Model(&user).Save(user).Error

	if err != nil {
		log.Printf("Failed to update: %v", err)
		return &profile.UpdateUserResponse{
			Status: http.StatusBadRequest, Error: "Failed to update",
		}, nil
	}

	err = server.H.DB.Model(&models.User{}).
		Preload("Roles").
		Preload("Account").
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

func (server *ProfileServer) UpdateBusiness(_ context.Context, request *profile.UpdateBusinessRequest) (*profile.UpdateBusinessResponse, error) {
	var b models.Business
	log.Printf("update business (id = %d)\n data=%v", request.Id, request.Data)

	err := server.H.DB.Model(&models.Business{}).
		Where("businesses.id = ?", request.Id).
		First(&b).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &profile.UpdateBusinessResponse{Status: http.StatusNotFound, Error: err.Error()}, nil
		}
		return &profile.UpdateBusinessResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	b.Name = request.Data.Name

	if request.Data.Landline != b.Landline {
		b.Landline = request.Data.Landline
	}
	if request.Data.Mobile != b.Mobile {
		b.Mobile = request.Data.Mobile
	}
	if request.Data.Street != b.Street {
		b.Street = request.Data.Street
	}
	if request.Data.City != b.City {
		b.City = request.Data.City
	}
	if request.Data.Country != b.Country {
		b.Country = request.Data.Country
	}
	if request.Data.PostalCode != b.PostalCode {
		b.PostalCode = request.Data.PostalCode
	}

	err = server.H.DB.Model(&b).Save(b).Error

	if err != nil {
		log.Printf("Failed to update: %v", err)
		return &profile.UpdateBusinessResponse{
			Status: http.StatusBadRequest, Error: "Failed to update",
		}, nil
	}

	err = server.H.DB.Model(&models.Business{}).
		Where("businesses.id = ?", request.Id).
		First(&b).
		Error

	if err != nil {
		log.Printf("Failed to get updated data: %v", err)
		return &profile.UpdateBusinessResponse{
			Status: http.StatusBadRequest, Error: "Failed to get updated data",
		}, nil
	}

	response, err := mapBusiness(&b)

	if err != nil {
		log.Printf("Mapping to map to proto type failed: %v", err)
		return &profile.UpdateBusinessResponse{
			Status: http.StatusBadRequest, Error: "Failed to map to proto type",
		}, nil
	}

	return &profile.UpdateBusinessResponse{
		Status: http.StatusOK,
		Data:   response,
	}, nil
}

func (server *ProfileServer) GetBusiness(_ context.Context, request *profile.GetBusinessRequest) (*profile.GetBusinessResponse, error) {
	var u models.Business
	log.Printf("get business profile with id, %d\n", request.Id)

	err := server.H.DB.Model(&models.Business{}).
		Where("businesses.id = ?", request.Id).
		First(&u).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &profile.GetBusinessResponse{Status: http.StatusNotFound, Error: err.Error()}, nil
		}
		return &profile.GetBusinessResponse{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	log.Printf("business found: %v", u)

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

func mapBusiness(d *models.Business) (*model.BusinessData, error) {
	log.Printf("Marsha to proto type")
	messageInBytes, err := json.Marshal(d)
	if err != nil {
		log.Printf("Marshal Error: %v,", err)
		return nil, err
	}
	log.Printf("raw data:- %s", messageInBytes)
	response := model.BusinessData{}

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

func mapUser(d *models.User) (*model.UserData, error) {
	log.Printf("Marsha to proto type")
	messageInBytes, err := json.Marshal(d)
	if err != nil {
		log.Printf("Marshal Error: %v,", err)
		return nil, err
	}
	log.Printf("raw data:- %s", messageInBytes)
	response := model.UserData{}

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
