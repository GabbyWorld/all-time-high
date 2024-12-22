import React from 'react'
import axios, { AxiosRequestConfig, AxiosResponse } from 'axios'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 60000,
  httpAgent: true,
  httpsAgent: false,
  validateStatus: function (status: number) {
    return status >= 200
  }
})

request.interceptors.request.use((config: any) => {
  const Authorization = localStorage.getItem("Authorization")
  if (Authorization) {
    config.headers['Authorization'] = `Bearer ${Authorization}`
  } 
 
  return config
})


request.interceptors.response.use(
  (response: AxiosResponse) => {
    return response.data   
  },
  (error: any) => {
    return Promise.reject(error)
  }
)

interface RequestParams {
  [key: string]: any
}

interface Headers {
  [key: string]: string
}

export const get = <T>(url: string, params?: RequestParams): Promise<T> => {
  return request({
    method: 'GET',
    url,
    params
  })
}

export const deletes = <T>(url: string, params?: RequestParams): Promise<T> => {
  return request({
    method: 'DELETE',
    url,
    params
  })
}

export const post = <T>(url: string, data?: any, headers?: Headers): Promise<T> => {
  return request({
    method: 'POST',
    url,
    data,
    headers
  })
}

export const patch = <T>(url: string, data?: any, headers?: Headers): Promise<T> => {
  return request({
    method: 'PATCH',
    url,
    data,
    headers
  })
}
