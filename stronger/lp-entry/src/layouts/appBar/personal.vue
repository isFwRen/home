<template>
	<transition name="fade">
		<div>
			<div class="pop_avatar_box" v-if="pop" v-click-outside="popBlur">
				<h4 class="pop_avatar_title pt-2 pr-4" align="center">账户信息</h4>
				<div class="pop_input_content pb-3 pt-3 pr-4">
					<div class="img_round" @mouseenter="mouseEnter" @mouseleave="mouseLeave">
						<label for="avatarInput">
							<img :src="newImgUrl" />
						</label>
						<input
							type="file"
							id="avatarInput"
							accept=".jpg, .jpeg, .png, .tif"
							name="avatar"
							@change="onFileChange"
						/>
						<transition name="fade">
							<div class="img_icon" v-if="showmask">
								<div class="icon_camera">
									<v-icon color="#fff" size="50"> mdi-camera </v-icon>
								</div>
							</div>
						</transition>
					</div>
				</div>
				<div class="userinfo">
					<div class="jobNum z-flex justify-center">
						<div class="info-left">工号</div>
						<div class="info-right">{{ user.code }}</div>
					</div>
					<div class="nickName z-flex justify-center">
						<div class="info-left">姓名</div>
						<div class="info-right">{{ user.name }}</div>
					</div>
					<div class="phone z-flex justify-center">
						<div class="info-left">手机号</div>
						<div class="info-right">{{ user.phone | filterPhone }}</div>
					</div>
				</div>
				<div class="z-flex justify-center mb-4 mt-4" style="width: 100%">
					<v-btn class="mr-2" color="error" @click="closeAccount" small> 注销账号 </v-btn>
					<v-btn color="primary" class="ml-2" small @click="handleUpload"> 上传头像 </v-btn>
				</div>
			</div>
			<v-dialog v-model="dialog" max-width="320">
				<v-card>
					<v-card-title class="text-h5"> 注销账号 </v-card-title>
					<v-card-text>
						<p :class="[isLogout ? 'text-center' : 'text-left']">
							<v-icon color="error" v-if="!isLogout">mdi-alert-circle</v-icon>
							{{ content }}
						</p>
						<v-spacer></v-spacer>
						<p v-if="isLogout" :class="[isLogout ? 'text-center' : 'text-left']">
							<span class="font-weight-bold">{{ number }}</span
							>s
						</p>
					</v-card-text>
					<v-card-actions v-show="!isLogout">
						<v-spacer></v-spacer>
						<v-btn color="primary" @click="dialog = false"> 再想想 </v-btn>
						<v-btn color="error" @click="deletAccount"> 去意已决 </v-btn>
					</v-card-actions>
				</v-card>
			</v-dialog>
		</div>
	</transition>
