<template>
	<div>
		<lp-dialog
			ref="dialog"
			cardTextClass="pa-0 ma-0"
			fullscreen
			persistent
			@dialog="handleDialog"
		>
			<!-- 顶部信息 BEGIN -->
			<v-toolbar slot="toolbar" color="primary" dark>
				<v-btn icon dark @click="closeDialog">
					<v-icon>mdi-close</v-icon>
				</v-btn>

				<!-- 工序 BEGIN -->
				<v-chip-group>
					<template v-for="(entry, index) in entries">
						<v-chip
							v-show="entry.show"
							:key="`entry_${index}`"
							:class="entry.class"
							:color="entry.actived ? '#fb8c00' : ''"
							:input-value="entry.actived"
							outlined
							@click="navPage(entry)"
						>
							{{ entry.label }}({{ task.topInfo[entry.value] | ifLousyValue }})
						</v-chip>
					</template>
				</v-chip-group>
				<!-- 工序 END -->

				<v-spacer></v-spacer>

				<!-- 页数、单号 BEGIN -->
				<v-toolbar-items v-if="op !== -1">
					<ul class="z-flex align-center">
						<li v-if="op === 'op0' && bill.phoTotal" class="mr-4">
							<label>页数：</label>
							<span>{{ bill.thumbIndex + 1 }}/{{ bill.phoTotal }}</span>
						</li>

						<li class="z-flex align-center">
							<v-icon>mdi-file-document-outline</v-icon>
							<span @dblclick="dblBillNum">{{ bill.billNum }}</span>
							<span v-if="bill.agency">、{{ bill.agency }}</span>
						</li>
					</ul>
				</v-toolbar-items>
				<!-- 页数、单号 END -->

				<v-spacer></v-spacer>

				<v-toolbar-items>
					<!-- 字符数、准确率、倒计时 BEGIN -->
					<ul v-if="op !== -1" class="z-flex align-center">
						<li class="mr-4">
							<label>字符数：</label>
							<span>{{ task.topInfo.character | ifLousyValue }}</span>
						</li>

						<li class="mr-4">
							<label>准确率：</label>
							<span>{{ (task.topInfo.accuracy * 100) | decimalFormat }}%</span>
						</li>

						<li class="z-flex align-center mr-4">
							<v-icon>mdi-clock-time-three-outline</v-icon>
							<span>{{ freeTime }}</span>
						</li>
					</ul>
					<!-- 字符数、准确率、倒计时 END -->

					<!-- 用户信息 BEGIN -->
					<div class="z-flex align-center">
						<v-icon>mdi-crown-outline</v-icon>
						<span>{{ user.code }} - {{ user.nickName }}</span>
					</div>
					<!-- 用户信息 END -->
				</v-toolbar-items>
			</v-toolbar>
			<!-- 顶部信息 END -->

			<template v-if="tools.isYummy(task.topInfo) && tools.isYummy(task.config)">
				<!-- 当前工序(op) BEGIN -->
				<div
					slot="main"
					class="main"
					:style="{
						height: task.mainHeight > 0 ? `${task.mainHeight}px` : null
					}"
				>
					<router-view
						ref="opRouter"
						@gotTaskResponse="handleGotTaskResponse"
						@submittedTaskResponse="handleSubmittedTaskResponse"
					></router-view>
				</div>
				<!-- 当前工序(op) END -->

				<!-- 当前 field 录入提示 BEGIN -->
				<div class="prompt" ref="prompt" slot="bottom">
					{{ task.prompt }}
				</div>
				<!-- 当前 field 录入提示 END -->
			</template>
		</lp-dialog>

		<!-- 删单 BEGIN -->
		<z-dynamic-form
			ref="delDynamic"
			:fieldList="cells.deleteFields"
			:width="550"
			@confirm="handleDelConfirm"
		></z-dynamic-form>
		<!-- 删单 END -->
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import { tools } from "vue-rocket";
import TaskDialogMixins from "./TaskDialogMixins";
import { newWaterMark, delWatermark } from "@/plugins/watermark";
import cells from "./cells";

const TOP_TIME = 15;
const OP_ITEMS = ["op0", "op1", "op2", "opq"];

