import io from "socket.io-client";
const isIntranet = location.hostname.includes("192.168") ? true : false;

const { VUE_APP_TASK_BASE_URL: taskUrl, VUE_APP_TASK_BASE_IN_URL: taskInUrl } = process.env;

const url = isIntranet ? taskInUrl : taskUrl;

const socket = io(`${url}/global-export`);

export default {
	created() {
		this.subscribeRelease();
	},

	mounted() {
		window.addEventListener("keydown", this.commonFuckEvents);
	},

	beforeDestroy() {
		window.removeEventListener("keydown", this.commonFuckEvents);
	},

	methods: {
		// 顶部信息
		async getTaskTopInfo() {
			this.user = this.storage.get("user") || {};

			const data = {
				code: this.user.code
			};

			const result = await this.$store.dispatch("GET_LP_TASKS_INFO", data);

			if (result.code === 200) {
				this.$store.commit("UPDATE_CHANNEL", { topInfo: result.data });
			}
		},

		// 删单
		async handleDelConfirm({}, form) {
			const project = this.storage.get("project");
			const body = {
				proCode: project.code,
				id: this.bill.ID,
				delRemarks: form.delRemarks + "--管理"
			};

			if (form.password !== "huiliu2022") {
				this.toasted.warning("删单密码不正确!", { position: "top" });
				return;
			}

			const result = await this.$store.dispatch("DELETE_CASE_ITEM", body);

			if (result.code === 200) {
				this.getTaskTopInfo();
				this.$refs.delDynamic.close();
				return;
			}

			if (result.code !== 200) {
				this.toasted.warning(result.msg, { position: "top" });
			}
		},

		commonFuckEvents(event) {
			const { keyCode } = event || window.event;

			switch (keyCode) {
				// 删单(F10)
				case 121:
					event.preventDefault();
					this.$refs.delDynamic.open({ title: "删单" });
					break;
			}
		},

		/**
		 * 释放分块/工序
		 * @description 若 blockId 不为空，释放分块
		 * @description 若 op 不为空，释放工序
		 */
		subscribeRelease() {
			socket.on("release", result => {
				const { billId, blockId, op } = result.data;

				if (result.code === 200) {
					if (this.bill.ID === billId) {
						this.redirectTask("任务已删除!");
						return;
					}

					if (blockId && op) {
						if (this.op === op) {
							this.redirectTask("任务已释放!");
							return;
						}
					}

					if (blockId && !op) {
						if (this.blockId === blockId) {
							this.redirectTask("任务已释放!");
							return;
						}
					}

					if (!blockId && op) {
						if (this.op === op) {
							this.redirectTask("任务已释放!");
						}
					}
				}
			});
		},

		// 释放任务
		async releaseTask() {
			const body = {
				id: this.block.ID,
				op: this.op,
				code: ""
			};

			const result = await this.$store.dispatch("TASK_ALLOCATION_TASK", body);

			if (result.code === 200) {
				this.redirectTask("超时，释放分块!");
			}
		},

		// 退出当前工序
		redirectTask(message) {
			this.toasted.warning(message, { position: "top" });
			this.$router.replace({ path: "/main/entry/channel", query: { op: -1 } });
		}
	}
};
