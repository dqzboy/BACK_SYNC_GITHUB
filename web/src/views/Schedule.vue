<template>
  <el-card class="ops-card" v-loading="loading" element-loading-background="transparent">
    <template #header>
      <span class="ops-card__title"><el-icon class="ops-card__icon"><Clock /></el-icon>{{ t('schedule.title') }}</span>
    </template>

    <el-form :model="form" label-width="140px" style="max-width: 760px">
      <div class="ops-section"><span>{{ t('schedule.section') }}</span></div>

      <el-form-item :label="t('schedule.enabled')">
        <el-switch v-model="form.schedule_enabled" @change="clearNext" />
      </el-form-item>

      <el-form-item :label="t('schedule.freq')">
        <el-select v-model="mode" style="width: 240px" @change="onModeChange">
          <el-option :label="t('schedule.daily')" value="daily" />
          <el-option :label="t('schedule.hourly')" value="hourly" />
          <el-option :label="t('schedule.custom')" value="custom" />
        </el-select>
      </el-form-item>

      <el-form-item v-if="mode === 'daily'" :label="t('schedule.dailyTime')">
        <el-time-picker
          v-model="dailyTime"
          format="HH:mm"
          value-format="HH:mm"
          :placeholder="t('schedule.pickTime')"
          @change="rebuildCron"
        />
        <span class="ops-hint">{{ t('schedule.dailyHint') }}</span>
      </el-form-item>

      <el-form-item v-if="mode === 'hourly'" :label="t('schedule.interval')">
        <el-input-number v-model="hourInterval" :min="1" :max="24" @change="rebuildCron" />
        <span class="ops-hint">{{ t('schedule.hourlyHint', { n: hourInterval }) }}</span>
      </el-form-item>

      <el-form-item v-if="mode === 'custom'" :label="t('schedule.cronExpr')">
        <el-input v-model="form.schedule_cron" :placeholder="t('schedule.cronPh')" @input="clearNext" />
        <span class="ops-hint">{{ t('schedule.cronHint') }}</span>
      </el-form-item>

      <el-form-item :label="t('schedule.next')">
        <transition name="ops-next" mode="out-in">
          <el-tag
            v-if="nextRunText"
            :key="'run'"
            class="ops-pill ops-pill--success ops-pill--live"
          >
            <el-icon class="ops-pill__icon"><Calendar /></el-icon>{{ nextRunText }}
          </el-tag>
          <el-tag v-else :key="'none'" class="ops-pill ops-pill--idle">{{ t('common.none') }}</el-tag>
        </transition>
      </el-form-item>

      <el-form-item :label="t('schedule.lastRun')">
        <span class="ops-muted">{{ form.schedule_last_run || t('schedule.neverRun') }}</span>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" :loading="saving" @click="save">{{ t('schedule.save') }}</el-button>
        <el-button :loading="testing" @click="testRun">{{ t('schedule.test') }}</el-button>
      </el-form-item>
    </el-form>

    <el-alert
      v-if="form.schedule_enabled"
      class="ops-alert"
      type="success"
      :closable="false"
      show-icon
      :title="t('schedule.on')"
      :description="t('schedule.onDesc', { cron: form.schedule_cron })"
    />
  </el-card>
</template>

<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Clock, Calendar } from '@element-plus/icons-vue'
import api from '../api'
import { t } from '../i18n'

const form = reactive({
  schedule_enabled: false,
  schedule_cron: '0 2 * * *',
  schedule_last_run: ''
})
const mode = ref('daily')
const dailyTime = ref('02:00')
const hourInterval = ref(2)
const saving = ref(false)
const testing = ref(false)
const loading = ref(true)

// 根据模式推导出 cron 表达式
function rebuildCron() {
  if (mode.value === 'daily') {
    const [h, m] = (dailyTime.value || '02:00').split(':')
    form.schedule_cron = `${m} ${h} * * *`
  } else if (mode.value === 'hourly') {
    form.schedule_cron = `0 */${hourInterval.value} * * *`
  }
  clearNext()
}

function onModeChange() {
  rebuildCron()
}

function clearNext() {
  // 触发下次执行预览重算
  nextRunTick.value++
}
const nextRunTick = ref(0)

// 由 cron 表达式计算下一次执行时间（支持 *、*/n、a-b、a,b）
function parseField(field, min, max) {
  const set = new Set()
  if (field === '*') {
    for (let i = min; i <= max; i++) set.add(i)
    return set
  }
  for (const part of field.split(',')) {
    let step = 1
    let base = part
    if (part.includes('/')) {
      const [b, s] = part.split('/')
      step = parseInt(s, 10)
      base = b === '*' ? '*' : b
    }
    if (base === '*') {
      for (let i = min; i <= max; i += step) set.add(i)
    } else if (base.includes('-')) {
      const [a, bb] = base.split('-').map(Number)
      for (let i = a; i <= bb; i += step) set.add(i)
    } else {
      set.add(parseInt(base, 10))
    }
  }
  return set
}

