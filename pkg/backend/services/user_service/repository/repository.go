package repository

import (
	"fmt"
	"reflect"

	api "github.com/aburifat/go-agro/apis/agro"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func serialize[T any](source, destination *T) {
	srcValue := reflect.ValueOf(source).Elem()
	dstValue := reflect.ValueOf(destination).Elem()

	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Field(i)
		dstField := dstValue.Field(i)

		if !srcField.IsZero() && dstField.CanSet() {
			dstField.Set(srcField)
		}
	}
}

func Create(db *gorm.DB, user *api.User) (string, error) {
	result := db.Create(&user)
	if result.Error != nil {
		return "", fmt.Errorf("failed to insert: %v", result.Error)
	}
	return user.ID, nil
}

func Update[T interface{}](db *gorm.DB, id string, model *T) error {
	result := db.Model(&api.User{}).Where("id = ?", id).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no user found with ID: %s", id)
	}
	return nil
}

func GetById(db *gorm.DB, id string) (*api.User, error) {
	var user api.User
	result := db.First(&user, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Println("No user found with ID:", id)
		} else {
			panic("failed to fetch user: " + result.Error.Error())
		}
	} else {
		fmt.Printf("Fetched user: %+v\n", user)
		fmt.Println("User ID:", user.ID)
	}
	return &user, nil
}

func GetAll[T interface{}](db *gorm.DB, pageNumber, pageSize int) ([]*T, error) {
	skip := int64((pageNumber - 1) * pageSize)
	limit := int64(pageSize)

	var data []*T
	result := db.Limit(int(limit)).Offset(int(skip)).Find(&data)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get users: %v", result.Error)
	}

	return data, nil
}

func Delete(db *gorm.DB, id string) error {
	_, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid uid format: %v", err)
	}

	// Delete the user with the given uid
	result := db.Where("id = ?", id).Delete(&api.User{})
	if err := result.Error; err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return fmt.Errorf("no user found with uid: %s", id)
	}
	return nil
}
