<template>
  <el-card class="ops-card">
    <template #header><span class="ops-card__title">备份配置</span></template>
    <el-form :model="form" label-width="140px" style="max-width: 760px">
      <div class="ops-section"><span>GitHub 仓库</span></div>
      <el-form-item label="GitHub 用户名">
        <el-input v-model="form.git_user" />
      </el-form-item>
      <el-form-item label="GitHub Token">
        <el-input v-model="form.git_token" type="password" show-password placeholder="留空表示不修改（已设置时显示为 ********）" />
      </el-form-item>
      <el-form-item label="仓库名称">
        <el-input v-model="form.repo_name" />
      </el-form-item>
      <el-form-item label="分支">
        <el-input v-model="form.branch" />
      </el-form-item>

      <div class="ops-section"><span>备份设置</span></div>
      <el-form-item label="备份目录">
        <el-input v-model="form.backup_dir" />
      </el-form-item>
      <el-form-item label="服务器标识">
        <el-input v-model="form.server_name" placeholder="留空则自动探测本机 IP" />
      </el-form-item>
      <el-form-item label="备份源路径">
        <div style="width: 100%">
          <div v-for="(src, i) in form.backup_sources" :key="i" style="display: flex; margin-bottom: 10px">
            <el-input v-model="form.backup_sources[i]" placeholder="/path/to/file_or_dir" />
            <el-button type="danger" :icon="Delete" circle style="margin-left: 10px" @click="removeSource(i)" />
          </div>
          <el-button :icon="Plus" @click="addSource">添加路径</el-button>
        </div>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" :loading="saving" @click="save">保存配置</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Delete } from '@element-plus/icons-vue'
import api from '../api'

const form = ref({
  git_user: '',
  git_token: '',
  repo_name: '',
  branch: '',
  backup_dir: '',
  server_name: '',
  backup_sources: []
})
const saving = ref(false)

onMounted(async () => {
  const { data } = await api.get('/config')
  form.value = data
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
    ElMessage.success('配置已保存')
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '保存失败')
  } finally {
    saving.value = false
  }
}
</script>
