<template>
  <div class="login-container">
    <div class="login-content">
      <div class="login-header">
        <img src="../img/logo.png" alt="Logo" class="logo">
        <h1>Easy Agent</h1>
        <p class="subtitle">智能助手，随时待命</p>
      </div>
      <el-card class="login-card">
        <template #header>
          <h2 class="card-header">登录</h2>
        </template>
        <el-form
          ref="loginForm"
          :model="loginData"
          :rules="rules"
          label-position="top"
        >
          <el-form-item label="用户名" prop="username">
            <el-input
              v-model="loginData.username"
              placeholder="请输入用户名"
              prefix-icon="User"
            />
          </el-form-item>
          <el-form-item label="密码" prop="password">
            <el-input
              v-model="loginData.password"
              type="password"
              placeholder="请输入密码"
              prefix-icon="Lock"
              show-password
            />
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              :loading="loading"
              class="login-button"
              @click="handleLogin"
            >
              登录
            </el-button>
          </el-form-item>
          <div class="register-link">
            还没有账号？
            <router-link to="/register">立即注册</router-link>
          </div>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const loginForm = ref(null)
const loading = ref(false)

const loginData = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginForm.value) return
  
  await loginForm.value.validate(async (valid) => {
    if (valid) {
      try {
        loading.value = true
        await userStore.login(loginData)
        ElMessage.success('登录成功')
        router.push('/')
      } catch (error) {
        ElMessage.error(error.response?.data?.error || '登录失败')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: flex-start;  /* 改为靠上对齐 */
  background: linear-gradient(135deg, #e0f2ff 0%, #f5f7fa 100%);
  padding: 20px;
  padding-top: 12vh;  /* 添加顶部间距 */
}

.login-header {
  text-align: center;
  margin-bottom: 16px;
}

.logo {
  width: 48px;  /* 进一步缩小 logo */
  height: 48px;
  margin-bottom: 8px;
}

.login-header h1 {
  font-size: 22px;
  margin: 0 0 2px;
}

.subtitle {
  font-size: 13px;
  margin: 0;
  color: #666;
}

.login-content {
  width: 100%;
  max-width: 360px;  /* 稍微减小宽度 */
  padding: 0 20px;
}

.login-card {
  width: 100%;
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.9);
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05) !important;
  padding: 8px 16px;  /* 增加内边距 */
  height: 44px;  /* 增加高度 */
}

:deep(.el-input__inner) {
  font-size: 15px;  /* 增加字体大小 */
}

.card-header {
  text-align: center;
  margin: 0;
  color: #303133;
  font-size: 20px;
  font-weight: 500;
}

:deep(.el-button) {
  height: 44px;
  font-size: 16px;
  border-radius: 8px;
  width: 100%;
  margin-top: 20px;
}

.register-link {
  text-align: center;
  margin-top: 20px;
  font-size: 14px;
  color: #606266;
}

.register-link a {
  color: #409eff;
  text-decoration: none;
  font-weight: 500;
  margin-left: 4px;
}

.register-link a:hover {
  color: #79bbff;
  text-decoration: underline;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  padding-bottom: 8px;
}
</style>