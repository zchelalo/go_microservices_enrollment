package enrollment

import (
	"context"
	"log"

	"github.com/zchelalo/go_microservices_domain/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(ctx context.Context, enrollment *domain.Enrollment) error
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Enrollment, error)
		Update(ctx context.Context, id string, status *string) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	repository struct {
		log *log.Logger
		db  *gorm.DB
	}
)

func NewRepository(log *log.Logger, db *gorm.DB) Repository {
	return &repository{
		log: log,
		db:  db,
	}
}

func (repo *repository) Create(ctx context.Context, enrollment *domain.Enrollment) error {
	if err := repo.db.WithContext(ctx).Create(enrollment).Error; err != nil {
		repo.log.Println(err)
		return err
	}

	repo.log.Println("enrollment created with id: ", enrollment.Id)
	return nil
}

func (repo *repository) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Enrollment, error) {
	var enrollments []domain.Enrollment

	tx := repo.db.WithContext(ctx).Model(&enrollments)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	if err := tx.Order("created_at desc").Find(&enrollments).Error; err != nil {
		repo.log.Println(err)
		return nil, err
	}

	return enrollments, nil
}

func (repo *repository) Update(ctx context.Context, id string, status *string) error {
	values := make(map[string]interface{})

	if status != nil {
		values["status"] = *status
	}

	result := repo.db.WithContext(ctx).Model(&domain.Enrollment{}).Where("id = ?", id).Updates(values)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		repo.log.Printf("enrollment with id %s doesn't exists", id)
		return ErrNotFound{id}
	}

	return nil
}

func (repo *repository) Count(ctx context.Context, filters Filters) (int, error) {
	var count int64
	tx := repo.db.WithContext(ctx).Model(&domain.Enrollment{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		repo.log.Println(err)
		return 0, err
	}

	return int(count), nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.UserId != "" {
		tx = tx.Where("user_id = ?", filters.UserId)
	}

	if filters.CourseId != "" {
		tx = tx.Where("course_id = ?", filters.CourseId)
	}

	return tx
}
