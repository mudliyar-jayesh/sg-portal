package util

import (
	"gorm.io/gorm"
)

type Repository[T any] struct {
	db *gorm.DB
}

// NewRepository creates a new instance of the repository for any model
func NewRepository[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

func (r *Repository[T]) Create(entry *T) error {
	return r.db.Create(entry).Error
}

func (r *Repository[T]) CreateMultiple(entries *[]T) error {
	return r.db.Create(entries).Error
}

func (r *Repository[T]) GetByField(field string, value interface{}) (*T, error) {
	var entry T
	err := r.db.Where(field+" = ?", value).First(&entry).Error
	return &entry, err
}

func (r *Repository[T]) GetAllByCondition(condition string, args ...interface{}) ([]T, error) {
	var entries []T
	// Pass the condition and any arguments (e.g., "age > ?", 20)
	err := r.db.Where(condition, args...).Find(&entries).Error
	return entries, err
}
func (r *Repository[T]) UpdateOne(field string, value interface{}, updates map[string]interface{}) error {
	return r.db.Model(new(T)).Where(field+" = ?", value).Updates(updates).Error
}

func (r *Repository[T]) UpdateMany(field string, value interface{}, updates map[string]interface{}) error {
	return r.db.Model(new(T)).Where(field+" = ?", value).Updates(updates).Error
}

func (r *Repository[T]) Joins(join string, condition string, args ...interface{}) ([]T, error) {
	var entries []T
	err := r.db.Joins(join).Where(condition, args...).Find(&entries).Error
	return entries, err
}
