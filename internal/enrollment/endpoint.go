package enrollment

import (
	"context"
	"errors"

	courseSdk "github.com/zchelalo/go_microservices_course_sdk/course"
	"github.com/zchelalo/go_microservices_meta/meta"
	"github.com/zchelalo/go_microservices_response/response"
	userSdk "github.com/zchelalo/go_microservices_user_sdk/user"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	Endpoints struct {
		Create Controller
		GetAll Controller
		Update Controller
	}

	CreateRequest struct {
		UserId   string `json:"user_id"`
		CourseId string `json:"course_id"`
	}

	GetAllRequest struct {
		UserId   string
		CourseId string
		Limit    int
		Page     int
	}

	UpdateRequest struct {
		Id     string
		Status *string `json:"status"`
	}

	Config struct {
		LimPageDef string
	}
)

func MakeEndpoints(service Service, config Config) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(service),
		GetAll: makeGetAllEndpoint(service, config),
		Update: makeUpdateEndpoint(service),
	}
}

func makeCreateEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)

		if req.UserId == "" {
			return nil, response.BadRequest(ErrUserIdRequired.Error())
		}

		if req.CourseId == "" {
			return nil, response.BadRequest(ErrCouseIdRequired.Error())
		}

		enrollment, err := service.Create(ctx, req.UserId, req.CourseId)
		if err != nil {
			if errors.As(err, &userSdk.ErrNotFound{}) || errors.As(err, &courseSdk.ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", enrollment, nil), nil
	}
}

func makeGetAllEndpoint(service Service, config Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAllRequest)

		filters := Filters{
			UserId:   req.UserId,
			CourseId: req.CourseId,
		}

		count, err := service.Count(ctx, filters)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		meta, err := meta.New(req.Page, req.Limit, count, config.LimPageDef)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		enrollments, err := service.GetAll(ctx, filters, meta.Offset(), meta.Limit())
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", enrollments, meta), nil
	}
}

func makeUpdateEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)

		if req.Status != nil && *req.Status == "" {
			return nil, response.BadRequest(ErrStatusRequired.Error())
		}

		err := service.Update(ctx, req.Id, req.Status)
		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			if errors.As(err, &ErrInvalidStatus{}) {
				return nil, response.BadRequest(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", "enrollment updated successfully", nil), nil
	}
}
