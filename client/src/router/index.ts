import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { UserRole } from '@/types/user'
import { ElMessage } from 'element-plus'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/user/login.vue'),
    meta: { guest: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/user/register.vue'),
    meta: { guest: true }
  },
  {
    path: '/',
    component: () => import('@/layouts/DefaultLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/dashboard'
      },
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue')
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/user/index.vue')
      },
      {
        path: 'users/stats',
        name: 'UserStats',
        component: () => import('@/views/user/stats.vue')
      },
      {
        path: 'products',
        name: 'Products',
        component: () => import('@/views/product/index.vue')
      },
      {
        path: 'products/create',
        name: 'CreateProduct',
        component: () => import('@/views/product/form.vue')
      },
      {
        path: 'products/:id/edit',
        name: 'EditProduct',
        component: () => import('@/views/product/form.vue')
      },
      {
        path: 'orders',
        name: 'Orders',
        component: () => import('@/views/order/index.vue')
      },
      {
        path: 'orders/stats',
        name: 'OrderStats',
        component: () => import('@/views/order/stats.vue')
      },
      {
        path: 'products/stats',
        name: 'ProductStats',
        component: () => import('@/views/product/stats.vue')
      },
      {
        path: 'advertisements/stats',
        name: 'AdvertisementStats',
        component: () => import('@/views/advertisement/stats.vue')
      },
      {
        path: 'orders/:id',
        name: 'OrderDetail',
        component: () => import('@/views/order/detail.vue'),
        meta: { title: '订单详情' }
      },
      {
        path: 'orders/:id/logistics',
        name: 'OrderLogistics',
        component: () => import('@/views/order/logistics.vue'),
        meta: { title: '物流信息' }
      },
      {
        path: 'reviews',
        name: 'Reviews',
        component: () => import('@/views/review/index.vue')
      },
      {
        path: 'advertisements',
        name: 'Advertisements',
        component: () => import('@/views/advertisement/index.vue'),
        meta: { title: '广告列表' }
      },
      {
        path: 'advertisements/create',
        name: 'CreateAdvertisement',
        component: () => import('@/views/advertisement/form.vue'),
        meta: { title: '新增广告' }
      },
      {
        path: 'advertisements/:id/edit',
        name: 'EditAdvertisement',
        component: () => import('@/views/advertisement/form.vue'),
        meta: { title: '编辑广告' }
      },
      {
        path: 'settings/profile',
        name: 'Profile',
        component: () => import('@/views/settings/profile.vue')
      },
      {
        path: 'settings/password',
        name: 'Password',
        component: () => import('@/views/settings/password.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth) {
    if (!authStore.token) {
      next('/login')
      return
    }
    
    try {
      if (!authStore.user) {
        await authStore.getInfo()
      }
      
      if (authStore.user && authStore.user.role !== UserRole.Admin) {
        ElMessage.error('非管理员账号无法访问')
        await authStore.logout()
        next('/login')
        return
      }
    } catch (error) {
      next('/login')
      return
    }
  }
  
  if (to.meta.guest && authStore.token) {
    next('/dashboard')
    return
  }
  
  next()
})

export default router