<template>
	<div id="lpLogin" class="lp-login" @keyup.enter="onLogin">
		<div class="py-8 rounded-lg elevation-6 animate__animated animate__shakeX login-form">
			<div class="z-flex justify-center align-center top">
				<div class="logo">
					<img src="@/assets/logo.png" />
				</div>
				<h5 class="ml-4 text-h5">珠海汇流理赔数据处理平台2.0</h5>
			</div>

			<ul class="pl-0 mt-8 mx-auto form">
				<li class="mb-1">
					<z-text-field
						:formId="formId"
						:formKey="userField.formKey"
						:label="userField.label"
						:validation="userField.validation"
						:defaultValue="form[userField.formKey]"
					>
					</z-text-field>
				</li>

				<li class="mb-1">
					<z-text-field
						:formId="formId"
						formKey="password"
						:append-icon="eye ? 'mdi-eye' : 'mdi-eye-off'"
						label="密码"
						:type="eye ? 'text' : 'password'"
						:validation="[{ rule: 'required', message: '请输入密码!' }]"
						@click:append="eye = !eye"
						:defaultValue="form.password"
					>
					</z-text-field>
				</li>

				<!-- <li class="z-flex align-end mb-4" v-if="!isIntranet">
					<z-text-field
						:formId="formId"
						formKey="captcha"
						class="flex-grow-1"
						label="验证码"
						:validation="[{ rule: 'required', message: '请输入验证码!' }]"
					>
					</z-text-field>
					<z-btn
						class="pb-5 ml-4"
						:color="color"
						:lockedTime="2500"
						:disabled="!validAccount || counting"
						@click="sendCode"
						>{{ text }}</z-btn
					>
				</li> -->

				<li class="mb-0">
					<z-btn
						:formId="formId"
						btnType="validate"
						block
						:color="color"
						:disabled="pending"
						:loading="pending"
						@click="onLogin"
						>登录</z-btn
					>
				</li>

				<li class="z-flex justify-between px-1 mt-4 bottom">
					<span class="z-flex align-center">
						<v-checkbox v-model="memoPassword" class="ma-0 pa-0" hide-details @change="onMemo"></v-checkbox>
						记住密码
					</span>

					<span v-if="!isIntranet" class="mr-12" @click="onForgot"> 忘记密码 </span>

					<span @click="onChange"> 修改密码 </span>
				</li>
			</ul>

			<div class="z-flex justify-between contact">
				<span>联系方式：0756-3630987</span>
				<span>在线客服：关注公众号"珠海汇流HR"</span>
			</div>
		</div>

		<forgot-password ref="forgot"></forgot-password>

		<change-password ref="change"></change-password>

		<div class="extra-bar">
			<ul class="z-flex align-center">
				<!-- <li v-if="!isIntranet" @click="authenticationFun">身份验证</li> -->
				<span v-if="!isIntranet" class="mx-2"></span>
				<li v-if="!isIntranet" @click="navOfficial">公司官网</li>
				<span v-if="!isIntranet" class="mx-2 line"></span>
				<li @click="navUsage">使用说明</li>
			</ul>
		</div>

		<bubble-background></bubble-background>
		<authentication ref="authDom" :networkIdetify="isIntranet" />
	</div>
</template>

<script>
import { localStorage, tools } from "vue-rocket";
import ButtonMixins from "@/mixins/ButtonMixins";
import LoginMixins from "./LoginMixins";
import AuthMixins from "./AuthMixins";
import { sessionStorage } from "vue-rocket";