function nextRunFromCron(expr) {
  const f = expr.trim().split(/\s+/)
  if (f.length !== 5) return null
  const mins = parseField(f[0], 0, 59)
  const hrs = parseField(f[1], 0, 23)
  const doms = parseField(f[2], 1, 31)
  const mons = parseField(f[3], 1, 12)
  const dows = parseField(f[4], 0, 7)
  if (dows.has(7)) dows.add(0)
  const start = new Date(Date.now() + 60000)
  start.setSeconds(0, 0)
  for (let i = 0; i < 4 * 366 * 24 * 60; i++) {
    const d = new Date(start.getTime() + i * 60000)
    if (!mons.has(d.getMonth() + 1)) continue
    if (!doms.has(d.getDate())) continue
    if (!dows.has(d.getDay())) continue
    if (!hrs.has(d.getHours())) continue
    if (!mins.has(d.getMinutes())) continue
    return d
  }
  return null
}

const nextRunText = computed(() => {
  void nextRunTick.value
  if (!form.schedule_enabled || !form.schedule_cron) return ''
  const d = nextRunFromCron(form.schedule_cron)
  if (!d) return ''
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
})

// 根据已有 cron 反推当前模式
function inferMode(cron) {
  const f = (cron || '').trim().split(/\s+/)
  if (f.length !== 5) return 'custom'
  if (f[2] === '*' && f[3] === '*' && f[4] === '*') {
    if (f[1].startsWith('*/')) {
      hourInterval.value = parseInt(f[1].slice(2), 10) || 2
      return 'hourly'
    }
    if (/^\d+$/.test(f[1]) && /^\d+$/.test(f[0])) {
      dailyTime.value = `${pad(f[1])}:${pad(f[0])}`
      return 'daily'
    }
  }
  return 'custom'
}
function pad(n) {
  return String(n).padStart(2, '0')
}

onMounted(async () => {
  try {
    const { data } = await api.get('/config')
    form.schedule_enabled = data.schedule_enabled
    form.schedule_cron = data.schedule_cron
    form.schedule_last_run = data.schedule_last_run
    mode.value = inferMode(data.schedule_cron)
  } catch (e) {
    ElMessage.error(e.response?.data?.error || t('schedule.loadError'))
  } finally {
    loading.value = false
  }
})

// cron 变化时（自定义输入）刷新预览
watch(() => form.schedule_cron, () => clearNext())

async function save() {
  saving.value = true
  try {
    const payload = {
      schedule_enabled: form.schedule_enabled,
      schedule_cron: form.schedule_cron
    }
    await api.put('/config', payload)
    ElMessage.success(t('schedule.saved'))
  } catch (e) {
    ElMessage.error(e.response?.data?.error || t('schedule.saveError'))
  } finally {
    saving.value = false
  }
}

async function testRun() {
  testing.value = true
  try {
    const { data } = await api.post('/backup/run')
    ElMessage.success(t('schedule.triggered', { id: data.id }))
  } catch (e) {
    ElMessage.error(e.response?.data?.error || t('schedule.triggerError'))
  } finally {
    testing.value = false
  }
}
</script>

<style scoped>
.ops-hint {
  margin-left: 12px;
  color: var(--text-muted);
  font-size: 12px;
  font-family: var(--font-display);
}
.ops-muted {
  color: var(--text-muted);
  font-family: var(--font-display);
  font-size: 13px;
}
.ops-alert {
  margin-top: 8px;
}
.ops-pill__icon {
  margin-right: 2px;
  vertical-align: -1px;
}
/* subtle static cue that the value is "live" once scheduled backup is on */
.ops-pill--live {
  box-shadow: 0 0 0 1px color-mix(in srgb, var(--accent) 30%, transparent);
}
/* enter/exit animation for the next-run preview (low-frequency toggle) */
.ops-next-enter-active {
  transition: opacity 0.32s cubic-bezier(0.2, 0, 0, 1), transform 0.32s cubic-bezier(0.2, 0, 0, 1);
}
.ops-next-leave-active {
  transition: opacity 0.2s cubic-bezier(0.2, 0, 0, 1), transform 0.2s cubic-bezier(0.2, 0, 0, 1);
}
.ops-next-enter-from {
  opacity: 0;
  transform: translateY(-6px) scale(0.92);
}
.ops-next-leave-to {
  opacity: 0;
  transform: translateY(4px) scale(0.98);
}
</style>
