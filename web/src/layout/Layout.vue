<template>
  <el-container class="ops-shell">
    <!-- 移动端遮罩 -->
    <div v-if="isMobile && sidebarOpen" class="ops-backdrop" @click="sidebarOpen = false" />

    <el-aside
      width="250px"
      class="ops-side"
      :class="{ 'ops-side--mobile': isMobile, 'ops-side--open': isMobile && sidebarOpen }"
    >
      <div class="ops-brand">
        <div class="ops-logo-mark">⎇</div>
        <div>
          <div class="ops-brand-name">{{ t('app.brand') }}</div>
          <div class="ops-brand-sub">{{ t('app.sub') }}</div>
        </div>
      </div>

      <el-menu :default-active="activeMenu" router class="ops-nav">
        <el-menu-item index="/dashboard">
          <el-icon><DataLine /></el-icon><span>{{ t('nav.dashboard') }}</span>
        </el-menu-item>
        <el-menu-item index="/config">
          <el-icon><Setting /></el-icon><span>{{ t('nav.config') }}</span>
        </el-menu-item>
        <el-menu-item index="/backup">
          <el-icon><Upload /></el-icon><span>{{ t('nav.backup') }}</span>
        </el-menu-item>
        <el-menu-item index="/jobs">
          <el-icon><List /></el-icon><span>{{ t('nav.jobs') }}</span>
        </el-menu-item>
        <el-menu-item index="/users">
          <el-icon><User /></el-icon><span>{{ t('nav.users') }}</span>
        </el-menu-item>
        <el-menu-item index="/schedule">
          <el-icon><Clock /></el-icon><span>{{ t('nav.schedule') }}</span>
        </el-menu-item>
      </el-menu>

      <div class="ops-side-foot">
        <el-tooltip :content="t('app.repoTooltip', { repo: repoName })" placement="top">
          <a class="ops-repo-icon" :href="repoUrl" target="_blank" rel="noopener" :aria-label="t('app.repoTooltip', { repo: repoName })">
            <svg viewBox="0 0 16 16" width="18" height="18" fill="currentColor" aria-hidden="true">
              <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0 0 16 8c0-4.42-3.58-8-8-8z" />
            </svg>
          </a>
        </el-tooltip>
        <span class="ops-online"><span class="ops-dot"></span> {{ t('app.online') }}</span>
      </div>
    </el-aside>

    <el-container>
      <el-header class="ops-header">
        <div class="ops-head-left">
          <el-button
            v-if="isMobile"
            class="ops-menu-btn"
            circle
            :icon="Menu"
            :aria-label="t('responsive.menu')"
            @click="sidebarOpen = !sidebarOpen"
          />
          <div class="ops-crumb">
            <span class="ops-crumb-title">{{ pageTitle }}</span>
            <span class="ops-crumb-path">~/{{ current }}</span>
          </div>
        </div>

        <div class="ops-header-right">
          <!-- 语言切换 -->
          <el-dropdown trigger="click" @command="onLang">
            <button class="ops-lang-trigger" :aria-label="t('lang.label')">
              <span class="ops-lang-text">{{ currentLangLabel }}</span>
              <el-icon class="ops-lang-caret"><ArrowDown /></el-icon>
            </button>
            <template #dropdown>
              <el-dropdown-menu class="ops-lang-menu">
                <div class="ops-lang-menu-head">{{ t('lang.label') }}</div>
                <el-dropdown-item
                  v-for="l in availableLocales"
                  :key="l.code"
                  :command="l.code"
                  :class="{ 'is-active': locale === l.code }"
                >
                  <span class="ops-lang-item">
                    <el-icon v-if="locale === l.code" class="ops-lang-check"><Select /></el-icon>
                    <span v-else class="ops-lang-check ops-lang-check--ph" />
                    <span class="ops-lang-name">{{ l.label }}</span>
                  </span>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>

          <!-- 主题切换 -->
          <el-button
            class="ops-theme-btn"
            circle
            :icon="theme === 'dark' ? Sunny : Moon"
            aria-label="theme"
            @click="onToggle"
          />

          <el-dropdown @command="onCommand">
            <span class="ops-user">
              <span class="ops-avatar">{{ initial }}</span>
              <span class="ops-user-name">{{ username }}</span>
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">{{ t('common.logout') }}</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main class="ops-main">
        <router-view v-slot="{ Component }">
          <transition name="ops-route" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed, ref, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  DataLine, Setting, Upload, List, User, Clock, Sunny, Moon, ArrowDown, Select, Menu
} from '@element-plus/icons-vue'
import { getStoredTheme, toggleTheme } from '../theme'
import { GITHUB_REPO_URL, GITHUB_REPO_NAME } from '../github'
import { useI18n, t } from '../i18n'

