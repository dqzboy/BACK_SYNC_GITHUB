<template>
  <el-card class="ops-card">
    <template #header><span class="ops-card__title">任务历史</span></template>
    <el-table :data="jobs" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="server_name" label="服务器" width="170" />
      <el-table-column label="状态" width="130">
        <template #default="{ row }">
          <span class="ops-pill" :class="pillClass(row.status)">{{ statusText(row.status) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="message" label="结果" />
      <el-table-column prop="started_at" label="开始时间" width="180" />
      <el-table-column prop="finished_at" label="结束时间" width="180" />
      <el-table-column label="操作" width="100">
        <template #default="{ row }">
          <el-button link type="primary" @click="viewLog(row)">日志</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialog" title="备份日志" width="72%">
      <pre class="ops-log">{{ activeLog }}</pre>
    </el-dialog>
  </el-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const jobs = ref([])
const dialog = ref(false)
const activeLog = ref('')

function statusText(s) {
  return s === 'success' ? '成功' : s === 'failed' ? '失败' : s === 'running' ? '进行中' : '未知'
}
function pillClass(s) {
  if (s === 'success') return 'ops-pill--success'
  if (s === 'failed') return 'ops-pill--failed'
  if (s === 'running') return 'ops-pill--running'
  return 'ops-pill--idle'
}

async function load() {
  const { data } = await api.get('/jobs')
  jobs.value = data
}

function viewLog(row) {
  activeLog.value = row.log || '(无日志)'
  dialog.value = true
}

onMounted(load)
</script>
