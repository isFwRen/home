<template>
	<div>
		<lp-dialog
			ref="dialog"
			card-text-class="pa-0 ma-0"
			fullscreen
			hide-overlay
			persistent
			:retain-focus="false"
			@dialog="handleDialog"
		>
			<!-- 顶部信息 BEGIN -->
			<v-toolbar slot="toolbar" color="primary" dark height="49">
				<v-btn icon dark @click="closeDialog">
					<v-icon>mdi-close</v-icon>
				</v-btn>

				<!-- 工序 BEGIN -->
				<v-chip-group>
					<template v-for="(op, index) in opTabs">
						<v-chip
							v-show="op.show"
							:key="`op_${index}`"
							:class="index < opTabs.length - 1 ? 'mr-3' : ''"
							:color="op.actived ? '#fb8c00' : ''"
							:input-value="op.actived"
							outlined
							@click="navigatePage(op)"
						>
							{{ op.label }}({{ task.topInfo[op.value] | ifLousyValue }})
						</v-chip>
					</template>
				</v-chip-group>
				<!-- 工序 END -->

				<v-spacer></v-spacer>

				<!-- 页数、单号 BEGIN -->
				<v-toolbar-items v-if="op !== -1">
					<ul class="z-flex align-center" v-if="task.displayPrompt">
						<li v-if="op === 'op0' && bill.phoTotal" class="mr-4">
							<label>页数：</label>
							<span>{{ bill.thumbIndex + 1 }}/{{ bill.phoTotal }}</span>
						</li>

						<li class="z-flex align-center">
							<v-icon>mdi-file-document-outline</v-icon>
							<span @dblclick="onBillNum">{{ bill.billNum }}</span>
							<span v-if="bill.agency">、{{ bill.agency }}</span>
						</li>
					</ul>

					<ul class="z-flex align-center" v-if="task.displayPrompt">
						<li class="z-flex align-center" v-if="op === 'opp'">
							<span>本次录入练习倒计时：</span>
							<span>{{ practiceTimes }}</span>
						</li>
					</ul>

					<ul class="z-flex align-center" v-if="!task.displayPrompt">
						<li class="z-flex align-center" v-if="op === 'ope'">
							<span>考核倒计时：</span>
							<span>{{ examTimes }}</span>
						</li>
					</ul>
				</v-toolbar-items>
				<!-- 页数、单号 END -->

				<v-spacer></v-spacer>

				<v-toolbar-items>
					<!-- 字符数、准确率、倒计时 BEGIN -->
					<ul v-if="op !== -1 && task.displayPrompt && !task.displayRight" class="z-flex align-center">
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

					<ul v-if="op !== -1 && task.displayRight" class="z-flex align-center">
						<li class="mr-4">
							<label>字符数：</label>
							<span>{{ this.$store.state["entry/task"].task.character }}</span>
						</li>

						<li class="mr-4">
							<label>准确率：</label>
							<span>{{ this.$store.state["entry/task"].task.accuracyRate * 100 + "%" }}</span>
						</li>

						<li class="z-flex align-center mr-4" style="cursor: pointer">
							<v-icon>mdi-exit-to-app</v-icon>
							<span @click="modal = true">结束练习</span>
						</li>
					</ul>
					<!-- 字符数、准确率、倒计时 END -->

					<!-- 用户信息 BEGIN -->
					<div class="z-flex align-center">
						<v-icon>mdi-crown-outline</v-icon>
						<span>{{ user.code }} - {{ user.name }}</span>
					</div>
					<!-- 用户信息 END -->
				</v-toolbar-items>
			</v-toolbar>
			<!-- 顶部信息 END -->

			<template v-if="tools.isYummy(task.topInfo) && tools.isYummy(task.config)">
				<!-- 当前工序(op) BEGIN -->
				<div
					slot="main"
					:class="['main', task.mainHeight == 0 ? 'h' : null]"
					:style="{
						height: task.mainHeight > 0 ? `${task.mainHeight}px` : null
					}"
				>
					<router-view
						ref="opRouter"
						@gotTaskResponse="handleGotTaskResponse"
						@submittedTaskResponse="handleSubmittedTaskResponse"
						@examContent="handleExam"
						@practiceContent="handlePractice"
					></router-view>
				</div>
				<!-- 当前工序(op) END -->

				<!-- 当前 field 录入提示 BEGIN -->
				<div ref="prompt" slot="bottom" class="prompt" v-if="task.displayPrompt">
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

		<v-dialog v-model="modal" transition="dialog-top-transition" max-width="600" persistent>
			<template>
				<v-card>
					<v-toolbar color="primary" dark style="font-size: 22px">是否结束练习</v-toolbar>
					<div class="font">
						结束练习不会保留当前的练习记录，再次进入将重新开始。如需保留记录，请点击左侧X按钮退出当前练习
					</div>
					<v-card-actions class="justify-end">
						<v-btn color="primary" @click="confirm">确定</v-btn>
						<v-btn color="primary" @click="modal = false">关闭</v-btn>
					</v-card-actions>
				</v-card>
			</template>
		</v-dialog>

		<SocketDialog ref="sockets" :data="socketMessages" />
	</div>