const route = useRoute()
const router = useRouter()
const { locale, setLocale, availableLocales } = useI18n()

const activeMenu = computed(() => route.path)
const current = computed(() => (route.path === '/' ? 'dashboard' : route.path.replace('/', '')))
const pageTitle = computed(() => t(route.meta.title || 'app.name'))
const username = localStorage.getItem('username') || 'admin'
const initial = username.charAt(0).toUpperCase()
const theme = ref(getStoredTheme() || 'dark')

const repoUrl = GITHUB_REPO_URL
const repoName = GITHUB_REPO_NAME

// 响应式：移动端折叠侧边栏
const isMobile = ref(false)
const sidebarOpen = ref(false)
let mq = null
function applyMQ(e) {
  isMobile.value = e.matches
  if (!e.matches) sidebarOpen.value = false
}

// 当前语言短标签 / 旗帜字符
const langMeta = {
  'zh-CN': { label: '简体', flag: '🇨🇳' },
  'zh-TW': { label: '繁體', flag: '🇹🇼' },
  en: { label: 'EN', flag: '🇬🇧' }
}
const currentLangLabel = computed(() => (langMeta[locale.value] || langMeta['zh-CN']).label)
const currentLangFlag = computed(() => (langMeta[locale.value] || langMeta['zh-CN']).flag)

function onLang(code) {
  setLocale(code)
}
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

onMounted(() => {
  if (typeof window !== 'undefined' && window.matchMedia) {
    mq = window.matchMedia('(max-width: 860px)')
    applyMQ(mq)
    mq.addEventListener('change', applyMQ)
  }
})
onBeforeUnmount(() => {
  if (mq) mq.removeEventListener('change', applyMQ)
})
</script>

<style>
/* ============================ 布局骨架 ============================ */
.ops-shell { height: 100vh; }
.ops-side {
  background: var(--surface);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  padding: 20px 0;
  transition: transform 0.28s ease;
  z-index: 30;
}
.ops-brand { display: flex; align-items: center; gap: 12px; padding: 0 20px 22px; }
.ops-nav { flex: 1; background: transparent; border: none; }
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
.ops-online { display: flex; align-items: center; gap: 8px; color: var(--text-muted); letter-spacing: 0.14em; }
.ops-dot {
  width: 8px; height: 8px; border-radius: 50%;
  background: var(--accent);
  box-shadow: 0 0 8px var(--accent-glow);
  animation: ops-pulse 1.6s infinite;
}
.ops-header {
  display: flex; align-items: center; justify-content: space-between;
  background: transparent; height: 64px; border-bottom: 1px solid var(--border); padding: 0 30px;
}
.ops-head-left { display: flex; align-items: center; gap: 14px; min-width: 0; }
.ops-crumb { display: flex; flex-direction: column; line-height: 1.15; min-width: 0; }
.ops-crumb-title {
  font-family: var(--font-display); color: var(--text); font-size: 16px; font-weight: 600;
  letter-spacing: 0.02em; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}
