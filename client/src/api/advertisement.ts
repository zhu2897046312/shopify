import request from '@/utils/request'
import type {
  Advertisement,
  ListAdvertisementParams,
  ListAdvertisementResponse,
  CreateAdvertisementRequest,
  UpdateAdvertisementRequest,
  UpdateAdvertisementStatusRequest
} from '@/types/advertisement'

// 获取广告列表
export function getAdvertisements(params: ListAdvertisementParams): Promise<ListAdvertisementResponse> {
  return request.get('/admin/advertisements', { params })
}

// 获取广告详情
export function getAdvertisement(id: number): Promise<Advertisement> {
  return request.get(`/admin/advertisements/${id}`)
}

// 创建广告
export function createAdvertisement(data: CreateAdvertisementRequest): Promise<Advertisement> {
  return request.post('/admin/advertisements', data)
}

// 更新广告
export function updateAdvertisement(id: number, data: UpdateAdvertisementRequest): Promise<Advertisement> {
  return request.put(`/admin/advertisements/${id}`, data)
}

// 删除广告
export function deleteAdvertisement(id: number): Promise<void> {
  return request.delete(`/admin/advertisements/${id}`)
}

// 更新广告状态
export function updateAdvertisementStatus(id: number, data: UpdateAdvertisementStatusRequest): Promise<Advertisement> {
  return request.put(`/admin/advertisements/${id}/status`, data)
} 