export default {
	name: "Login",
	mixins: [ButtonMixins, LoginMixins, AuthMixins],

	data() {
		return {
			formId: "Login",
			eye: false,
			memoPassword: false,
			form: {}
		};
	},

	created() {
		this.getMemo();
		localStorage.delete("secret");
	},
	mounted() {
		if (!this.memoPassword) {
			this.form = {};
		}
	},
	methods: {
		// 登录
		async onLogin() {
			const accountKey = this.isIntranet ? "username" : "phone";

			this.isPending(true);

			const form = {
				isIntranet: this.isIntranet,
				accountKey,
				captchaId: this.captchaId,
				...this.forms[this.formId]
			};

			const result = await this.$store.dispatch("LOGIN", form);

			this.isPending(false);
			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				localStorage.set({
					token: result.data.token,
					user: result.data.user,
					secret: result.data.secret,
					account: { account: form.account, password: form.password },
					viewport: tools.getViewportSize()
				});
				console.log(tools.getViewportSize());
				this.getRoleSysMenu();
			}
		},
		async getRoleSysMenu() {
			const result = await this.$store.dispatch("GET_ROLE_SYS_MENU");
			console.log("getRoleSysMenu", result);
			if (result.code === 200) {
				const { Menus, Perm, proCode = "" } = result.data;
				const [mapPro, proItems] = this.resolvePerm(Perm, proCode);

				sessionStorage.set("proCode", proCode);
				const { isApp } = this.$route.query;
				let isApps = isApp || false;
				sessionStorage.set("isApp", { isApp: isApps });

				if (isApp && isApp === "true") {
					this.$router.push({ path: "/transit" });
				} else {
					this.$router.replace({ name: "Channel", query: { sync: 1 } });
				}
				this.$store.commit("UPDATE_AUTH", {
					menus: Menus,
					perm: Perm,
					mapPro,
					proItems,
					isApp: isApp
				});
			}
		},
		resolvePerm(permissions, proCode) {
			const mapPro = {};
			const proItems = [];

			permissions.map(permission => {
				if (permission.proCode === proCode && this.hasAuth(permission)) {
					mapPro[permission.proCode] = permission;

					proItems.push({
						label: permission.proCode,
						value: permission.proCode
					});
				}
			});

			return [mapPro, proItems];
		},
		hasAuth({ hasOp0, hasOp1, hasOp2, hasOpq }) {
			if (hasOp0 || hasOp1 || hasOp2 || hasOpq) {
				return true;
			}

			return false;
		},
		onMemo() {
			localStorage.set("memoPassword", this.memoPassword);
			if (!this.memoPassword) {
				localStorage.delete("account");
			}
		},

		getMemo() {
			const memoPassword = localStorage.get("memoPassword");
			const account = localStorage.get("account");

			this.memoPassword = memoPassword;
			this.form = account || {};
		},

		// 打开忘记密码弹框
		onForgot() {
			this.$refs.forgot.onOpen();
		},

		// 打开修改密码弹框
		onChange() {
			this.$refs.change.onOpen();
		},
		// 认证
		authenticationFun() {
			this.$refs.authDom.onOpen();
		},
		navOfficial() {
			window.open("http://www.i-confluence.com");
		},

		navUsage() {
			this.$router.push({ path: "usage" });
		}
	},

	components: {
		"forgot-password": () => import("./forgotPassword"),
		"change-password": () => import("./changePassword"),
		"bubble-background": () => import("./bubbleBackground"),
		authentication: () => import("./authentication")
	}
};
</script>

<style lang="scss" scoped>
.lp-login {
	position: relative;
	width: 100%;
	height: 100%;
	background-color: #4e54c8;
	overflow: hidden;

	.login-form {
		position: absolute;
		top: 50%;
		left: 50%;
		margin-top: -225px;
		margin-left: -250px;
		width: 500px;
		background-color: #fff;
		z-index: 1;

		.logo,
		.logo img {
			width: 65px;
			height: 65px;
		}

		ul.form {
			width: 400px;

			.bottom {
				cursor: pointer;
			}
		}
	}

	.contact {
		position: absolute;
		width: inherit;
		left: 0;
		bottom: -80px;
		color: #fff;
		font-size: 14px;
	}

	.extra-bar {
		position: absolute;
		top: 24px;
		right: 24px;
		color: #fafafa;
		cursor: pointer;
		z-index: 1;

		& li:hover {
			color: #fafafa;
			text-decoration: underline;
		}

		& span.line {
			display: block;
			width: 2px;
			height: 14px;
			background-color: #fafafa;
		}
	}
}
</style>