.ops-crumb-path { font-family: var(--font-display); color: var(--text-muted); font-size: 11px; letter-spacing: 0.08em; }
.ops-header-right { display: flex; align-items: center; gap: 12px; }
.ops-theme-btn {
  background: var(--surface-2); border: 1px solid var(--border); color: var(--text-muted);
  transition: color 0.2s ease, border-color 0.2s ease, background-color 0.2s ease;
}
.ops-theme-btn:hover { color: var(--accent); border-color: var(--accent); background: var(--accent-dim); }
.ops-user {
  display: flex; align-items: center; gap: 10px; cursor: pointer; color: var(--text);
  font-family: var(--font-display); font-size: 13px; outline: none;
}
.ops-user-name { max-width: 120px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.ops-avatar {
  width: 30px; height: 30px; border-radius: 8px; background: var(--accent-dim); color: var(--accent);
  display: flex; align-items: center; justify-content: center; font-weight: 700;
  border: 1px solid color-mix(in srgb, var(--accent) 30%, transparent);
}
.ops-main { padding: 30px; background: transparent; overflow-x: hidden; }

/* ============================ 语言切换 ============================ */
.ops-lang-trigger {
  display: inline-flex; align-items: center; gap: 7px; height: 36px; padding: 0 12px;
  border-radius: 999px; background: var(--surface-2); border: 1px solid var(--border);
  color: var(--text); font-family: var(--font-display); font-size: 13px; cursor: pointer;
  transition: color 0.2s ease, border-color 0.2s ease, background-color 0.2s ease;
}
.ops-lang-trigger:hover { color: var(--accent); border-color: var(--accent); background: var(--accent-dim); }
.ops-lang-flag { font-size: 15px; line-height: 1; }
.ops-lang-caret { font-size: 12px; color: var(--text-muted); }
.ops-lang-menu { padding: 8px !important; min-width: 184px; }
.ops-lang-menu-head {
  font-family: var(--font-display); font-size: 10px; letter-spacing: 0.18em; text-transform: uppercase;
  color: var(--text-muted); padding: 8px 12px 10px;
}
.ops-lang-menu .el-dropdown-menu__item { border-radius: 8px; padding: 10px 12px; }
.ops-lang-menu .el-dropdown-menu__item.is-active { background: var(--accent-dim) !important; color: var(--accent) !important; }
.ops-lang-item { display: flex; align-items: center; gap: 8px; width: 100%; }
.ops-lang-check { font-size: 14px; color: var(--accent); width: 16px; }
.ops-lang-check--ph { visibility: hidden; }
.ops-lang-name { flex: 1; font-family: var(--font-display); }

/* ============================ 移动端 ============================ */
.ops-menu-btn { background: var(--surface-2); border: 1px solid var(--border); color: var(--text-muted); }
.ops-menu-btn:hover { color: var(--accent); border-color: var(--accent); }
.ops-backdrop {
  position: fixed; inset: 0; background: rgba(2, 6, 23, 0.55); z-index: 25; backdrop-filter: blur(2px);
}
@media (max-width: 860px) {
  .ops-side {
    position: fixed; top: 0; bottom: 0; left: 0; width: 250px;
    transform: translateX(-100%); box-shadow: 0 30px 80px -20px rgba(0, 0, 0, 0.6);
  }
  .ops-side.ops-side--open { transform: translateX(0); }
  .ops-header { padding: 0 16px; }
  .ops-main { padding: 16px; }
  .ops-user-name { display: none; }
}

/* ============================ 路由切换动画 ============================ */
.ops-route-enter-active, .ops-route-leave-active { transition: opacity 0.22s ease, transform 0.22s ease; }
.ops-route-enter-from { opacity: 0; transform: translateY(8px); }
.ops-route-leave-to { opacity: 0; transform: translateY(-8px); }
@media (prefers-reduced-motion: reduce) {
  .ops-route-enter-active, .ops-route-leave-active { transition: none; }
  .ops-backdrop { backdrop-filter: none; }
}
</style>