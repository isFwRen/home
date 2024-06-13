<template>
	<lp-dialog ref="dialog" fullscreen persistent @dialog="handleDialog">
		<router-view slot="card" @back="handleBack" @close="handleClose"></router-view>
	</lp-dialog>
</template>

<script>
import { tools } from "vue-rocket";
import DialogMixins from "@/mixins/DialogMixins";
import cells from "./cells";

export default {
	name: "CaseDialog",
	mixins: [DialogMixins],

	data() {
		return {
			cells,
			rotate: 0
		};
	},

	methods: {
		handleBack() {
			this.$router.go(-1);
			this.onClose();
		},

		handleClose() {
			this.storage.set("cases", {});
			this.$router.push({ path: "/main/PM/case" });
			this.onClose();
		},

		// 存储案件信息
		rememberCaseInfo(dialog) {
			if (dialog) {
				const rowInfo = tools.isYummy(this.rowInfo) ? this.rowInfo : {};

				const caseInfo = {
					createAt: rowInfo.CreatedAt,
					caseId: rowInfo.ID,
					billNum: rowInfo.billNum,
					proCode: rowInfo.proCode,
					billInfo: rowInfo
					//wrongNote:rowInfo.wrongNote
				};

				this.storage.set("caseInfo", caseInfo);

				this.$store.commit("UPDATE_CASE", { caseInfo });
			}
		}
	},

	watch: {
		dialog: {
			handler(dialog) {
				this.rememberCaseInfo(dialog);
			},
			immediate: true
		},

		$route: {
			handler(route) {
				const { path } = route.meta;
				if (path) {
					this.$nextTick(() => {
						this.$refs.dialog.onOpen();
					});
				}
			},
			immediate: true
		}
	}
};
</script>
