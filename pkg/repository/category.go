package repository

import (
	"WatchHive/pkg/domain"
	interfaces "WatchHive/pkg/repository/interface"
	"errors"
	"strconv"

	"gorm.io/gorm"
)

type categoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(DB *gorm.DB) interfaces.CategoryRepository {

	return &categoryRepository{DB}
}

func (cr *categoryRepository) AddCategory(c domain.Category) (domain.Category, error) {

	var b domain.Category
	err := cr.DB.Raw("INSERT INTO categories (category) VALUES (?) RETURNING *", c.Category).Scan(&b).Error
	if err != nil {
		return domain.Category{}, err
	}

	return b, nil
}

func (cr *categoryRepository) GetCategories() ([]domain.Category, error) {
	var Model []domain.Category
	err := cr.DB.Raw("SELECT * FROM categories").Scan(&Model).Error
	if err != nil {
		return []domain.Category{}, err
	}

	return Model, nil
}

func (cr *categoryRepository) CheckCategory(currentId int) (bool, error) {
	var i int
	err := cr.DB.Raw("SELECT COUNT(*) FROM categories WHERE id =?", currentId).Scan(&i).Error
	if err != nil {
		return false, err
	}

	if i == 0 {
		return false, err
	}

	return true, err
}

func (cr *categoryRepository) UpdateCategory(currentId int, new string) (domain.Category, error) {

	// Check the database connection
	if cr.DB == nil {
		return domain.Category{}, errors.New("database connection is nil")
	}

	// Update the category
	if err := cr.DB.Exec("UPDATE categories SET category = $1 WHERE id = $2", new, currentId).Error; err != nil {
		return domain.Category{}, err
	}

	// Retrieve the updated category
	var newcat domain.Category
	if err := cr.DB.First(&newcat, "category = ?", new).Error; err != nil {
		return domain.Category{}, err
	}

	return newcat, nil
}

func (cr *categoryRepository) DeleteCategory(catergoryID string) error {
	id, err := strconv.Atoi(catergoryID)

	if err != nil {
		return errors.New("converting into integers is not happen")
	}

	result := cr.DB.Exec("DELETE FROM categories WHERE id = ?", id)

	if result.RowsAffected < 1 {
		return errors.New("now rows with that id exist")
	}
	return nil
}