export default {
	name: "EntryModal",
	mixins: [TaskDialogMixins],

	data() {
		return {
			op: "-1",

			cells,
			entries: cells.entries,

			user: {},
			bill: {},
			block: {},

			// 释放任务倒计时
			freeTime: 0,
			timer: null,

			// 间隔30秒获取一次顶部信息
			topTime: TOP_TIME,
			topTimer: null,

			// 问题件 查看图片
			thumbs: []
		};
	},

	methods: {
		// 当前工序领取任务后返回结果
		handleGotTaskResponse({ code, bill, block }) {
			this.bill = tools.isYummy(bill) ? bill : {};
			this.block = tools.isYummy(block) ? block : {};

			this.resetCountdown();
			this.resetTopCountdown();

			if (code === 200) {
				this.getTaskTopInfo();

				this.countdown();
				this.topCountdown();
			}
		},

		// 当前工序提交任务后返回结果
		handleSubmittedTaskResponse({ code }) {
			if (code === 200) {
				this.getTaskTopInfo();
			}
		},

		// 控制当前项目 op0 op1 op2 opq 显示/隐藏
		setEntries(meta) {
			const project = this.storage.get("project");
			const result = tools.find(this.auth.perm, { proCode: project.code });

			if (result) {
				this.entries.map(entry => {
					// 权限
					entry["show"] = tools.isYummy(result[entry.authKey]) ? true : false;

					// 选中当前op
					entry["actived"] = tools.isYummy(entry.value === meta.path) ? true : false;
				});
			}

			this.entries = [...this.entries];
		},

		// 弹窗状态
		handleDialog(dialog) {
			if (dialog) {
				this.getTaskTopInfo();

				const { code, nickName } = this.user;

				newWaterMark(`${code}-${nickName}`);
			} else {
				this.resetCountdown();
				this.resetTopCountdown();

				delWatermark();
			}
		},

		// 获取当前工序
		getOp(key) {
			this.op = OP_ITEMS[OP_ITEMS.indexOf(key)] || -1;
		},

		// 切换页面
		navPage(item) {
			this.$router.push({ path: item.link });
		},

		// 关闭弹框
		closeDialog() {
			this.$router.push({ path: "/main/entry/channel" });
			this.$refs.dialog.onClose();
		},

		// 倒计时
		countdown() {
			this.freeTime = this.block.freeTime;

			if (!this.freeTime) return;

			this.timer = setInterval(() => {
				--this.freeTime;

				if (this.freeTime < 0) {
					// 2022/5/12，泽如说倒计时为0要将用户踢出当前工序
					this.releaseTask();
					this.redirectTask();

					this.entries.map(entry => (entry["actived"] = false));

					this.resetCountdown();
				}
			}, 1000);
		},

		// 清空倒计时
		resetCountdown() {
			clearInterval(this.timer);
			this.freeTime = 0;
		},

		// 顶部信息倒计时
		topCountdown() {
			this.topTimer = setInterval(() => {
				--this.topTime;

				if (this.topTime <= 0) {
					this.resetTopCountdown();
					this.getTaskTopInfo();
					this.topCountdown();
				}
			}, 1000);
		},

		// 清空顶部信息倒计时
		resetTopCountdown() {
			clearInterval(this.topTimer);
			this.topTime = TOP_TIME;
		},

		// 设置 main 的高度
		setMainHeight() {
			const { height } = tools.getViewportSize();
			const promptHeight = this.$refs.prompt?.offsetHeight;

			const mainHeight = height - 64 - promptHeight;

			this.$store.commit("UPDATE_CHANNEL", { mainHeight });
		},

		dblBillNum() {
			if (this.$route.name === "Opq") {
				this.$refs.opRouter.navToViewImages();
			}
		}
	},

	computed: {
		...mapGetters(["auth", "task"])
	},

	watch: {
		$route: {
			handler(route) {
				const { meta, query } = route;

				this.getOp(meta.key);

				this.$nextTick(() => {
					this.setEntries(meta);

					if (meta.path || query.op) {
						this.$refs.dialog.onOpen();
					} else {
						this.$refs.dialog.onClose();
					}
				});
			},
			immediate: true
		},

		// 根据录入提示动态设置 main 的高度
		"task.prompt": {
			handler() {
				this.$nextTick(this.setMainHeight);
			},
			immediate: true
		}
	}
};
</script>

<style lang="scss">
.main {
	height: calc(100vh - 121px);
	padding-bottom: 8px;
	overflow: hidden;
}

.prompt {
	width: 100%;
	padding: 16px;
	border-top: 1px solid rgba(0, 0, 0, 0.12);
	background-color: #fff;
}
</style>
