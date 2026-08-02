<template>
  <el-card class="ops-card">
    <template #header>
      <span class="ops-card__title"><el-icon class="ops-card__icon"><List /></el-icon>{{ t('jobs.title') }}</span>
    </template>
    <el-table :data="jobs" style="width: 100%" v-loading="loading" element-loading-background="transparent">
      <el-table-column prop="id" :label="t('jobs.id')" width="80" />
      <el-table-column prop="server_name" :label="t('jobs.server')" width="170" />
      <el-table-column :label="t('jobs.status')" width="130">
        <template #default="{ row }">
          <span class="ops-pill" :class="pillClass(row.status)">{{ statusText(row.status) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="message" :label="t('jobs.result')" />
      <el-table-column prop="started_at" :label="t('jobs.started')" width="180" />
      <el-table-column prop="finished_at" :label="t('jobs.finished')" width="180" />
      <el-table-column :label="t('jobs.actions')" width="100">
        <template #default="{ row }">
          <el-button link type="primary" @click="viewLog(row)">{{ t('jobs.log') }}</el-button>
        </template>
      </el-table-column>
      <template #empty>
        <EmptyState :title="t('jobs.empty')" />
      </template>
    </el-table>

    <el-dialog v-model="dialog" :title="t('jobs.logTitle')" width="72%">
      <pre class="ops-log">{{ activeLog }}</pre>
    </el-dialog>
  </el-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { List } from '@element-plus/icons-vue'
import EmptyState from '../components/EmptyState.vue'
import api from '../api'
import { t } from '../i18n'

const jobs = ref([])
const dialog = ref(false)
const activeLog = ref('')
const loading = ref(true)

function statusText(s) {
  return s === 'success' ? t('status.success') : s === 'failed' ? t('status.failed') : s === 'running' ? t('status.running') : t('status.unknown')
}
function pillClass(s) {
  if (s === 'success') return 'ops-pill--success'
  if (s === 'failed') return 'ops-pill--failed'
  if (s === 'running') return 'ops-pill--running'
  return 'ops-pill--idle'
}

async function load() {
  loading.value = true
  try {
    const { data } = await api.get('/jobs')
    jobs.value = data
  } finally {
    loading.value = false
  }
}

function viewLog(row) {
  activeLog.value = row.log || t('jobs.noLog')
  dialog.value = true
}

onMounted(load)
</script>
