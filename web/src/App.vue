<template>
  <el-container class="app-container">
    <el-header height="60px" class="app-header">
      <div class="logo">
        <h1>Easy Agent</h1>
      </div>
      <div class="user-info" v-if="userStore.isLoggedIn">
        <el-dropdown>
          <span class="user-dropdown">
            {{ userStore.username }}
            <el-icon><CaretBottom /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-header>
    <el-main class="app-main">
      <router-view></router-view>
    </el-main>
  </el-container>
</template>

<script setup>
import { useUserStore } from './stores/user'
import { useRouter } from 'vue-router'

const userStore = useUserStore()
const router = useRouter()

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.app-container {
  height: 100vh;
}

.app-header {
  background-color: #ffffff;
  border-bottom: 2px solid #529eb7;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10 10px;
  height: auto;
}

.logo h1 {
  margin: 0;
  font-size: 24px;
  color: #409eff;
}

.user-dropdown {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
}

.app-main {
  background-color: #f5f7fa;
  padding: 20px;
  overflow: hidden; /* 防止卡片出现滚动条 */
}
</style>