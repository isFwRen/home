import Vue from 'vue'
import axios from 'axios'
import { localStorage } from 'vue-rocket'
import { tools as lpTools } from '@/libs/util'

const baseURLApi = 'http://127.0.0.1:30080'

/**
 * @description 创建请求实例
 */
function createService() {
	// 创建一个 axios 实例
	const service = axios.create()

	// 请求拦截
	service.interceptors.request.use(
		config => {
			return config
		},
		error => {
			return Promise.reject(error)
		}
	)

	// 响应拦截
	service.interceptors.response.use(
		response => {
			const result = response.data

			return result
		},
		error => {
			if (error.response === undefined) {
				console.error(error.message)

				error.response = { status: 502 }
			}

			const status = error.response.status
			switch (status) {
				case 400: error.message = '请求错误'; break
				case 401: error.message = '未授权，请登录'; break
				case 403: error.message = '拒绝访问'; break
				case 404: error.message = '请求地址不存在'; break
				case 408: error.message = '请求超时'; break
				case 500: error.message = '服务器内部错误'; break
				case 501: error.message = '服务未实现'; break
				case 502: error.message = '网关错误'; break
				case 503: error.message = '服务不可用'; break
				case 504: error.message = '网关超时'; break
				case 505: error.message = 'HTTP版本不受支持'; break
				default: break
			}

			Vue.prototype.toasted.error(error.message)

			return Promise.reject(error)
		}
	)

	return service
}

/**
 * @description 创建请求方法
 * @param {Object} service axios 实例
 */
function createRequest(service) {
	return function (config) {
		const token = localStorage.get('token')
		const user = localStorage.get('user')
		const project = localStorage.get('project') || {}
		const secret = localStorage.get("secret") || '';

		const headers = {
			'x-token': token,
			'Content-Type': 'application/json;charset=UTF-8'
		}
		if (secret) {
			const code = lpTools.GetCode(secret)
			headers["x-code"] = String(code);
		}
		if (user) {
			headers['x-user-id'] = user.id
		}

		if (project.code) {
			headers['pro-code'] = project.code
		}

		const defaultConfig = {
			method: 'GET',
			baseURL: baseURLApi,
			headers
		}

		const assignConfig = Object.assign(defaultConfig, config)

		if (config.headers) {
			assignConfig.headers = {
				...headers,
				...config.headers
			}
		}

		return service(assignConfig)
	}
}

export const service = createService()
export const reqConst = createRequest(service)










