package service

import (
	"errors"
	"shopify/models"
)

type ProductService struct {
	*Service
}

func NewProductService(base *Service) *ProductService {
	return &ProductService{Service: base}
}

func (s *ProductService) GetProduct(id uint) (*models.Product, error) {
	return s.repoFactory.GetProductRepository().GetByID(id)
}

func (s *ProductService) ListProducts(page, pageSize int) ([]models.Product, int64, error) {
	return s.repoFactory.GetProductRepository().List(page, pageSize)
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.Price.IsZero() {
		return errors.New("product price is required")
	}
	if product.Stock < 0 {
		return errors.New("product stock cannot be negative")
	}

	if product.Status != "" && product.Status != "active" && product.Status != "inactive" {
		product.Status = "active"
	}

	if product.Images == nil {
		product.Images = make([]string, 0)
	}
	if product.Tags == nil {
		product.Tags = make([]string, 0)
	}

	return s.repoFactory.GetProductRepository().Create(product)
}

func (s *ProductService) UpdateProduct(product *models.Product) error {
	if _, err := s.GetProduct(product.ID); err != nil {
		return errors.New("product not found")
	}

	if product.Status != "" && product.Status != "active" && product.Status != "inactive" {
		return errors.New("invalid status value")
	}

	return s.repoFactory.GetProductRepository().Update(product)
}

func (s *ProductService) DeleteProduct(id uint) error {
	if _, err := s.GetProduct(id); err != nil {
		return errors.New("product not found")
	}

	return s.repoFactory.GetProductRepository().Delete(id)
}

func (s *ProductService) ListProductsByCategory(category string, page, pageSize int) ([]models.Product, int64, error) {
	return s.repoFactory.GetProductRepository().ListByCategory(category, page, pageSize)
}

func (s *ProductService) ListProductsByPriceRange(minPrice, maxPrice float64, page, pageSize int) ([]models.Product, int64, error) {
	if minPrice > maxPrice {
		return nil, 0, errors.New("invalid price range")
	}
	return s.repoFactory.GetProductRepository().ListByPriceRange(minPrice, maxPrice, page, pageSize)
}

func (s *ProductService) ListProductsByTags(tags []string, page, pageSize int) ([]models.Product, int64, error) {
	if len(tags) == 0 {
		return nil, 0, errors.New("tags cannot be empty")
	}
	return s.repoFactory.GetProductRepository().ListByTags(tags, page, pageSize)
}

func (s *ProductService) SearchProducts(keyword string, page, pageSize int) ([]models.Product, int64, error) {
	if keyword == "" {
		return nil, 0, errors.New("search keyword cannot be empty")
	}
	return s.repoFactory.GetProductRepository().Search(keyword, page, pageSize)
} 