// 轻量 i18n：响应式 locale、localStorage 持久化、浏览器语言探测、缺失回退到简体中文。
// 不依赖第三方库，便于在 Docker 构建中零额外网络依赖集成。
import { reactive, computed } from 'vue'
import zhCN from './locales/zh-CN'
import zhTW from './locales/zh-TW'
import en from './locales/en'

const messages = { 'zh-CN': zhCN, 'zh-TW': zhTW, 'en': en }
const STORAGE_KEY = 'ops-lang'

export const availableLocales = [
  { code: 'zh-CN', label: '简体中文' },
  { code: 'zh-TW', label: '繁體中文' },
  { code: 'en', label: 'English' }
]

function lookup(dict, key) {
  if (!dict) return undefined
  return key.split('.').reduce((o, k) => (o == null ? undefined : o[k]), dict)
}

function resolveInitial() {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved && messages[saved]) return saved
  } catch (e) {
    /* 隐私模式下 localStorage 不可用 */
  }
  const nav = (typeof navigator !== 'undefined' && navigator.language ? navigator.language : 'zh-CN').toLowerCase()
  if (nav.startsWith('zh')) {
    if (nav.includes('tw') || nav.includes('hk') || nav.includes('hant')) return 'zh-TW'
    return 'zh-CN'
  }
  return 'en'
}

const state = reactive({ locale: resolveInitial() })

// 初始化 <html lang>
try {
  document.documentElement.setAttribute('lang', state.locale)
} catch (e) {}

export function t(key, params) {
  const locale = state.locale
  let str = lookup(messages[locale], key)
  if (str == null) str = lookup(messages['zh-CN'], key)
  if (str == null) return key
  if (params && typeof str === 'string') {
    str = str.replace(/\{(\w+)\}/g, (m, n) => (params[n] != null ? params[n] : m))
  }
  return str
}

export function setLocale(code) {
  if (!messages[code]) return
  state.locale = code
  try {
    localStorage.setItem(STORAGE_KEY, code)
    document.documentElement.setAttribute('lang', code)
  } catch (e) {}
}

export function currentLocale() {
  return state.locale
}

export function useI18n() {
  return {
    t,
    locale: computed(() => state.locale),
    setLocale,
    availableLocales
  }
}
