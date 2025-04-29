import axios from 'axios'
import { ElMessage } from 'element-plus'
import useUserStore from '@/stores/useUserStore'

const http = axios.create({
  timeout: 30000,
  withCredentials: true
})

http.interceptors.request.use((config: any) => {
  const userStore = useUserStore()
  config.headers['token'] = userStore.token
  return config
})

http.interceptors.response.use((response: any) => {
  const { data } = response
  if (data.status !== 200) {
    // 只有当 message 存在且 hideError 不为 true 时才显示错误消息
    if (data.message && !data.hideError) {
      ElMessage.error(data.message)
    }
    return Promise.reject(data)
  }
  return data
})

export default http
