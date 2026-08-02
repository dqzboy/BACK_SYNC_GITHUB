<template>
  <div v-loading="loading" element-loading-background="transparent">
    <el-row :gutter="16">
      <el-col :xs="12" :sm="12" :md="6">
        <el-card class="ops-card ops-fade-up" style="animation-delay: 0.04s">
          <div class="ops-stat-label">{{ t('dashboard.repo') }}</div>
          <div class="ops-stat-value">{{ cfg.repo_name || '-' }}</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="12" :md="6">
        <el-card class="ops-card ops-fade-up" style="animation-delay: 0.1s">
          <div class="ops-stat-label">{{ t('dashboard.branch') }}</div>
          <div class="ops-stat-value">{{ cfg.branch || '-' }}</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="12" :md="6">
        <el-card class="ops-card ops-fade-up" style="animation-delay: 0.16s">
          <div class="ops-stat-label">{{ t('dashboard.sources') }}</div>
          <div class="ops-stat-value">{{ cfg.backup_sources?.length || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="12" :md="6">
        <el-card class="ops-card ops-fade-up" style="animation-delay: 0.22s">
          <div class="ops-stat-label">{{ t('dashboard.lastJob') }}</div>
          <div class="ops-stat-value">
            <span class="ops-pill" :class="pillClass(lastJob?.status)">
              {{ lastJob ? statusText(lastJob.status) : t('common.none') }}
            </span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="ops-card" style="margin-top: 16px">
      <template #header>
        <span class="ops-card__title"><el-icon class="ops-card__icon"><DataLine /></el-icon>{{ t('dashboard.overview') }}</span>
      </template>
      <el-descriptions :column="2" border>
        <el-descriptions-item :label="t('dashboard.gitUser')">{{ cfg.git_user || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="t('dashboard.repoName')">{{ cfg.repo_name || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="t('dashboard.branchL')">{{ cfg.branch || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="t('dashboard.backupDir')">{{ cfg.backup_dir || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="t('dashboard.serverId')">{{ cfg.server_name || t('dashboard.autoIp') }}</el-descriptions-item>
        <el-descriptions-item :label="t('dashboard.admin')">{{ cfg.admin_user || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="t('dashboard.sourcesL')" :span="2">
          <el-tag v-for="s in cfg.backup_sources" :key="s" style="margin-right: 6px; margin-bottom: 4px">{{ s }}</el-tag>
          <span v-if="!cfg.backup_sources?.length" style="color: var(--text-muted)">{{ t('dashboard.notConfigured') }}</span>
        </el-descriptions-item>
      </el-descriptions>

      <el-alert
        v-if="!cfg.git_token"
        type="warning"
        :title="t('dashboard.tokenWarn')"
        :closable="false"
        style="margin-top: 14px"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { DataLine } from '@element-plus/icons-vue'
import api from '../api'
import { t } from '../i18n'

const cfg = ref({})
const lastJob = ref(null)
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

onMounted(async () => {
  try {
    const c = await api.get('/config')
    cfg.value = c.data
    const j = await api.get('/jobs')
    lastJob.value = j.data[0] || null
  } finally {
    loading.value = false
  }
})
</script>
