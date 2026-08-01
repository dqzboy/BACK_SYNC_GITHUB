<template>
  <el-container class="ops-shell">
    <el-aside width="250px" class="ops-side">
      <div class="ops-brand">
        <div class="ops-logo-mark">⎇</div>
        <div>
          <div class="ops-brand-name">GIT·BACKUP</div>
          <div class="ops-brand-sub">ops console</div>
        </div>
      </div>

      <el-menu :default-active="activeMenu" router class="ops-nav">
        <el-menu-item index="/dashboard">
          <el-icon><DataLine /></el-icon><span>仪表盘</span>
        </el-menu-item>
        <el-menu-item index="/config">
          <el-icon><Setting /></el-icon><span>备份配置</span>
        </el-menu-item>
        <el-menu-item index="/backup">
          <el-icon><Upload /></el-icon><span>执行备份</span>
        </el-menu-item>
        <el-menu-item index="/jobs">
          <el-icon><List /></el-icon><span>任务历史</span>
        </el-menu-item>
        <el-menu-item index="/users">
          <el-icon><User /></el-icon><span>用户中心</span>
        </el-menu-item>
        <el-menu-item index="/schedule">
          <el-icon><Clock /></el-icon><span>定时任务</span>
        </el-menu-item>
      </el-menu>

      <div class="ops-side-foot">
        <el-tooltip :content="`GitHub 仓库 · ${repoName}`" placement="top">
          <a class="ops-repo-icon" :href="repoUrl" target="_blank" rel="noopener" aria-label="查看 GitHub 仓库">
            <svg viewBox="0 0 16 16" width="18" height="18" fill="currentColor" aria-hidden="true">
              <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0 0 16 8c0-4.42-3.58-8-8-8z" />
            </svg>
          </a>
        </el-tooltip>
        <span class="ops-online"><span class="ops-dot"></span> SYSTEM ONLINE</span>
      </div>
    </el-aside>

    <el-container>
      <el-header class="ops-header">
        <div class="ops-crumb">~/{{ current }}</div>
        <div class="ops-header-right">
          <el-button
            class="ops-theme-btn"
            circle
            :icon="theme === 'dark' ? Sunny : Moon"
            aria-label="切换深浅主题"
            @click="onToggle"
          />
          <el-dropdown @command="onCommand">
            <span class="ops-user">
              <span class="ops-avatar">{{ initial }}</span>
              <span>{{ username }}</span>
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main class="ops-main"><router-view /></el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { DataLine, Setting, Upload, List, User, Clock, Sunny, Moon, ArrowDown } from '@element-plus/icons-vue'
import { getStoredTheme, setTheme, toggleTheme } from '../theme'
import { GITHUB_REPO_URL, GITHUB_REPO_NAME } from '../github'

const route = useRoute()
const router = useRouter()
const activeMenu = computed(() => route.path)
const current = computed(() => (route.path === '/' ? 'dashboard' : route.path.replace('/', '')))
const username = localStorage.getItem('username') || 'admin'
const initial = username.charAt(0).toUpperCase()
const theme = ref(getStoredTheme() || 'dark')

const repoUrl = GITHUB_REPO_URL
const repoName = GITHUB_REPO_NAME

function onToggle() {
  theme.value = toggleTheme()
}
function onCommand(cmd) {
  if (cmd === 'logout') {
    localStorage.removeItem('token')
    localStorage.removeItem('username')
    router.push('/login')
  }
}
</script>

<style scoped>
.ops-shell {
  height: 100vh;
}
.ops-side {
  background: var(--surface);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  padding: 20px 0;
}
.ops-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 20px 22px;
}
.ops-nav {
  flex: 1;
  background: transparent;
  border: none;
}
.ops-side-foot {
  margin-top: auto;
  padding: 14px 18px;
  border-top: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  gap: 10px;
  font-family: var(--font-display);
  font-size: 11px;
  letter-spacing: 0.04em;
}
.ops-repo-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  margin-left: -12px;
  border-radius: 10px;
  color: var(--text-muted);
  text-decoration: none;
  transition: color 0.2s ease, background-color 0.2s ease, border-color 0.2s ease;
  border: 1px solid transparent;
}
.ops-repo-icon:hover {
  color: var(--accent);
  background: var(--accent-dim);
  border-color: color-mix(in srgb, var(--accent) 35%, transparent);
}
.ops-online {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-muted);
  letter-spacing: 0.14em;
}
.ops-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--accent);
  box-shadow: 0 0 8px var(--accent-glow);
  animation: ops-pulse 1.6s infinite;
}
.ops-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: transparent;
  height: 64px;
  border-bottom: 1px solid var(--border);
  padding: 0 30px;
}
.ops-crumb {
  font-family: var(--font-display);
  color: var(--text-muted);
  font-size: 13px;
  letter-spacing: 0.05em;
}
.ops-header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}
.ops-theme-btn {
  background: var(--surface-2);
  border: 1px solid var(--border);
  color: var(--text-muted);
  transition: color 0.2s ease, border-color 0.2s ease, background-color 0.2s ease;
}
.ops-theme-btn:hover {
  color: var(--accent);
  border-color: var(--accent);
  background: var(--accent-dim);
}
.ops-user {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  color: var(--text);
  font-family: var(--font-display);
  font-size: 13px;
  outline: none;
}
.ops-avatar {
  width: 30px;
  height: 30px;
  border-radius: 8px;
  background: var(--accent-dim);
  color: var(--accent);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  border: 1px solid color-mix(in srgb, var(--accent) 30%, transparent);
}
.ops-main {
  padding: 30px;
  background: transparent;
}
</style>
