// 主题管理：data-theme 属性 + Element Plus 的 .dark 类
const KEY = 'ops-theme'

export function applyTheme(theme) {
  const root = document.documentElement
  root.dataset.theme = theme
  root.classList.toggle('dark', theme === 'dark')
  try {
    localStorage.setItem(KEY, theme)
  } catch (e) {
    /* 忽略隐私模式下的写入失败 */
  }
}

export function getStoredTheme() {
  try {
    return localStorage.getItem(KEY)
  } catch (e) {
    return null
  }
}

export function initTheme() {
  let theme = getStoredTheme()
  if (!theme) {
    theme = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
  }
  applyTheme(theme)
  return theme
}

export function setTheme(theme) {
  applyTheme(theme)
  return theme
}

export function toggleTheme() {
  const cur = document.documentElement.dataset.theme === 'dark' ? 'dark' : 'light'
  return setTheme(cur === 'dark' ? 'light' : 'dark')
}
