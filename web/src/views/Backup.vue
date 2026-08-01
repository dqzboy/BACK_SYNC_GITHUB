<template>
  <el-card class="ops-card">
    <template #header>
      <div style="display: flex; justify-content: space-between; align-items: center">
        <span class="ops-card__title">执行备份</span>
        <el-button type="primary" :icon="VideoPlay" :loading="running" @click="run">开始备份</el-button>
      </div>
    </template>

    <el-alert
      v-if="!currentJob"
      type="info"
      title="点击「开始备份」将按当前配置，把备份源拷贝到以服务器 IP 命名的目录并提交推送到 GitHub 仓库"
      :closable="false"
    />

    <div v-else>
      <div class="ops-result" :class="'ops-result--' + currentJob.status">
        <div class="ops-result-icon">{{ resultIcon(currentJob.status) }}</div>
        <div>
          <div class="ops-result-title">{{ statusText(currentJob.status) }}</div>
          <div class="ops-result-sub">{{ currentJob.message || '任务进行中，请稍候…' }}</div>
          <div class="ops-result-meta">任务 #{{ currentJob.id }} · 服务器 {{ currentJob.server_name }}</div>
        </div>
      </div>
      <pre class="ops-log">{{ currentJob.log }}</pre>
    </div>
  </el-card>
</template>

<script setup>
import { ref, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { VideoPlay } from '@element-plus/icons-vue'
import api from '../api'

const running = ref(false)
const currentJob = ref(null)
let timer = null

function statusText(s) {
  return s === 'success' ? '备份成功' : s === 'failed' ? '备份失败' : s === 'running' ? '备份进行中' : '未知'
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
    ElMessage.error(e.response?.data?.error || '启动失败')
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
