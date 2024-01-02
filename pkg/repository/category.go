package repository

import (
	"WatchHive/pkg/domain"
	interfaces "WatchHive/pkg/repository/interface"
	"WatchHive/pkg/utils/errmsg"
	"WatchHive/pkg/utils/models"
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

func (cr *categoryRepository) AddCategory(c models.CategoryAdd) (domain.Category, error) {

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
	var count int
	err := cr.DB.Raw("SELECT COUNT(*) FROM categories WHERE id =?", currentId).Scan(&count).Error
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, err
	}

	return true, err
}

func (cr *categoryRepository) UpdateCategory(currentId int, new string) (domain.Category, error) {

	// Check the database connection
	if cr.DB == nil {
		return domain.Category{}, errors.New(errmsg.ErrDbConnect)
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
		return errors.New(errmsg.ErrDatatypeConversion)
	}

	result := cr.DB.Exec("DELETE FROM categories WHERE id = ?", id)

	if result.RowsAffected < 1 {
		return errors.New(errmsg.ErrIdExist)
	}
	return nil
}

func (cr *categoryRepository) CheckCategoryByName(name string) bool {
	var count int
	err := cr.DB.Raw("SELECT COUNT(*) FROM categories WHERE category =?", name).Scan(&count).Error
	if err != nil {
		return false
	}

	if count == 0 {
		return false
	}

	return true
}

func (cr *categoryRepository) GetCategoryId(productId int) (int, error) {
	var catId int
	err := cr.DB.Raw("select category_id from products where id = ?", productId).Scan(&catId).Error
	if err != nil {
		return 0, errors.New(errmsg.ErrGetDB)
	}
	return catId, nil
}
