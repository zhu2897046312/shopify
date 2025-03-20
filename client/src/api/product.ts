import request from '@/utils/request'
import type {
  Product,
  ListProductParams,
  ListProductResponse,
  CreateProductRequest,
  UpdateProductRequest,
  ReviewResponse,
  ListReviewsParams
} from '@/types/product'

// 商品管理 API
export function getProducts(params: ListProductParams): Promise<ListProductResponse> {
  return request.get('/admin/products', { params })
}

export function getProduct(id: number): Promise<Product> {
  return request.get(`/admin/products/${id}`)
}

export function createProduct(data: CreateProductRequest): Promise<Product> {
  return request.post('/admin/products', data)
}

export function updateProduct(id: number, data: UpdateProductRequest): Promise<Product> {
  return request.put(`/admin/products/${id}`, data)
}

export function deleteProduct(id: number): Promise<void> {
  return request.delete(`/admin/products/${id}`)
}

// 商品筛选 API
export function getProductsByCategory(category: string, params: ListProductParams): Promise<ListProductResponse> {
  return request.get(`/products/filter/category`, { 
    params: { ...params, category }
  })
}

export function getProductsByPrice(minPrice: number, maxPrice: number, params: ListProductParams): Promise<ListProductResponse> {
  return request.get(`/products/filter/price`, {
    params: { ...params, min_price: minPrice, max_price: maxPrice }
  })
}

export function getProductsByTags(tags: string[], params: ListProductParams): Promise<ListProductResponse> {
  return request.get(`/products/filter/tags`, {
    params: { ...params, tags: tags.join(',') }
  })
}

export function searchProducts(keyword: string, params: ListProductParams): Promise<ListProductResponse> {
  return request.get(`/products/filter/keyword`, {
    params: { ...params, keyword }
  })
}

// 商品评论 API
export function getProductReviews(productId: number, params: ListReviewsParams): Promise<ReviewResponse> {
  return request.get(`/admin/reviews/products/${productId}`, { params });
}

export function deleteProductReview(reviewId: number) {
  return request.delete(`/admin/reviews/reviews/${reviewId}`)
}

// 商品图片上传
export function uploadProductImage(file: File): Promise<{ url: string }> {
  const formData = new FormData()
  formData.append('file', file)
  return request.post('/admin/upload/product', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
} 