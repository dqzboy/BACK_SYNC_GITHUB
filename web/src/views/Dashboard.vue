<template>
  <div>
    <el-row :gutter="16">
      <el-col :span="6">
        <el-card class="ops-card ops-fade-up" style="animation-delay: 0.04s">
          <div class="ops-stat-label">仓库</div>
          <div class="ops-stat-value">{{ cfg.repo_name || '-' }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="ops-card ops-fade-up" style="animation-delay: 0.1s">
          <div class="ops-stat-label">分支</div>
          <div class="ops-stat-value">{{ cfg.branch || '-' }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="ops-card ops-fade-up" style="animation-delay: 0.16s">
          <div class="ops-stat-label">备份源数量</div>
          <div class="ops-stat-value">{{ cfg.backup_sources?.length || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="ops-card ops-fade-up" style="animation-delay: 0.22s">
          <div class="ops-stat-label">最近任务</div>
          <div class="ops-stat-value">
            <span class="ops-pill" :class="pillClass(lastJob?.status)">
              {{ lastJob ? statusText(lastJob.status) : '无' }}
            </span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="ops-card" style="margin-top: 16px">
      <template #header><span class="ops-card__title">配置概览</span></template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="GitHub 用户">{{ cfg.git_user || '-' }}</el-descriptions-item>
        <el-descriptions-item label="仓库名">{{ cfg.repo_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="分支">{{ cfg.branch || '-' }}</el-descriptions-item>
        <el-descriptions-item label="备份目录">{{ cfg.backup_dir || '-' }}</el-descriptions-item>
        <el-descriptions-item label="服务器标识">{{ cfg.server_name || '自动探测(IP)' }}</el-descriptions-item>
        <el-descriptions-item label="管理员">{{ cfg.admin_user || '-' }}</el-descriptions-item>
        <el-descriptions-item label="备份源" :span="2">
          <el-tag v-for="s in cfg.backup_sources" :key="s" style="margin-right: 6px; margin-bottom: 4px">{{ s }}</el-tag>
          <span v-if="!cfg.backup_sources?.length" style="color: var(--text-muted)">未配置</span>
        </el-descriptions-item>
      </el-descriptions>

      <el-alert
        v-if="!cfg.git_token"
        type="warning"
        title="尚未配置 GitHub Token，请先到「备份配置」页面填写后再执行备份"
        :closable="false"
        style="margin-top: 14px"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const cfg = ref({})
const lastJob = ref(null)

function statusText(s) {
  return s === 'success' ? '成功' : s === 'failed' ? '失败' : s === 'running' ? '进行中' : '未知'
}
function pillClass(s) {
  if (s === 'success') return 'ops-pill--success'
  if (s === 'failed') return 'ops-pill--failed'
  if (s === 'running') return 'ops-pill--running'
  return 'ops-pill--idle'
}

onMounted(async () => {
  const c = await api.get('/config')
  cfg.value = c.data
  const j = await api.get('/jobs')
  lastJob.value = j.data[0] || null
})
</script>
