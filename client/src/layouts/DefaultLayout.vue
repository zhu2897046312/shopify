<template>
  <el-container class="layout-container">
    <!-- 侧边栏 -->
    <el-aside :width="isCollapse ? '64px' : '240px'" class="aside">
      <div class="logo">
        <img src="@/assets/logo.png" alt="Logo">
        <span v-if="!isCollapse">后台管理系统</span>
      </div>

      <div class="collapse-btn" @click="toggleCollapse">
        <el-icon :size="20">
          <Fold v-if="!isCollapse" />
          <Expand v-else />
        </el-icon>
      </div>
      
      <el-menu
        :default-active="route.path"
        :collapse="isCollapse"
        :unique-opened="true"
        router
        class="menu"
      >
        <el-menu-item index="/dashboard">
          <el-icon><DataBoard /></el-icon>
          <template #title>控制台</template>
        </el-menu-item>

        <!-- 用户管理 -->
        <el-sub-menu index="/users">
          <template #title>
            <el-icon><User /></el-icon>
            <span>用户管理</span>
          </template>
          <el-menu-item index="/users">用户列表</el-menu-item>
          <el-menu-item index="/users/stats">用户统计</el-menu-item>
        </el-sub-menu>

        <!-- 商品管理 -->
        <el-sub-menu index="/products">
          <template #title>
            <el-icon><Goods /></el-icon>
            <span>商品管理</span>
          </template>
          <el-menu-item index="/products">商品列表</el-menu-item>
          <el-menu-item index="/products/stats">商品统计</el-menu-item>
        </el-sub-menu>

        <!-- 订单管理 -->
        <el-sub-menu index="/orders">
          <template #title>
            <el-icon><List /></el-icon>
            <span>订单管理</span>
          </template>
          <el-menu-item index="/orders">订单列表</el-menu-item>
          <el-menu-item index="/orders/stats">订单统计</el-menu-item>
        </el-sub-menu>

        <!-- 广告管理 -->
        <el-sub-menu index="/advertisements">
          <template #title>
            <el-icon><Picture /></el-icon>
            <span>广告管理</span>
          </template>
          <el-menu-item index="/advertisements">广告列表</el-menu-item>
          <el-menu-item index="/advertisements/stats">广告统计</el-menu-item>
        </el-sub-menu>

        <!-- 系统设置 -->
        <el-sub-menu index="/settings">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>系统设置</span>
          </template>
          <el-menu-item index="/settings/profile">个人信息</el-menu-item>
          <el-menu-item index="/settings/password">修改密码</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>

    <!-- 主要内容区 -->
    <el-container>
      <el-header class="header">
        <div class="header-left">
          <el-breadcrumb>
            <el-breadcrumb-item v-for="item in breadcrumbs" :key="item.path" :to="item.path">
              {{ item.title }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-tooltip content="全屏" placement="bottom">
            <el-icon class="action-icon" @click="toggleFullscreen">
              <FullScreen v-if="!isFullscreen" />
              <Aim v-else />
            </el-icon>
          </el-tooltip>
          <user-avatar class="user-avatar" />
        </div>
      </el-header>

      <el-main class="main">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useFullscreen } from '@vueuse/core'
import {
  DataBoard,
  User,
  Goods,
  List,
  Picture,
  Setting,
  FullScreen,
  Aim,
  Fold,
  Expand
} from '@element-plus/icons-vue'
import UserAvatar from '@/components/user/UserAvatar.vue'

const route = useRoute()
const isCollapse = ref(false)
const { isFullscreen, toggle: toggleFullscreen } = useFullscreen()

const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

const breadcrumbs = computed(() => {
  return route.matched.map(item => ({
    title: item.meta.title,
    path: item.path
  }))
})
</script>

<style scoped lang="scss">
.layout-container {
  height: 100vh;

.aside {
  background: #e9ecf0;
  transition: width 0.3s;
  overflow: hidden;
  
  .logo {
    height: 60px;
    display: flex;
    align-items: center;
    padding: 0 20px;
    color: #fff;
    background: #e4eaef;
    
    img {
      height: 32px;
      margin-right: 12px;
    }
    
    span {
      font-size: 18px;
      font-weight: 600;
      white-space: nowrap;
    }
  }
  
  .collapse-btn {
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    color: #333;
    
    &:hover {
      background: #d3dce5;
    }
  }

  .menu {
    border-right: none;
    background-color: #e3e9ef;

    :deep(.el-menu-item) {
      &:hover {
        background-color: #1890ff !important;
        color: #fff !important;
      }
      
      &.is-active {
        background-color: #1890ff !important;
        color: #fff !important;
      }
    }

    :deep(.el-sub-menu) {
      .el-sub-menu__title {
        &:hover {
          background-color: #d3dce5 !important;
        }
      }
    }
  }
}

.header {
  background: #fff;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  box-shadow: 0 1px 4px rgba(0,21,41,.08);

  .header-left {
    display: flex;
    align-items: center;
    gap: 20px;
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 20px;

    .action-icon {
      font-size: 20px;
      cursor: pointer;
      color: #666;
      
      &:hover {
        color: #1890ff;
      }
    }

    .user-avatar {
      cursor: pointer;
    }
  }
}

  .main {
    background: #f0f2f5;
    padding: 20px;
  }
}

// 路由切换动画
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>