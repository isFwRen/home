<template>
	<div class="lp-notification animate__animated" :class="{ animate__swing: total > 0 }">
		<v-badge
			color="#ed561b"
			:content="messageCount"
			offset-x="14"
			offset-y="14"
			overlap
			:value="messageCount"
		>
			<z-btn icon small @click="handleMessage">
				<v-icon color="#fff"> mdi-bell </v-icon>
			</z-btn>
		</v-badge>

		<div
			ref="globalNoticeRef"
			v-show="visible && globalNotice.count > 0"
			class="message-box animate__animated"
			:class="{
				animate__backInRight: visible
			}"
		>
			<v-card light width="360">
				<v-card-text class="pb-0">
					<div class="z-flex solid-bottom">
						<div class="mr-2">
							<v-avatar color="#4b90e1" size="48">
								<span class="white--text text-14">{{ globalNotice.proCode }}</span>
							</v-avatar>
						</div>

						<div class="flex-grow-1">
							<div class="z-flex justify-between">
								<h6 class="text-primary">{{ globalNotice.fileName }}</h6>
								<span
									v-show="globalNotice.sendTime"
									class="text-12 text-secondary"
									>{{
										globalNotice.sendTime | dateFormat("YYYY-MM-DD HH:mm:ss")
									}}</span
								>
							</div>
							<p class="text-bold">{{ globalNotice.content }}</p>
						</div>
					</div>
				</v-card-text>

				<v-card-actions class="z-flex justify-end">
					<z-btn class="mr-4" text @click="handleLater">
						<v-icon left> mdi-update </v-icon>
						稍后处理
					</z-btn>

					<z-btn text @click="handleReply">
						<v-icon left> mdi-message-reply-text-outline </v-icon>
						回复
					</z-btn>
				</v-card-actions>
			</v-card>
		</div>

		<div
			ref="businessNoticeRef"
			class="speed-card animate__animated"
			:style="computeStyle"
			:class="{
				animate__backInRight: visible
			}"
			v-show="show && businessNotice.count > 0"
		>
			<div class="speed-left">
				<div class="speed-avatar">{{ businessNotice.proCode }}</div>
			</div>
			<div class="speed-right">
				<div class="speed-right-notify">
					<div class="notify">{{ businessNotice.title }}</div>
					<div class="time">
						{{ businessNotice.CreatedAt | formTime }}
					</div>
				</div>
				<div class="speed-subtitle mt-1 text-bold">
					{{ businessNotice.msg }}
				</div>
				<div class="speed-action mt-1">
					<z-btn class="mr-4" text @click="onLater">
						<v-icon left> mdi-download-box </v-icon>
						稍后处理
					</z-btn>

					<z-btn text @click="onBusiness">
						<v-icon left> mdi-card-search </v-icon>
						快速查看
					</z-btn>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
import { localStorage } from "vue-rocket";
import { mapGetters } from "vuex";
import io from "socket.io-client";
import { tools as lpTools } from "@/libs/util";
import moment from "moment";
const { baseURL } = lpTools.baseURL();
const user = localStorage.get("user");
const userId = user?.id ?? "";
// const socket = io.connect(`${baseURL}global-notice?userId=${userId}`);
const Messsocket = io.connect(`${baseURL}global-message?userId=${userId}`);
export default {
	name: "LPNotification",

	data() {
		return {
			total: 1,
			visible: false,
			show: false
		};
	},
	filters: {
		formTime(time) {
			if (!time) return "";
			return moment(time).format("YYYY-MM-DD HH:mm:ss");
		}
	},
	computed: {
		...mapGetters(["globalNotice", "businessNotice"]),
		messageCount() {
			return this.globalNotice.count + this.businessNotice.count;
		},
		computeStyle() {
			const obj = {};
			if (this.visible && this.globalNotice.count > 0) {
				obj.top = "170px";
			} else {
				obj.top = "36px";
			}
			return obj;
		}
	},

	created() {
		// 时效简报
		Messsocket.on("briefing", result => {
			this.$store.commit("NOTIFICATION_TIMEBRIEFS_UPDATE", JSON.parse(result.ageingBriefing));

		});
		// socket.on("customerNotice", result => {
		// 	this.$store.commit("GLOBAL_NOTIFICATION_UPDATE_ITEM", result);
		// });
		// // 新能消息通知功能8.21
		// Messsocket.on("download", result => {
		// 	this.addMessage(result)
		// 	console.log(result, "download");
		// });
		// Messsocket.on("upload", result => {
		// 	console.log(result, "upload");
		// 	this.addMessage(result)
		// });
	},

	methods: {
		addMessage(result) {
			this.$store.commit("INCREAT_COUNT");
			const { proCode, title, CreatedAt, msg, ID } = result.businessPush;
			const pushId = result.businessPushSend.pushId;
			this.$store.commit("BUSINESS_NOTIFICATION_UPDATE_ITEM", {
				proCode,
				title,
				CreatedAt,
				msg,
				ID,
				pushId
			});
			this.show = true;
		},
		handleRouteLink(path) {
			if (this.$route.fullPath !== path) {
				this.$router.push({
					path
				});
			}
		},
		handleMessage() {
			this.visible = !this.visible;
			this.show = !this.show;
		},
		onBusiness() {
			this.show = false;
			this.handleRouteLink("/main/notices/business");
			this.$store.commit("BUSINESS_NOTIFICATION_UPDATE_ITEM", {
				count: 0
			});
		},
		onLater() {
			this.show = false;
		},
		handleLater() {
			this.visible = false;
		},
		handleReply() {
			this.visible = false;
			this.handleRouteLink("/main/notices/customer");
			//this.$router.push({ path: "/main/notices/customer" });
		}
	}
};
</script>

<style scoped lang="scss">
.lp-notification {
	position: relative;

	.message-box {
		position: absolute;
		top: 36px;
		right: 0;
	}
}

.speed-card {
	width: 360px;
	background-color: #fff;
	padding: 15px 10px 5px 10px;
	display: flex;
	font-size: 14px;
	border-radius: 5px;
	box-shadow: 1px 1px 2px 2px #eee;
	position: absolute;
	top: 0;
	right: 0;
	z-index: 100;
	transition-delay: 0.5s;
	.speed-avatar {
		min-width: 50px;
		min-height: 50px;
		background-color: #4b90e1;
		display: flex;
		justify-content: center;
		align-items: center;
		border-radius: 50%;
	}
	.speed-right {
		flex: 1;
		color: #00000099;
		padding-left: 5px;
		.speed-right-notify {
			display: flex;
			justify-content: space-between;
			color: #000;
			.time {
				font-size: 12px;
				color: gray;
			}
		}
		.speed-subtitle {
			border-bottom: 2px solid #eee;
			padding-bottom: 10px;
			width: 280px;
			white-space: nowrap;
			text-overflow: ellipsis;
			overflow: hidden;
		}
		.speed-action {
			display: flex;
			width: 70%;
			justify-content: space-between;
		}
	}
}
.focus-list {
	border: 2px solid #3b82f6;
}
</style>
