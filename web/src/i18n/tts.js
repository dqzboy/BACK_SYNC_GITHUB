// 语音播报（TTS）：基于浏览器原生 Web Speech API，无需任何依赖。
// 开启后，切换语言或路由时朗读当前界面关键文案，实现「多语种 + 语音」支持。
import { reactive } from 'vue'
import { currentLocale } from './index'

const SPEECH_LANG = { 'zh-CN': 'zh-CN', 'zh-TW': 'zh-TW', 'en': 'en-US' }
const KEY = 'ops-tts'

function readStored() {
  try {
    return localStorage.getItem(KEY) === '1'
  } catch (e) {
    return false
  }
}

const state = reactive({ enabled: readStored() })

export function isTTSEnabled() {
  return state.enabled
}

export function setTTSEnabled(v) {
  state.enabled = !!v
  try {
    localStorage.setItem(KEY, state.enabled ? '1' : '0')
  } catch (e) {}
  if (!state.enabled && typeof window !== 'undefined' && window.speechSynthesis) {
    window.speechSynthesis.cancel()
  }
}

export function speak(text, langCode) {
  if (typeof window === 'undefined' || !('speechSynthesis' in window)) return
  if (!text) return
  try {
    window.speechSynthesis.cancel()
    const u = new SpeechSynthesisUtterance(text)
    u.lang = SPEECH_LANG[langCode || currentLocale()] || 'zh-CN'
    u.rate = 0.95
    u.pitch = 1
    window.speechSynthesis.speak(u)
  } catch (e) {
    /* 部分浏览器在用户手势前会拒绝，忽略即可 */
  }
}

// 朗读（仅当播报开关开启时）
export function narrate(text, langCode) {
  if (!state.enabled) return
  speak(text, langCode)
}

export function useTTS() {
  return { state, setEnabled: setTTSEnabled, speak, narrate, isTTSEnabled }
}
