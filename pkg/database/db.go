package database

import (
	"fmt"

	categoryModel "boilerplate/internal/category/model"
	userModel "boilerplate/internal/user/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(host, port, user, password, dbname string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto Migrate
	err = db.AutoMigrate(&userModel.User{}, &categoryModel.Category{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
