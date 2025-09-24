import axios from 'axios'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    // 只有在非登录/注册请求的401错误才跳转首页
    if (error.response?.status === 401) {
      const url = error.config?.url || ''
      const isAuthRequest = url.includes('/auth/login') || url.includes('/auth/register')
      
      if (!isAuthRequest) {
        // Token过期或无效，清除本地存储
        localStorage.removeItem('token')
        // 跳转到首页，让用户重新登录
        window.location.href = '/'
      }
    }
    return Promise.reject(error)
  }
)

export default api
