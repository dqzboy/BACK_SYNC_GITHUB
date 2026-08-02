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

    <!-- 数据可视化：状态分布 + 近期活动 -->
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :xs="24" :md="10">
        <el-card class="ops-card ops-fade-up" style="animation-delay: 0.28s">
          <template #header>
            <span class="ops-card__title"><el-icon class="ops-card__icon"><PieChart /></el-icon>{{ t('dashboard.distTitle') }}</span>
          </template>
          <div v-if="total" class="ops-donut-wrap">
            <svg class="ops-donut" viewBox="0 0 120 120">
              <circle cx="60" cy="60" r="52" fill="none" stroke="var(--surface-2)" stroke-width="14" />
              <circle
                v-for="(seg, i) in donutSegments"
                :key="i"
                cx="60" cy="60" r="52" fill="none"
                :stroke="seg.color" stroke-width="14" stroke-linecap="round"
                :stroke-dasharray="`${seg.dash} ${seg.gap}`"
                :stroke-dashoffset="seg.offset"
                transform="rotate(-90 60 60)"
              />
            </svg>
            <div class="ops-donut-center">
              <div class="ops-donut-total">{{ total }}</div>
              <div class="ops-donut-cap">{{ t('dashboard.totalJobs') }}</div>
            </div>
          </div>
          <EmptyState v-else :title="t('dashboard.noJobs')" />
          <div v-if="total" class="ops-legend">
            <span class="ops-legend__item"><i class="ops-legend__dot" style="background: var(--accent)" />{{ t('status.success') }} · {{ statusCounts.success }}</span>
            <span class="ops-legend__item"><i class="ops-legend__dot" style="background: var(--danger)" />{{ t('status.failed') }} · {{ statusCounts.failed }}</span>
            <span class="ops-legend__item"><i class="ops-legend__dot" style="background: var(--warn)" />{{ t('status.running') }} · {{ statusCounts.running }}</span>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :md="14">
        <el-card class="ops-card ops-fade-up" style="animation-delay: 0.34s">
          <template #header>
            <span class="ops-card__title"><el-icon class="ops-card__icon"><Histogram /></el-icon>{{ t('dashboard.recentTitle') }}</span>
          </template>
          <div v-if="recentBars.length" class="ops-bars">
            <div
              v-for="(b, i) in recentBars"
              :key="i"
              class="ops-bar"
              :style="{ height: b.h + '%', background: b.color }"
              :title="b.title"
            />
          </div>
          <EmptyState v-else :title="t('dashboard.noJobs')" />
        </el-card>
      </el-col>
    </el-row>

    <el-card class="ops-card ops-fade-up" style="margin-top: 16px; animation-delay: 0.4s">
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
import { ref, computed, onMounted } from 'vue'
import { DataLine, PieChart, Histogram } from '@element-plus/icons-vue'
import EmptyState from '../components/EmptyState.vue'
import api from '../api'
import { t } from '../i18n'

const cfg = ref({})
const jobs = ref([])
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
function colorOf(s) {
  if (s === 'success') return 'var(--accent)'
  if (s === 'failed') return 'var(--danger)'
  if (s === 'running') return 'var(--warn)'
  return 'var(--text-muted)'
}

const total = computed(() => jobs.value.length)
const statusCounts = computed(() => {
  const c = { success: 0, failed: 0, running: 0 }
  for (const j of jobs.value) {
    if (j.status === 'success') c.success++
    else if (j.status === 'failed') c.failed++
    else if (j.status === 'running') c.running++
  }
  return c
})

// 环形图分段
const donutSegments = computed(() => {
  const r = 52
  const C = 2 * Math.PI * r
  const order = [
    { key: 'success', color: 'var(--accent)' },
    { key: 'failed', color: 'var(--danger)' },
    { key: 'running', color: 'var(--warn)' }
  ]
  const segs = []
  let acc = 0
  for (const o of order) {
    const frac = total.value ? statusCounts.value[o.key] / total.value : 0
    const dash = frac * C
    segs.push({ color: o.color, dash, gap: C - dash, offset: -acc })
    acc += dash
  }
  return segs
})

// 近期活动迷你柱（最近 16 个任务），高度按状态区分，形成节奏
const recentBars = computed(() => {
  const list = jobs.value.slice(0, 16)
  return list.map((j) => {
    let h = 55
    if (j.status === 'success') h = 100
    else if (j.status === 'running') h = 86
    else if (j.status === 'failed') h = 68
    return { h, color: colorOf(j.status), title: `${j.id} · ${statusText(j.status)}` }
  })
})

onMounted(async () => {
  try {
    const c = await api.get('/config')
    cfg.value = c.data
    const j = await api.get('/jobs')
    jobs.value = j.data || []
    lastJob.value = jobs.value[0] || null
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.ops-donut-wrap {
  position: relative;
  width: 120px;
  height: 120px;
  margin: 6px auto 4px;
}
.ops-donut { width: 120px; height: 120px; }
.ops-donut-center {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  pointer-events: none;
}
.ops-donut-total {
  font-family: var(--font-display);
  font-size: 26px;
  font-weight: 700;
  color: var(--text);
  line-height: 1;
}
.ops-donut-cap {
  font-family: var(--font-display);
  font-size: 10px;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: var(--text-muted);
  margin-top: 4px;
}
.ops-legend {
  display: flex;
  flex-wrap: wrap;
  gap: 10px 18px;
  justify-content: center;
  margin-top: 12px;
  font-family: var(--font-display);
  font-size: 12px;
  color: var(--text-muted);
}
.ops-legend__item { display: inline-flex; align-items: center; gap: 7px; }
.ops-legend__dot { width: 9px; height: 9px; border-radius: 3px; display: inline-block; }
.ops-bars {
  display: flex;
  align-items: flex-end;
  gap: 6px;
  height: 132px;
  padding: 8px 2px 0;
}
.ops-bar {
  flex: 1;
  min-width: 6px;
  border-radius: 5px 5px 2px 2px;
  opacity: 0.85;
  transition: opacity 0.2s ease, transform 0.2s ease;
}
.ops-bar:hover { opacity: 1; transform: translateY(-2px); }
</style>
