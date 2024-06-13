<template>
	<div id="lpLogin" class="lp-login" @keyup.enter="onLogin">
		<div class="py-8 rounded-lg elevation-6 animate__animated animate__jello login-form">
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
						:validation="[
							{
								rule: 'required',
								message: '请输入验证码!'
							}
						]"
					>
					</z-text-field>
					<z-btn
            class="pb-5 ml-4"
            :color="color"
            :lockedTime="2500"
            :disabled="!validAccount || counting"
            @click="sendCode"
          >{{ text }}</z-btn>
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
						<v-checkbox
							v-model="memoPassword"
							class="ma-0 pa-0"
							hide-details
							@change="onMemo"
						></v-checkbox>
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

		<sliding-background></sliding-background>
		<authentication ref="authDom" :networkIdetify="isIntranet" />
	</div>
</template>

<script>
import { localStorage } from "vue-rocket";
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
	},
	mounted() {
		if (!this.memoPassword) {
			this.form = {};
		}
	},
	beforeRouteEnter(to, from, next) {
		next(vm => {
			sessionStorage.set("CaseSearch", {});
			let { isApp } = vm.$route.query;
			if (!isApp) {
				isApp = false;
			}
			if (isApp && isApp === "false") {
				isApp = false;
			}
			if (isApp && isApp === "true") {
				isApp = true;
			}

			sessionStorage.set("client", { isApp: Boolean(isApp) });
		});
	},
	methods: {
		// 认证
		authenticationFun() {
			this.$refs.authDom.onOpen();
		},
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
				localStorage.set("token", result.data.token);
				localStorage.set("secret", result.data.secret);
				localStorage.set("user", result.data.user);
				localStorage.set("account", {
					account: form.account,
					password: form.password
				});

				this.getRoleSysMenu();
				// this.$router.replace({ name: 'Case' })
				this.$router.replace({
					name: "Case",
					query: { sync: 1 }
				});
			}
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
		"sliding-background": () => import("./slidingBackground"),
		authentication: () => import("./authentication")
	}
};
</script>

<style lang="scss" scoped>
.lp-login {
	position: relative;
	width: 100%;
	height: 100%;
	background: rgb(8, 36, 47);
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
