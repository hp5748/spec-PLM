<template>
  <div class="login-container">
    <!-- 科技感背景动画 -->
    <div class="tech-background">
      <div class="grid-lines"></div>
      <div class="floating-particles">
        <div v-for="i in 20" :key="i" class="particle" :style="getParticleStyle(i)"></div>
      </div>
    </div>

    <!-- 登录卡片 -->
    <div class="login-box">
      <!-- 发光边框装饰 -->
      <div class="glow-border"></div>

      <div class="login-header">
        <div class="logo-wrapper">
          <div class="logo-icon">
            <img src="@/assets/logo.svg" alt="Logo" class="logo" />
          </div>
        </div>
        <h2 class="title neon-text">PLM 系统</h2>
        <p class="subtitle">产品生命周期管理平台</p>
      </div>

      <el-form ref="formRef" :model="loginForm" :rules="rules" class="login-form">
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
            prefix-icon="User"
            size="large"
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            prefix-icon="Lock"
            size="large"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>

        <el-form-item>
          <el-checkbox v-model="rememberMe">记住登录状态</el-checkbox>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            class="login-btn"
            @click="handleLogin"
          >
            <span v-if="!loading">登 录</span>
            <span v-else>登录中...</span>
          </el-button>
        </el-form-item>
      </el-form>

      <div class="login-footer">
        <div class="tech-divider">
          <span>默认管理员账号</span>
        </div>
        <p class="credential">admin / admin123</p>
      </div>
    </div>

    <!-- 底部版权信息 -->
    <div class="copyright">
      <span>PLM Product Lifecycle Management</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const formRef = ref<FormInstance>()
const loading = ref(false)
const rememberMe = ref(false)

const loginForm = reactive({
  username: '',
  password: '',
})

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' },
  ],
}

// 生成随机粒子样式
function getParticleStyle(index: number) {
  const size = Math.random() * 4 + 2
  const x = Math.random() * 100
  const y = Math.random() * 100
  const duration = Math.random() * 20 + 10
  const delay = Math.random() * 10

  return {
    width: `${size}px`,
    height: `${size}px`,
    left: `${x}%`,
    top: `${y}%`,
    animationDuration: `${duration}s`,
    animationDelay: `${delay}s`,
  }
}

async function handleLogin() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true

  try {
    await userStore.login(loginForm.username, loginForm.password)
    ElMessage.success('登录成功')

    // 跳转到重定向页面或首页
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (error: unknown) {
    const err = error as Error
    ElMessage.error(err.message || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
@use '@/styles/variables.scss' as *;

.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #0f172a 100%);
}

/* 科技感背景 */
.tech-background {
  position: absolute;
  inset: 0;
  overflow: hidden;

  .grid-lines {
    position: absolute;
    inset: 0;
    background-image:
      linear-gradient(rgba(0, 212, 255, 0.05) 1px, transparent 1px),
      linear-gradient(90deg, rgba(0, 212, 255, 0.05) 1px, transparent 1px);
    background-size: 60px 60px;
    animation: gridMove 20s linear infinite;
  }

  .floating-particles {
    position: absolute;
    inset: 0;

    .particle {
      position: absolute;
      background: $primary-color;
      border-radius: 50%;
      opacity: 0;
      animation: floatUp 15s ease-in-out infinite;
      box-shadow: 0 0 10px $primary-color;
    }
  }
}

@keyframes gridMove {
  0% {
    transform: translate(0, 0);
  }
  100% {
    transform: translate(60px, 60px);
  }
}

@keyframes floatUp {
  0% {
    opacity: 0;
    transform: translateY(0) scale(1);
  }
  10% {
    opacity: 0.8;
  }
  90% {
    opacity: 0.8;
  }
  100% {
    opacity: 0;
    transform: translateY(-100vh) scale(0.5);
  }
}

/* 登录卡片 */
.login-box {
  width: 420px;
  padding: 48px 40px;
  background: rgba(15, 23, 42, 0.85);
  border: 1px solid $border-color;
  border-radius: $border-radius-card;
  backdrop-filter: blur(20px);
  position: relative;
  z-index: 1;
  box-shadow:
    0 25px 50px -12px rgba(0, 0, 0, 0.5),
    0 0 0 1px rgba(0, 212, 255, 0.1);

  /* 发光边框 */
  .glow-border {
    position: absolute;
    inset: -1px;
    border-radius: inherit;
    background: linear-gradient(
      135deg,
      rgba(0, 212, 255, 0.4),
      transparent 30%,
      transparent 70%,
      rgba(14, 165, 233, 0.4)
    );
    z-index: -1;
    animation: borderGlow 3s ease-in-out infinite;
  }
}

@keyframes borderGlow {
  0%, 100% {
    opacity: 0.5;
  }
  50% {
    opacity: 1;
  }
}

.login-header {
  text-align: center;
  margin-bottom: 36px;

  .logo-wrapper {
    display: flex;
    justify-content: center;
    margin-bottom: 20px;
  }

  .logo-icon {
    width: 70px;
    height: 70px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: $primary-gradient;
    border-radius: 16px;
    box-shadow:
      0 0 30px rgba(0, 212, 255, 0.4),
      0 10px 30px rgba(0, 0, 0, 0.3);

    .logo {
      width: 42px;
      height: 42px;
      filter: brightness(0) invert(1);
    }
  }

  .title {
    font-size: 28px;
    font-weight: 700;
    margin: 0 0 12px;
    letter-spacing: 4px;
  }

  .subtitle {
    font-size: 14px;
    color: $text-secondary;
    margin: 0;
    letter-spacing: 2px;
  }
}

.login-form {
  .login-btn {
    width: 100%;
    height: 48px;
    font-size: 16px;
    font-weight: 600;
    letter-spacing: 4px;
    border-radius: $border-radius-base;
    background: $primary-gradient;
    border: none;
    box-shadow: 0 4px 20px rgba(0, 212, 255, 0.4);
    transition: all $transition-fast;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 25px rgba(0, 212, 255, 0.5);
    }

    &:active {
      transform: translateY(0);
    }
  }
}

.login-footer {
  margin-top: 28px;

  .tech-divider {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 12px;

    &::before,
    &::after {
      content: '';
      flex: 1;
      height: 1px;
      background: linear-gradient(90deg, transparent, $border-color, transparent);
    }

    span {
      font-size: 12px;
      color: $text-placeholder;
      white-space: nowrap;
    }
  }

  .credential {
    text-align: center;
    font-size: 13px;
    color: $text-secondary;
    margin: 0;
    font-family: 'Courier New', monospace;
    letter-spacing: 1px;
  }
}

.copyright {
  position: absolute;
  bottom: 24px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 12px;
  color: $text-placeholder;
  letter-spacing: 1px;

  span {
    opacity: 0.6;
  }
}
</style>
