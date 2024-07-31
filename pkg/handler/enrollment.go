package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/zchelalo/go_microservices_enrollment/internal/enrollment"
	"github.com/zchelalo/go_microservices_response/response"
)

func NewEnrollmentHTTPServer(ctx context.Context, endpoints enrollment.Endpoints) http.Handler {
	router := http.NewServeMux()

	opts := []httpTransport.ServerOption{
		httpTransport.ServerErrorEncoder(encodeError),
	}

	router.Handle("POST /enrollments", httpTransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateEnrollment,
		encodeResponse,
		opts...,
	))
	router.Handle("GET /enrollments", httpTransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllEnrollment,
		encodeResponse,
		opts...,
	))
	router.Handle("PATCH /enrollments/{id}", httpTransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateEnrollment,
		encodeResponse,
		opts...,
	))

	return router
}

func decodeCreateEnrollment(_ context.Context, request *http.Request) (interface{}, error) {
	var req enrollment.CreateRequest

	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	return req, nil
}

func decodeGetAllEnrollment(_ context.Context, request *http.Request) (interface{}, error) {
	queries := request.URL.Query()

	limit, _ := strconv.Atoi(queries.Get("limit"))
	page, _ := strconv.Atoi(queries.Get("page"))

	req := enrollment.GetAllRequest{
		UserId:   queries.Get("user_id"),
		CourseId: queries.Get("course_id"),
		Limit:    limit,
		Page:     page,
	}

	return req, nil
}

func decodeUpdateEnrollment(_ context.Context, request *http.Request) (interface{}, error) {
	var req enrollment.UpdateRequest

	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	id := request.PathValue("id")
	req.Id = id

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}
