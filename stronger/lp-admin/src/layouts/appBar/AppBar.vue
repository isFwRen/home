<template>
	<v-app-bar app color="#1976d2" clipped-left elevate-on-scroll style="z-index: 7; color: #fff">
		<v-app-bar-nav-icon color="#fff" @click="onDrawer"></v-app-bar-nav-icon>

		<v-avatar class="mr-4" size="52">
			<img src="@/assets/logo.png" alt="汇流" />
		</v-avatar>

		<template v-slot:img="{ props }">
			<v-img
				v-bind="props"
				gradient="to top right, rgba(19,84,122,.5), rgba(128,208,199,.8)"
			></v-img>
		</template>

		<v-toolbar-title>珠海汇流理赔数据处理平台2.0 - 管理系统</v-toolbar-title>

		<v-spacer></v-spacer>

		<template v-if="moduleName !== 'usage'">
			<div class="z-flex align-center mr-3 user-info" style="position: relative">
				<Personal
					ref="PersonalRef"
					:user="user"
					:avatarUrl="avatarUrl"
					@emitUploadImg="onEmitUploadImg"
				/>
				<v-avatar
					v-if="user.headerImg"
					color="indigo"
					size="42"
					@click="showPop"
					style="cursor: pointer"
				>
					<img :src="avatarUrl" />
				</v-avatar>

				<v-icon v-else size="26">mdi-account-circle</v-icon>

				<span class="pl-1">{{ user.name }}</span>

				<v-icon class="pl-4" color="#fff">mdi-card-account-details</v-icon>

				<span class="pl-1">{{ user.code }}</span>
			</div>

			<lp-notification></lp-notification>

			<z-btn icon @click="onSignOut">
				<v-icon color="#fff">mdi-export</v-icon>
			</z-btn>
		</template>
	</v-app-bar>
</template>

<script>
import { localStorage, sessionStorage } from "vue-rocket";
import Personal from "./personal.vue";
import { tools as lpTools } from "@/libs/util";
const { baseURLApi } = lpTools.baseURL();

export default {
	name: "NavigateBar",

	data() {
		return {
			user: {},
			moduleName: "",
			avatarUrl: ""
		};
	},

	mounted() {
		const user = this.storage.get("user");
		if (user) {
			this.getAvatar();
			this.user = user;
		}
	},

	methods: {
		onEmitUploadImg({ base64Url }) {
			this.avatarUrl = base64Url;
		},
		async getAvatar() {
			const user = localStorage.get("user");
			if (!user.headerImg) {
				return;
			}
			const url = `https://www.i-confluence.com:31111/api/${user.headerImg}`;
			const newBase64 = await lpTools.getTokenImg(url);
			if (newBase64) {
				lpTools.getBase64(newBase64).then(base64String => {
					this.avatarUrl = base64String;
				});
			}
		},
		showPop() {
			this.$refs.PersonalRef.showPop();
		},
		onDrawer() {
			this.$emit("drawer");
		},
		onSignOut() {
			const client = JSON.parse(window.sessionStorage.getItem("client"));
			const isApp = client.isApp;
			this.$modal({
				visible: true,
				title: "退出提示",
				content: `请确认是否要退出？`,
				confirm: () => {
					localStorage.delete([
						"secret",
						"token",
						"caseInfo",
						"auth",
						"project",
						"user",
						"lp2ConstantBaseInfo"
					]);

					sessionStorage.delete(["CaseSearch", "thumbs", "client"]);
					location.replace(`${location.origin}/login?isApp=${isApp}`);
				}
			});
		}
	},
	components: {
		Personal
	},
	watch: {
		$route: {
			handler(route) {
				const { pKey } = route.meta;
				this.moduleName = pKey;
			},
			immediate: true
		}
	}
};
</script>