</template>

<script>
import SocketDialog from "./socketDialog/index.vue";
import { tools as lpTools } from "@/libs/util";
import { mapGetters } from "vuex";
import { tools, sessionStorage, localStorage } from "vue-rocket";
import { newWaterMark, delWatermark } from "@/plugins/watermark";
import TaskDialogMixins from "./mixins/TaskDialogMixins";
import cells from "./cells";
import moment from "moment";

const isIntranet = lpTools.isIntranet();

const TOP_TIME = 15;
const OP_ITEMS = ["op0", "op1", "op2", "opq", "opp", "ope"];

const [toolbarHeight, defaultPromptHeight] = [64, 54];

export default {
	name: "EntryModal",
	mixins: [TaskDialogMixins],
	inject: ["socketMessages"],
	data() {
		return {
			messages: [],
			dialog: false,
			modal: false,
			op: "-1",
			cells,
			opTabs: cells.opTabs,

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
			thumbs: [],

			// 考核倒计时
			examTimes: "",

			// 考核计时器
			examFlag: "",

			// 练习倒计时
			practiceTimes: "",

			// 练习计时器
			practiceFlag: "",

			// 练习剩余时间
			surplus: "",

			// // 准确率
			// rateP: this.$store.state["entry/task"].task.accuracyRate * 100 + "%",

			// // 字符
			// charP: this.$store.state["entry/task"].task.character
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

		// 考核
		handleExam({ code, bill, block }) {
			if (code === 200) {
				this.examTime();
			}
		},

		// 练习
		handlePractice({ code, applyAt }) {
			if (code === 200) {
				let endPracticeTime = moment(applyAt, "YYYY-MM-DD HH:mm:ss").valueOf() + 432000 * 1000;
				let nowTime = moment(new Date(), "YYYY-MM-DD HH:mm:ss").valueOf();
				this.surplus = Math.floor((endPracticeTime - nowTime) / 1000);
				console.log(this.surplus);
				this.practiceTime();
			}
		},

		async confirm() {
			this.modal = false;
			let code = localStorage.get("user").code;
			await this.$store.dispatch("EXIT_PRACTICE_TASK", { code });
			this.$router.push("/main/practice/practiceChannel");
		},

		// 设置当前项目 opTabs 的状态
		setOpTabs(meta) {
			const project = localStorage.get("project");
			const perm = tools.find(this.auth.perm, { proCode: project?.code });

			if (meta && meta.key == "practiceChannel") {
				this.opTabs = [
					{
						authKey: "hasOpp",
						label: "练习",
						name: "Opp",
						value: "opp",
						show: true
					}
				];
			} else if (meta && meta.key == "opp") {
				this.opTabs = [
					{
						authKey: "hasOpp",
						label: "练习",
						name: "Opp",
						value: "opp",
						show: true,
						actived: true
					}
				];
			} else if (meta && meta.key == "examChannel") {
				this.opTabs = [
					{
						authKey: "hasOpe",
						label: "考核",
						name: "Ope",
						value: "ope",
						show: true
					}
				];
			} else if (meta && meta.key == "ope") {
				this.opTabs = [
					{
						authKey: "hasOpe",
						label: "考核",
						name: "Ope",
						value: "ope",
						show: true,
						actived: true
					}
				];
			} else {
				this.opTabs.map(op => {
					// 权限
					op["show"] = perm?.[op.authKey] ? true : false;

					// 选中当前op
					op["actived"] = op.value === meta?.path ? true : false;
				});
			}
		},

		// 弹窗状态
		handleDialog(dialog) {
			this.dialog = dialog;

			if (dialog) {
				if (this.socketMessages.length > 0) {
					this.$refs.sockets.openDialog();
				}
				this.getTaskTopInfo();
				const { code, name } = this.user;

				newWaterMark(`${code}-${name}`);
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
		navigatePage({ name }) {
			const { proCode } = this.$route.query;
			if (name != "Opp" && name != "Ope") {
				// 用于点击当前工序重新领取任务
				this.$router.replace({ name: "Channel", query: { proCode, op: -1 } });

				this.$nextTick(() => {
					this.$router.push({ name, query: { proCode } });
				});
			} else {
				this.$router.replace({ name: "Channel", query: { proCode, op: -1 } });
				this.$nextTick(() => {
					this.$router.push({ name, query: { proCode } });
				});
				// this.$router.push({ name, query: { proCode } });
			}
		},

		// 关闭弹框
		async closeDialog() {
			// this.$router.push({ name: this.$route.name });
			console.log(this.$route.name);
			if (this.$route.name == "Ope") {
				this.$router.push("/main/exam/examChannel");
			} else if (this.$route.name == "Opp") {
				this.$router.push("/main/practice/practiceChannel");
			} else if (this.$route.name == "PracticeChannel") {
				this.$router.push("/main/practice/practiceChannel");
			} else {
				this.$router.push("/main/entry/channel");
			}

			this.$refs.dialog.dialog = false;

			const { innerIp, inAppPort, outIp, outAppPort } = sessionStorage.get("task").rowInfo;
			const addr = isIntranet ? `${innerIp}:${inAppPort}` : `${outIp}:${outAppPort}`;

			this.subscribeRelease({ addr, dialog: false });
			const result = await this.$store.dispatch("ALLOCATION_ALL_TASK", { code: this.user.code, op: this.op });
			if (result.code === 200) {
				this.redirectTask("退出释放所有分块!");
			}
			sessionStorage.delete("task");
		},

		// 倒计时
		countdown() {
			this.freeTime = this.block.freeTime;

			if (!this.freeTime) return;

			this.timer = setInterval(() => {
				--this.freeTime;

				if (this.freeTime < 0) {
					this.resetCountdown();

					// 2022/5/12，泽如说倒计时为0要将用户踢出当前工序
					this.releaseTask();
					this.redirectTask();

					this.opTabs.map(entry => (entry["actived"] = false));
				}
			}, 1000);
		},

		// 清空倒计时
		resetCountdown() {
			clearInterval(this.timer);
			clearInterval(this.examFlag);
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
			const promptHeight = this.$refs.prompt?.offsetHeight || defaultPromptHeight;

			const mainHeight = height - (toolbarHeight + promptHeight);
			this.$store.commit("UPDATE_CHANNEL", { mainHeight });
		},

		// 双击单号
		onBillNum() {
			if (this.$route.name === "Opq") {
				this.$refs.opRouter.navToViewImages();
			}
		},

		// 考核倒计时
		examTime() {
			let time = 7200;
			let countdown = () => {
				let h = Math.floor(time / 3600);
				h = h < 10 ? "0" + h : h;
				let m = Math.floor((time % 3600) / 60);
				m = m < 10 ? "0" + m : m;
				let s = time % 60;
				s = s < 10 ? "0" + s : s;
				this.examTimes = h + ":" + m + ":" + s;
				time--;
				if (time == 0) {
					// alert("考核时间结束");
					clearInterval(this.examFlag);
					return false;
				}
			};
			this.examFlag = setInterval(countdown, 1000);
		},

		// 练习倒计时
		practiceTime() {
			let time = this.surplus;
			let countdown = () => {
				let day = Math.floor(time / 86400);
				let h = Math.floor((time % 86400) / 3600);
				h = h < 10 ? "0" + h : h;
				let m = Math.floor(((time % 86400) % 3600) / 60);
				m = m < 10 ? "0" + m : m;
				let s = time % 60;
				s = s < 10 ? "0" + s : s;
				this.practiceTimes = day + "天" + h + "小时" + m + "分钟" + s + "秒";
				time--;
				if (time == 0) {
					// alert("考核时间结束");
					clearInterval(this.practiceFlag);
					return false;
				}
			};
			this.practiceFlag = setInterval(countdown, 1000);
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
					this.setOpTabs(meta);

					if (meta.path || query.op) {
						this.$refs.dialog.dialog = true;
					} else {
						this.$refs.dialog.dialog = false;
					}
				});
			},
			immediate: true
		},

		// 根据录入提示动态设置 main 的高度
		"task.prompt": {
			handler(newValue) {
				this.$nextTick(() => {
					if (newValue == "") {
						this.$store.commit("UPDATE_CHANNEL", { mainHeight: 0 });
					} else {
						this.setMainHeight();
					}
				});
			},
			immediate: true
		},

		socketMessages: {
			handler(value) {
				console.log(value, this.dialog, "value");
				if (this.dialog && value.length > 0) {
					this.$refs.sockets.openDialog();
				}
			}
		}
	},
	components: {
		SocketDialog
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

.h {
	height: calc(100vh - 90px) !important;
}

.font {
	font-size: 16px;
	color: black;
	padding: 10px 10px;
}
</style>
