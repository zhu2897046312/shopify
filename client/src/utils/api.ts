import axios from 'axios'

// 创建两个axios实例，分别用于本地和远程API
export const localApi = axios.create({
  baseURL: import.meta.env.VITE_LOCAL_API_URL,
  timeout: 5000
})

export const remoteApi = axios.create({
  baseURL: import.meta.env.VITE_REMOTE_API_URL,
  timeout: 5000
})

// 创建一个函数来尝试两个API
export async function tryBothApis(requestConfig: any) {
  try {
    // 先尝试本地API
    const localResponse = await localApi(requestConfig)
    return localResponse
  } catch (localError) {
    console.log('本地API请求失败，尝试远程API')
    try {
      // 本地失败后尝试远程API
      const remoteResponse = await remoteApi(requestConfig)
      return remoteResponse
    } catch (remoteError) {
      // 两个都失败则抛出错误
      throw remoteError
    }
  }
}

// 导出默认的API实例
export default {
  local: localApi,
  remote: remoteApi,
  tryBoth: tryBothApis
}
