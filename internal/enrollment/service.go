package enrollment

import (
	"context"
	"log"

	"github.com/zchelalo/go_microservices_domain/domain"
)

type (
	Filters struct {
		UserId   string
		CourseId string
	}

	Service interface {
		Create(ctx context.Context, userId, courseId string) (*domain.Enrollment, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Enrollment, error)
		Update(ctx context.Context, id string, status *string) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	service struct {
		log        *log.Logger
		repository Repository
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:        log,
		repository: repo,
	}
}

func (srv *service) Create(ctx context.Context, userId, courseId string) (*domain.Enrollment, error) {
	srv.log.Println("create enrollment service")

	enrollment := domain.Enrollment{
		UserId:   userId,
		CourseId: courseId,
		Status:   "P",
	}
	if err := srv.repository.Create(ctx, &enrollment); err != nil {
		return nil, err
	}
	return &enrollment, nil
}

func (srv *service) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Enrollment, error) {
	srv.log.Println("get all enrollment service")
	enrollments, err := srv.repository.GetAll(ctx, filters, offset, limit)
	if err != nil {
		return nil, err
	}
	return enrollments, nil
}

func (srv *service) Update(ctx context.Context, id string, status *string) error {
	srv.log.Println("update enrollment service")
	return srv.repository.Update(ctx, id, status)
}

func (srv *service) Count(ctx context.Context, filters Filters) (int, error) {
	srv.log.Println("count enrollment service")
	return srv.repository.Count(ctx, filters)
}
