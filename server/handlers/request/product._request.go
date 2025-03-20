package request

import (
	"github.com/shopspring/decimal"
)

type ProductRequest struct {
	Name        string          `json:"name" binding:"required"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price" `
	Stock       int             `json:"stock" `
	Sales       int             `json:"sales" `
	Category    string          `json:"category" `
	Rating      float64         `json:"rating" `
	Images      []string        `json:"images"`
	Tags        []string        `json:"tags"`
	Status      string          `json:"status" binding:"omitempty,oneof=active inactive"`
}
/**

UpdataProductRequest + Param(id)

CreateProductRequest
*/
