<template>
	<v-app-bar app color="#1976d2" dark clipped-left elevate-on-scroll style="z-index: 7">
		<v-app-bar-nav-icon @click="onDrawer"></v-app-bar-nav-icon>

		<v-avatar class="mr-4" size="52">
			<img src="@/assets/logo.png" alt="汇流" />
		</v-avatar>

		<template v-slot:img="{ props }">
			<v-img v-bind="props" gradient="to top right, rgba(19,84,122,.5), rgba(128,208,199,.8)"></v-img>
		</template>

		<v-toolbar-title>珠海汇流理赔数据处理平台2.0 - 录入系统</v-toolbar-title>

		<v-spacer></v-spacer>

		<template v-if="moduleName !== 'usage'">
			<div class="z-flex align-center mr-3 user-info" style="position: relative">
				<Personal ref="PersonalRef" :user="user" :avatarUrl="avatarUrl" @emitUploadImg="onEmitUploadImg" />
				<v-avatar v-if="user.headerImg" color="indigo" size="42" @click="showPop" style="cursor: pointer">
					<img :src="avatarUrl" />
				</v-avatar>

				<v-icon v-else size="26">mdi-account-circle</v-icon>

				<span class="pl-1">{{ user.name }}</span>

				<v-icon class="pl-4">mdi-card-account-details</v-icon>

				<span class="pl-1">{{ user.code }}</span>
			</div>

			<z-btn icon @click="onSignOut">
				<v-icon>mdi-export</v-icon>
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
		const user = localStorage.get("user");
		if (user) {
			this.getAvatar();
			this.user = user;
		} else {
			const timer = setTimeout(() => {
				const user = localStorage.get("user");
				if (user) {
					this.user = user;
				}
				clearTimeout(timer);
			}, 200);
		}
	},
	methods: {
		onEmitUploadImg({ base64Url }) {
			this.avatarUrl = base64Url;
		},
		onDrawer() {
			this.$emit("drawer");
		},
		async getAvatar() {
			const user = localStorage.get("user");
			const url = `${baseURLApi}${user.headerImg}`;
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
		onSignOut() {
			let isApp = sessionStorage.get("isApp");
			this.$modal({
				visible: true,
				title: "退出提示",
				content: `请确认是否要退出？`,
				confirm: async () => {
					let url = localStorage.get("task")?.baseURL;
					if (url) {
						const result = await this.$store.dispatch("ALLOCATION_ALL_TASK", { code: this.user.code, op: "op0" });
						if (result.code === 200) {
							this.toasted.warning("退出释放所有分块!", 1000);
						}
					}
					localStorage.delete(["secret", "token", "caseInfo", "auth", "project", "user", "lp2ConstantBaseInfo"]);
					location.replace(`${location.origin}/login?isApp=${isApp.isApp}`);
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
