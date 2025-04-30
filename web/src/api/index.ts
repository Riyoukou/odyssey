import axios from 'axios'
import { ElMessage } from 'element-plus'
import useUserStore from '@/stores/useUserStore'

const http = axios.create({
  timeout: 30000,
  withCredentials: true
})

http.interceptors.request.use(
  (config: any) => {
    const userStore = useUserStore()
    config.headers['token'] = userStore.token
    return config
  },
  (error: any) => {
    return Promise.reject(error)
  }
)

http.interceptors.response.use(
  (response: any) => {
    const { data } = response
    return data
  },
  (error) => {
    const { data } = error.response
    if (data.status === 401) {
      ElMessage.error("登录状态已过期")
      useUserStore().logoutApp()
      return Promise.reject(data)
    } else if (data.status === 404) {
      ElMessage.error("请求连接超时")
      return Promise.reject(data)
    } else if (data.status === 400) {
      ElMessage.error(`请求失败:${data.message} Code:${data.status}`)
      return Promise.reject(data)
    } else if (data.status !== 200) {
      ElMessage.error(`请求失败:${data.message} Code:${data.status}`)
      return Promise.reject(data)
    } else {
      return data
    }
  }
)

export default http
