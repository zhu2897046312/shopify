package repository

import (
	"database/sql"
	"shopify/models"
	"time"

	"gorm.io/gorm"
)

type ProductRepository struct {
	*BaseRepository
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 创建产品
func (r *ProductRepository) Create(product *models.Product) error {

	// 使用事务来处理创建操作
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 创建产品记录
		if err := tx.Create(product).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetByID 获取产品详情
func (r *ProductRepository) GetProductViewByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Reviews.User").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Update 更新产品信息
func (r *ProductRepository) Update(product *models.Product) error {
	// 确保 Images 和 Tags 字段不为 nil
	if product.Images == nil {
		product.Images = make([]string, 0)
	}
	if product.Tags == nil {
		product.Tags = make([]string, 0)
	}

	// 使用 Model 和 Where 来指定更新的记录
	return r.db.Model(&models.Product{}).
		Where("id = ?", product.ID).
		Updates(models.Product{
			Name:        product.Name,
			Description: product.Description,
			Price:      product.Price,
			Stock:      product.Stock,
			Status:     product.Status,
			Category:   product.Category,
			Images:     product.Images,
			Tags:       product.Tags,
			UpdatedAt:  time.Now(),
		}).Error
}

// Delete 删除产品(软删除)
func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

// List 获取产品列表(支持分页)
func (r *ProductRepository) List(page, pageSize int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	if err := r.db.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&products).Error

	return products, total, err
}

// ListByCategory 按类别查询产品
func (r *ProductRepository) ListByCategory(category string, page, pageSize int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	if err := r.db.Model(&models.Product{}).Where("category = ?", category).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.Where("category = ?", category).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&products).Error

	return products, total, err
}

// ListByPriceRange 按价格区间查询产品
func (r *ProductRepository) ListByPriceRange(minPrice, maxPrice float64, page, pageSize int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	if err := r.db.Model(&models.Product{}).
		Where("price >= ? AND price <= ?", minPrice, maxPrice).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.Where("price >= ? AND price <= ?", minPrice, maxPrice).
		Offset(offset).
		Limit(pageSize).
		Order("price ASC").
		Find(&products).Error

	return products, total, err
}

// ListByTags 按标签查询产品
func (r *ProductRepository) ListByTags(tags []string, page, pageSize int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	if err := r.db.Model(&models.Product{}).
		Where("tags ?| array[:tags]", sql.Named("tags", tags)). // PostgreSQL语法
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.Where("tags ?| array[:tags]", sql.Named("tags", tags)).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&products).Error

	return products, total, err
}

// UpdateStock 更新库存
func (r *ProductRepository) UpdateStock(id uint, quantity int) error {
	return r.db.Model(&models.Product{}).
		Where("id = ? AND stock >= ?", id, -quantity). // 确保库存充足
		UpdateColumn("stock", gorm.Expr("stock + ?", quantity)).Error
}

// UpdateSales 更新销量
func (r *ProductRepository) UpdateSales(id uint, quantity int) error {
	return r.db.Model(&models.Product{}).
		Where("id = ?", id).
		UpdateColumn("sales", gorm.Expr("sales + ?", quantity)).Error
}

// Search 搜索产品
func (r *ProductRepository) Search(keyword string, page, pageSize int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.Model(&models.Product{}).
		Where("name LIKE ? OR description LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := query.Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&products).Error

	return products, total, err
}
