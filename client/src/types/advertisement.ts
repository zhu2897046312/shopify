import { PaginationParams, PaginationResponse } from "./base"

// 广告状态枚举
export enum AdvertisementStatus {
  Active = 'active',
  Inactive = 'inactive'
}

// 广告位置枚举
export enum AdvertisementPosition {
  BannerTop = 'banner_top',      // 顶部横幅
  HomeBanner = 'home_banner',    // 首页轮播图
  HomeMiddle = 'home_middle',    // 首页中部广告位
  BannerBottom = 'banner_bottom', // 底部横幅
  SidebarTop = 'sidebar_top',    // 侧边栏顶部
  SidebarBottom = 'sidebar_bottom' // 侧边栏底部
}

// 广告数据模型
export interface Advertisement {
  id: number
  title: string
  image: string
  url: string
  position: AdvertisementPosition
  start_time: string
  end_time: string
  status: AdvertisementStatus
  created_at: string
  updated_at: string
}

// 创建广告请求参数
export interface CreateAdvertisementRequest {
  title: string
  image: string
  url: string
  position: string
  start_time: string
  end_time: string
  status: string
}

// 更新广告请求参数
export interface UpdateAdvertisementRequest {
  title?: string
  image?: string
  url?: string
  position?: string
  start_time?: string
  end_time?: string
  status?: string
}

// 更新广告状态请求参数
export interface UpdateAdvertisementStatusRequest {
  status: AdvertisementStatus
}

// 广告列表查询参数
export interface ListAdvertisementParams extends PaginationParams {
  position?: AdvertisementPosition
  status?: AdvertisementStatus
  start_time?: string
  end_time?: string
}

// 广告列表响应
export type ListAdvertisementResponse = PaginationResponse<Advertisement> 