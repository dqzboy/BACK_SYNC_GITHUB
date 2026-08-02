<template>
  <el-card class="ops-card">
    <template #header>
      <div style="display: flex; justify-content: space-between; align-items: center">
        <span class="ops-card__title"><el-icon class="ops-card__icon"><Upload /></el-icon>{{ t('backup.title') }}</span>
        <el-button type="primary" :icon="VideoPlay" :loading="running" @click="run">{{ t('backup.start') }}</el-button>
      </div>
    </template>

    <el-alert
      v-if="!currentJob"
      type="info"
      :title="t('backup.info')"
      :closable="false"
    />

    <div v-else>
      <div class="ops-result" :class="'ops-result--' + currentJob.status">
        <div class="ops-result-icon">{{ resultIcon(currentJob.status) }}</div>
        <div>
          <div class="ops-result-title">{{ statusText(currentJob.status) }}</div>
          <div class="ops-result-sub">{{ currentJob.message || t('backup.runningMsg') }}</div>
          <div class="ops-result-meta">{{ t('backup.jobMeta', { id: currentJob.id, server: currentJob.server_name }) }}</div>
        </div>
      </div>
      <pre class="ops-log">{{ currentJob.log }}</pre>
    </div>
  </el-card>
</template>

<script setup>
import { ref, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { VideoPlay, Upload } from '@element-plus/icons-vue'
import api from '../api'
import { t } from '../i18n'

const running = ref(false)
const currentJob = ref(null)
let timer = null

function statusText(s) {
  return s === 'success' ? t('backup.success') : s === 'failed' ? t('backup.failed') : s === 'running' ? t('backup.running') : t('backup.unknown')
}
function resultIcon(s) {
  return s === 'success' ? '✓' : s === 'failed' ? '✕' : '◌'
}

async function run() {
  running.value = true
  try {
    const { data } = await api.post('/backup/run')
    poll(data.id)
  } catch (e) {
    ElMessage.error(e.response?.data?.error || t('backup.startError'))
    running.value = false
  }
}

function poll(id) {
  timer = setInterval(async () => {
    try {
      const { data } = await api.get('/jobs/' + id)
      currentJob.value = data
      if (data.status !== 'running') {
        running.value = false
        clearInterval(timer)
        timer = null
      }
    } catch (e) {
      /* 忽略轮询错误 */
    }
  }, 1000)
}

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>