</template>
<script>
import { localStorage, sessionStorage } from "vue-rocket";
export default {
	data() {
		return {
			showmask: false,
			pop: false,
			dialog: false,
			newImgUrl: "",
			isLogout: false,
			number: 5,
			base64Url: "",
			formData: {}
		};
	},
	computed: {
		content() {
			return this.isLogout ? "账号已注销,江湖再见" : "注销账号会清空所有信息和数据，您是否确认注销？";
		}
	},
	props: {
		user: {
			type: Object
		},
		avatarUrl: {
			type: String
		}
	},
	watch: {
		avatarUrl(base64) {
			this.newImgUrl = base64;
		}
	},
	filters: {
		filterPhone: function (value) {
			if (!value) return "";
			return value.replace(value.slice(3, 7), "****");
		}
	},
	methods: {
		async handleUpload() {
			if (!this.base64Url) {
				this.toasted.warning("请选择要上传的图片");
				return;
			}
			const result = await this.$store.dispatch("UPLOAD_AVATAR", this.formData);
			const token = localStorage.get("token");
			// 更新user
			const body = {
				token,
				userId: this.user.id
			};
			this.getUserInfo(body);

			this.$emit("emitUploadImg", {
				base64Url: this.base64Url
			});
		},
		async getUserInfo(data) {
			const form = {
				"x-token": data.token,
				"x-user-id": data.userId
			};
			const result = await this.$store.dispatch("GET_USER_INFO", form);
			if (result.code === 200) {
				localStorage.set("user", result.data.user);
			}
		},
		setTimeoutFun() {
			const timer = setInterval(() => {
				this.number = this.number - 1;
				if (this.number === 0) {
					clearInterval(timer);
					localStorage.delete(["secret", "token", "caseInfo", "auth", "project", "user", "lp2ConstantBaseInfo"]);
					sessionStorage.delete(["CaseSearch", "thumbs"]);
					location.replace(`${location.origin}/login`);
				}
			}, 1000);
		},
		async deletAccount() {
			this.isLogout = true;
			const data = {
				name: this.user.name,
				reason: "",
				username: this.user.code
			};

			// 注销成功调用计时器
			const result = await this.$store.dispatch("DELETE_ACCOUNT", data);
			if (result.code === 200) {
				this.setTimeoutFun();
			}
			this.setTimeoutFun();
		},
		closeAccount() {
			this.dialog = true;
		},
		mouseEnter() {
			this.showmask = true;
		},
		mouseLeave() {
			this.showmask = false;
		},
		showPop() {
			this.pop = true;
		},
		popBlur() {
			this.pop = false;
		},
		getBase64(img, callback) {
			const reader = new FileReader();
			reader.addEventListener("load", () => callback(reader.result));
			reader.readAsDataURL(img);
		},
		async onFileChange(event) {
			const file = event.target.files[0];
			if (!file) {
				this.newImgUrl = "";
				return;
			}
			const isLt2M = file.size / 1024 / 1024 < 2;
			if (!isLt2M) {
				this.toasted.warning("请上传2M以内图片");
				return;
			}

			// 上传头像
			this.formData = new FormData();
			this.formData.append("file", file);
			this.formData.append("jobNumber", this.user.code);
			this.formData.append("jobName", this.user.name);

			this.getBase64(file, base64Url => {
				this.newImgUrl = base64Url;
				this.base64Url = base64Url;
			});
		}
	}
};
</script>
<style scoped lang="scss">
.fade-enter-active,
.fade-leave-active {
	opacity: 1;
	transition: opacity 0.4s;
}

.fade-enter,
.fade-leave-to {
	opacity: 0;
}

.pop_avatar_box {
	border: 1px solid #0f172a1a;
	z-index: 999;
	width: 220px;
	position: absolute;
	bottom: 0;
	margin-left: -65%;
	transform: translateY(calc(100% + 20px));
	border-radius: 10px;
	color: #334155;
	font-weight: 600;
	font-size: 14px;
	background-color: #ffffff;
	box-shadow: 0 10px 15px -3px #0000001a, 0 4px 6px -4px #0000001a;

	.pop_avatar_title {
		color: #1976d2;
	}

	.pop_input_content {
		width: 100%;

		.img_round {
			width: 85px;
			height: 85px;
			margin: 0 auto;
			border: 1px solid #eee;
			border-radius: 50%;
			overflow: hidden;
			position: relative;

			& input[type="file"] {
				position: absolute;
				width: 100%;
				height: 100%;
				left: 0;
				top: 0;
				opacity: 0;
				cursor: pointer;
				z-index: 10;
			}

			& img {
				width: 100%;
				height: 100%;
				object-fit: cover;
				cursor: pointer;
			}

			.img_icon {
				position: absolute;
				left: 0;
				top: 0;
				width: 100%;
				height: 100%;
				color: #fff;
				background-color: rgba(0, 0, 0, 0.15);
				cursor: pointer;

				.icon_camera {
					position: absolute;
					left: 0;
					right: 0;
					top: 0;
					bottom: 0;
					margin: auto;
					width: 50px;
					height: 50px;
				}
			}
		}
	}

	.userinfo {
		letter-spacing: 1px;

		.info-left {
			flex: 1;
			text-align: center;
		}

		.info-right {
			flex: 1;
			color: #409eff;
		}
	}
}
</style>
