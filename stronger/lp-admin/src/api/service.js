import Vue from "vue";
import axios from "axios";
import { tools as lpTools, storage } from "@/libs/util";
import { localStorage } from "vue-rocket";

const { baseURLApi } = lpTools.baseURL();

/**
 * @description 创建请求实例
 */
function createService() {
	// 创建一个 axios 实例
	const service = axios.create();

	// 请求拦截
	service.interceptors.request.use(
		config => {
			return config;
		},
		error => {
			return Promise.reject(error);
		}
	);

	// 响应拦截
	service.interceptors.response.use(
		response => {
			const result = response.data;

			return result;
		},
		error => {
			if (error.response === undefined) {
				console.error(error.message);
				error.response = { status: 502 };
				return error;
			}

			const status = error.response.status;
			switch (status) {
				// case 400:
				// 	error.message = "请求错误";
				// 	break;
				case 401:
					// Vue.prototype.$router.push({path:"/login"});

					error.message = "登录过期或者校准系统时间后重新登陆";
					break;
				case 403:
					error.message = "拒绝访问";
					break;
				case 404:
					error.message = "请求地址不存在";
					break;
				case 408:
					error.message = "请求超时";
					break;
				// case 500:
				// 	error.message = "服务器内部错误";
				// 	break;
				case 501:
					error.message = "服务未实现";
					break;
				case 502:
				case 400:
				case 500:
					error.message = "响应中...";
					break;
				case 503:
					error.message = "服务不可用";
					break;
				case 504:
					error.message = "访问超时，请重试";
					break;
				case 505:
					error.message = "HTTP版本不受支持";
					break;
				default:
					break;
			}

			Vue.prototype.toasted.info(error.message);

			return Promise.reject(error);
		}
	);

	return service;
}

/**
 * @description 创建请求方法
 * @param {Object} service axios 实例
 */
function createRequest(service) {
	return function (config) {
		const token = localStorage.get("token");
		const user = localStorage.get("user");
		const project = storage.get("project") || {};
		const secret = localStorage.get("secret") || "";

		const headers = {
			"Content-Type": "application/json;charset=UTF-8",
			"x-is-intranet": String(lpTools.isIntranet())
		};

		if (secret) {
			const code = lpTools.GetCode(secret);
			headers["x-code"] = String(code);
		}

		if (user) {
			headers["x-user-id"] = user.id;
		}

		if (project.code) {
			headers["pro-code"] = project.code;
		}
		if(token){
			headers["x-token"] =  localStorage.get("token");
		}

		const defaultConfig = {
			method: "GET",
			baseURL: baseURLApi,
			headers
		};
	
		const assignConfig = Object.assign(defaultConfig, config);

		if (config.headers) {
			assignConfig.headers = {
				...headers,
				...config.headers
			};
		}

		return service(assignConfig);
	};
}

export const service = createService();
export const request = createRequest(service);
