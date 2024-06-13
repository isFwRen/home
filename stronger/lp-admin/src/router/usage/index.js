const UsageRoutes = {
	path: "/usage",
	name: "Usage",
	meta: {
		key: "usage",
		title: "使用说明"
	},
	component: () => import("@/views/usage"),

	children: [
		{
			path: "/usage",
			redirect: "registrationGuide"
		},

		{
			path: "registrationGuide",
			name: "RegistrationGuide",
			meta: {
				pKey: "usage",
				key: "registrationGuide",
				realm: "registrationGuide",
				title: "如何注册"
			},
			component: () => import("@/views/usage/registrationGuide")
		},

		{
			path: "loginGuide",
			name: "LoginGuide",
			meta: {
				pKey: "usage",
				key: "loginGuide",
				realm: "loginGuide",
				title: "如何登录"
			},
			component: () => import("@/views/usage/loginGuide")
		},

		{
			path: "forgotPassword",
			name: "ForgotPassword",
			meta: {
				pKey: "usage",
				key: "forgotPassword",
				realm: "forgotPassword",
				title: "忘记密码"
			},
			component: () => import("@/views/usage/forgotPassword")
		},

		{
			path: "code",
			name: "Code",
			meta: {
				pKey: "usage",
				key: "code",
				realm: "code",
				title: "无法获取验证码"
			},
			component: () => import("@/views/usage/code")
		},

		{
			path: "forgotJobNumber",
			name: "ForgotJobNumber",
			meta: {
				pKey: "usage",
				key: "forgotJobNumber",
				realm: "forgotJobNumber",
				title: "忘记工号"
			},
			component: () => import("@/views/usage/forgotJobNumber")
		},

		{
			path: "restoreJobNumber",
			name: "RestoreJobNumber",
			meta: {
				pKey: "usage",
				key: "restoreJobNumber",
				realm: "restoreJobNumber",
				title: "恢复工号"
			},
			component: () => import("@/views/usage/restoreJobNumber")
		}
	]
};

export default UsageRoutes;
