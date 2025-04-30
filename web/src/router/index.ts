import { ElMessage } from 'element-plus'
import useUserStore from '@/stores/useUserStore'
import { createRouter, createWebHashHistory } from 'vue-router'
import routes from './routes'

const router = createRouter({
  history: createWebHashHistory(import.meta.env.VITE_BASE_PATH),
  routes: routes
})

router.beforeEach((to) => {
  const whitelist = ['/login', '/register']  
  const meta = to.meta
  const userStore = useUserStore()

  const isWhitelist = whitelist.includes(to.path)
  const hasRole = Array.isArray(meta.roles)
    ? meta.roles.includes(userStore.role)
    : true
  const hasToken = userStore.token !== ''

  if (!hasToken && !isWhitelist) {
    ElMessage.warning('请先登录后再操作')
    return '/login'
  } 

  if (hasToken && isWhitelist) {
    return '/home'
  }

  if (!hasRole) {
    return '/401'
  }

  return true
})

export default router
