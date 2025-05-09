<template>
  <div class="register-container">
    <div class="register-content">
      <div class="register-header">
        <img src="../img/logo.png" alt="Logo" class="logo">
        <h1>Easy Agent</h1>
        <p class="subtitle">开启智能助手之旅</p>
      </div>
      <el-card class="register-card">
        <template #header>
          <h2 class="card-header">账号注册</h2>
        </template>
        <el-form
          ref="registerForm"
          :model="registerData"
          :rules="rules"
          label-position="top"
        >
          <el-form-item label="用户名" prop="username">
            <el-input
              v-model="registerData.username"
              placeholder="请输入用户名"
              prefix-icon="User"
            />
          </el-form-item>
          <el-form-item label="密码" prop="password">
            <el-input
              v-model="registerData.password"
              type="password"
              placeholder="请输入密码"
              prefix-icon="Lock"
              show-password
            />
          </el-form-item>
          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input
              v-model="registerData.confirmPassword"
              type="password"
              placeholder="请再次输入密码"
              prefix-icon="Lock"
              show-password
            />
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              :loading="loading"
              class="register-button"
              @click="handleRegister"
            >
              注册
            </el-button>
          </el-form-item>
          <div class="login-link">
            已有账号？
            <router-link to="/login">立即登录</router-link>
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
const registerForm = ref(null)
const loading = ref(false)

const registerData = reactive({
  username: '',
  password: '',
  confirmPassword: ''
})

const validatePass2 = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== registerData.password) {
    callback(new Error('两次输入密码不一致!'))
  } else {
    callback()
  }
}

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validatePass2, trigger: 'blur' }
  ]
}

const handleRegister = async () => {
  if (!registerForm.value) return
  
  await registerForm.value.validate(async (valid) => {
    if (valid) {
      try {
        loading.value = true
        await userStore.register({
          username: registerData.username,
          password: registerData.password
        })
        ElMessage.success('注册成功')
        router.push('/login')
      } catch (error) {
        ElMessage.error(error.response?.data?.error || '注册失败')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  background: linear-gradient(135deg, #e0f2ff 0%, #f5f7fa 100%);
  padding: 20px;
  padding-top: 12vh;
}

.register-header {
  text-align: center;
  margin-bottom: 16px;
}

.logo {
  width: 48px;
  height: 48px;
  margin-bottom: 8px;
}

.register-header h1 {
  font-size: 22px;
  margin: 0 0 2px;
}

.subtitle {
  font-size: 13px;
  margin: 0;
  color: #666;
}

.register-content {
  width: 100%;
  max-width: 360px;
  padding: 0 20px;
}

/* 其他样式与 Login.vue 保持一致 */
</style>