<template>
  <el-card class="ops-card">
    <template #header>
      <div class="ops-users-head">
        <span class="ops-card__title">用户中心</span>
        <el-button type="primary" :icon="Plus" @click="openCreate">新增用户</el-button>
      </div>
    </template>

    <el-table :data="users" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="username" label="用户名" />
      <el-table-column label="角色" width="150">
        <template #default="{ row }">
          <el-tag :type="row.role === 'admin' ? 'success' : 'info'" effect="dark">
            {{ row.role === 'admin' ? '管理员' : '观察者' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="200" />
      <el-table-column label="操作" width="160">
        <template #default="{ row }">
          <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button link type="danger" @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
      v-model="dialog"
      :title="editing ? '编辑用户' : '新增用户'"
      width="460px"
      append-to-body
      align-center
      destroy-on-close
    >
      <el-form :model="form" label-width="80px">
        <el-form-item label="用户名">
          <el-input v-model="form.username" :disabled="editing" placeholder="登录用户名" />
        </el-form-item>
        <el-form-item :label="editing ? '重置密码' : '密码'">
          <el-input
            v-model="form.password"
            type="password"
            show-password
            :placeholder="editing ? '留空则不修改' : '请输入密码'"
          />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="form.role" style="width: 100%">
            <el-option label="管理员" value="admin" />
            <el-option label="观察者" value="viewer" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialog = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="submit">保存</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import api from '../api'

const users = ref([])
const dialog = ref(false)
const editing = ref(false)
const saving = ref(false)
const form = ref({ id: 0, username: '', password: '', role: 'viewer' })

async function load() {
  const { data } = await api.get('/users')
  users.value = data
}

function openCreate() {
  editing.value = false
  form.value = { id: 0, username: '', password: '', role: 'viewer' }
  dialog.value = true
}
function openEdit(row) {
  editing.value = true
  form.value = { id: row.id, username: row.username, password: '', role: row.role }
  dialog.value = true
}

async function submit() {
  saving.value = true
  try {
    if (editing.value) {
      const payload = { role: form.value.role }
      if (form.value.password) payload.password = form.value.password
      await api.put('/users/' + form.value.id, payload)
    } else {
      await api.post('/users', {
        username: form.value.username,
        password: form.value.password,
        role: form.value.role
      })
    }
    ElMessage.success('已保存')
    dialog.value = false
    await load()
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '保存失败')
  } finally {
    saving.value = false
  }
}

async function remove(row) {
  try {
    await ElMessageBox.confirm(`确定删除用户「${row.username}」？`, '提示', { type: 'warning' })
  } catch (e) {
    return
  }
  try {
    await api.delete('/users/' + row.id)
    ElMessage.success('已删除')
    await load()
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '删除失败')
  }
}

onMounted(load)
</script>

<style scoped>
.ops-users-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>
