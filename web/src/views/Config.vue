<template>
  <el-card class="ops-card" v-loading="loading" element-loading-background="transparent">
    <template #header>
      <span class="ops-card__title"><el-icon class="ops-card__icon"><Setting /></el-icon>{{ t('config.title') }}</span>
    </template>
    <el-form :model="form" label-width="140px" style="max-width: 760px">
      <div class="ops-section"><span>{{ t('config.sectionRepo') }}</span></div>
      <el-form-item :label="t('config.gitUser')" required>
        <el-input v-model="form.git_user" />
      </el-form-item>
      <el-form-item :label="t('config.gitToken')" required>
        <el-input v-model="form.git_token" type="password" show-password :placeholder="t('config.tokenPh')" />
      </el-form-item>
      <el-form-item :label="t('config.repoName')" required>
        <el-input v-model="form.repo_name" />
      </el-form-item>
      <el-form-item :label="t('config.branch')">
        <el-input v-model="form.branch" />
      </el-form-item>

      <div class="ops-section"><span>{{ t('config.sectionBackup') }}</span></div>
      <el-form-item :label="t('config.backupDir')" required>
        <el-input v-model="form.backup_dir" />
      </el-form-item>
      <el-form-item :label="t('config.serverName')">
        <el-input v-model="form.server_name" :placeholder="t('config.serverNamePh')" />
      </el-form-item>
      <el-form-item :label="t('config.hostRoot')">
        <el-input v-model="form.host_root" :placeholder="t('config.hostRootPh')" />
        <div style="font-size: 12px; color: var(--el-text-color-secondary); margin-top: 6px; line-height: 1.5">{{ t('config.hostRootHelp') }}</div>
      </el-form-item>
      <el-form-item :label="t('config.sources')">
        <div style="width: 100%">
          <div v-for="(src, i) in form.backup_sources" :key="i" style="display: flex; margin-bottom: 10px">
            <el-input v-model="form.backup_sources[i]" :placeholder="t('config.sourcePh')" />
            <el-button type="danger" :icon="Delete" circle style="margin-left: 10px" @click="removeSource(i)" />
          </div>
          <el-button :icon="Plus" @click="addSource">{{ t('config.addPath') }}</el-button>
        </div>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" :loading="saving" @click="save">{{ t('config.save') }}</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Delete, Setting } from '@element-plus/icons-vue'
import api from '../api'
import { t } from '../i18n'

const form = ref({
  git_user: '',
  git_token: '',
  repo_name: '',
  branch: '',
  backup_dir: '',
  server_name: '',
  host_root: '',
  backup_sources: []
})
const saving = ref(false)
const loading = ref(true)

onMounted(async () => {
  try {
    const { data } = await api.get('/config')
    form.value = data
  } finally {
    loading.value = false
  }
})

function addSource() {
  form.value.backup_sources.push('')
}
function removeSource(i) {
  form.value.backup_sources.splice(i, 1)
}

async function save() {
  saving.value = true
  try {
    const payload = { ...form.value }
    payload.backup_sources = (payload.backup_sources || []).map((s) => s.trim()).filter(Boolean)
    await api.put('/config', payload)
    ElMessage.success(t('config.saved'))
  } catch (e) {
    ElMessage.error(e.response?.data?.error || t('config.saveError'))
  } finally {
    saving.value = false
  }
}
</script>
