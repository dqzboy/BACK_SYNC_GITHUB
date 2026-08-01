import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import Layout from '../layout/Layout.vue'
import Dashboard from '../views/Dashboard.vue'
import Config from '../views/Config.vue'
import Backup from '../views/Backup.vue'
import Jobs from '../views/Jobs.vue'
import Users from '../views/Users.vue'
import Schedule from '../views/Schedule.vue'

const routes = [
  { path: '/login', name: 'login', component: Login },
  {
    path: '/',
    component: Layout,
    meta: { requiresAuth: true },
    children: [
      { path: '', redirect: '/dashboard' },
      { path: 'dashboard', name: 'dashboard', component: Dashboard },
      { path: 'config', name: 'config', component: Config },
      { path: 'backup', name: 'backup', component: Backup },
      { path: 'jobs', name: 'jobs', component: Jobs },
      { path: 'users', name: 'users', component: Users },
      { path: 'schedule', name: 'schedule', component: Schedule }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth && !token) {
    return '/login'
  }
  if (to.path === '/login' && token) {
    return '/'
  }
})

export default router
