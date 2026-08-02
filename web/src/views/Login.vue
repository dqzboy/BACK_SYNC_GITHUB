<template>
  <div class="ops-login">
    <div class="ops-login-card ops-fade-up">
      <div class="ops-login-brand">
        <div class="ops-logo-mark">⎇</div>
        <div>
          <div class="ops-brand-name">{{ t('app.brand') }}</div>
          <div class="ops-brand-sub">{{ t('app.sub') }}</div>
        </div>
      </div>

      <div class="ops-prompt">&gt; auth.login --user<span class="caret">_</span></div>

      <el-form :model="form">
        <el-form-item>
          <el-input v-model="form.username" :placeholder="t('login.usernamePh')" :prefix-icon="User" size="large" />
        </el-form-item>
        <el-form-item>
          <el-input
            v-model="form.password"
            type="password"
            :placeholder="t('login.passwordPh')"
            :prefix-icon="Lock"
            show-password
            size="large"
            @keyup.enter="onSubmit"
          />
        </el-form-item>
        <el-button type="primary" size="large" :loading="loading" class="ops-login-btn" @click="onSubmit">
          {{ t('login.submit') }}
        </el-button>
      </el-form>

      <div class="ops-login-hint">{{ t('login.hint') }}</div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import api from '../api'
import { t } from '../i18n'

const router = useRouter()
const form = ref({ username: '', password: '' })
const loading = ref(false)

async function onSubmit() {
  loading.value = true
  try {
    const { data } = await api.post('/auth/login', form.value)
    localStorage.setItem('token', data.token)
    localStorage.setItem('username', form.value.username)
    router.push('/')
  } catch (e) {
    ElMessage.error(e.response?.data?.error || t('login.error'))
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.ops-login {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}
.ops-login-card {
  width: 384px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 18px;
  padding: 34px 30px;
  box-shadow: 0 40px 90px -36px rgba(0, 0, 0, 0.92), 0 0 0 1px rgba(52, 226, 155, 0.05);
}
.ops-login-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
}
.ops-prompt {
  font-family: var(--font-display);
  color: var(--accent);
  font-size: 13px;
  margin-bottom: 20px;
  letter-spacing: 0.04em;
}
.caret {
  animation: ops-pulse 1s steps(1) infinite;
}
.ops-login-btn {
  width: 100%;
  letter-spacing: 0.3em;
  margin-top: 4px;
}
.ops-login-hint {
  margin-top: 16px;
  text-align: center;
  font-family: var(--font-display);
  font-size: 11px;
  color: var(--text-muted);
  letter-spacing: 0.04em;
}
</style>
