<template>
  <el-card class="ops-card">
    <template #header><span class="ops-card__title">定时任务</span></template>

    <el-form :model="form" label-width="140px" style="max-width: 760px">
      <div class="ops-section"><span>调度设置</span></div>

      <el-form-item label="启用定时备份">
        <el-switch v-model="form.schedule_enabled" @change="clearNext" />
      </el-form-item>

      <el-form-item label="执行频率">
        <el-select v-model="mode" style="width: 240px" @change="onModeChange">
          <el-option label="每天固定时间" value="daily" />
          <el-option label="每隔 N 小时" value="hourly" />
          <el-option label="自定义 Cron" value="custom" />
        </el-select>
      </el-form-item>

      <el-form-item v-if="mode === 'daily'" label="每日时间">
        <el-time-picker
          v-model="dailyTime"
          format="HH:mm"
          value-format="HH:mm"
          placeholder="选择时间"
          @change="rebuildCron"
        />
        <span class="ops-hint">每天在该时刻执行一次备份</span>
      </el-form-item>

      <el-form-item v-if="mode === 'hourly'" label="间隔(小时)">
        <el-input-number v-model="hourInterval" :min="1" :max="24" @change="rebuildCron" />
        <span class="ops-hint">每隔 {{ hourInterval }} 小时执行一次</span>
      </el-form-item>

      <el-form-item v-if="mode === 'custom'" label="Cron 表达式">
        <el-input v-model="form.schedule_cron" placeholder="分 时 日 月 周，例如 0 2 * * *" @input="clearNext" />
        <span class="ops-hint">标准 5 段 Cron：分 时 日 月 周</span>
      </el-form-item>

      <el-form-item label="下次执行">
        <el-tag v-if="nextRunText" type="success" effect="dark" class="ops-pill">{{ nextRunText }}</el-tag>
        <el-tag v-else type="info" effect="plain" class="ops-pill">—</el-tag>
      </el-form-item>

      <el-form-item label="上次自动运行">
        <span class="ops-muted">{{ form.schedule_last_run || '尚未执行' }}</span>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" :loading="saving" @click="save">保存设置</el-button>
        <el-button :loading="testing" @click="testRun">立即测试运行</el-button>
      </el-form-item>
    </el-form>

    <el-alert
      v-if="form.schedule_enabled"
      class="ops-alert"
      type="success"
      :closable="false"
      show-icon
      title="定时备份已开启"
      :description="`Cron: ${form.schedule_cron}　|　服务端每 20 秒同步一次配置，无需重启。`"
    />
  </el-card>
</template>

<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import api from '../api'

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
    ElMessage.error(e.response?.data?.error || '读取配置失败')
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
    ElMessage.success('定时设置已保存')
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '保存失败')
  } finally {
    saving.value = false
  }
}

async function testRun() {
  testing.value = true
  try {
    const { data } = await api.post('/backup/run')
    ElMessage.success(`已触发一次备份，任务 #${data.id}`)
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '触发失败')
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
</style>
