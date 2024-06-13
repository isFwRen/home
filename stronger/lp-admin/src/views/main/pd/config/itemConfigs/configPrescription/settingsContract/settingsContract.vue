<template>
	<div class="settings-contrach">
		<div class="z-flex justify-space-between">
			<h5 class="font-weight-bold">基础设置</h5>
			<v-btn color="primary" @click="onAdd">
				<v-icon class="text-h6">mdi-plus</v-icon>
				新增
			</v-btn>
		</div>

		<div class="wrapper">
			<div
				class="list z-flex justify-space-between"
				v-for="(item, index) in datSource"
				:key="index"
			>
				<div class="item start-time">
					<div class="header">起始时间</div>
					<div class="value">{{ item.contractStartTime }}</div>
				</div>
				<div class="item end-time">
					<div class="header">结束时间</div>
					<div class="value">{{ item.contractEndTime }}</div>
				</div>
				<div class="item claim-type">
					<div class="header">理赔类型</div>
					<div class="value">{{ computedType(item.claimType) }}</div>
				</div>
				<div class="item out-start-time">
					<div class="header">时效外开始时间</div>
					<div class="value">{{ item.contractOutsideStartTime }}</div>
				</div>
				<div class="item out-late-time">
					<div class="header">时效外最晚时间</div>
					<div class="value">{{ item.contractOutsideEndTime }}</div>
				</div>
				<div class="item assessment-requirements">
					<div class="header">考核要求(min)</div>
					<div class="value">{{ item.requirementsTime }}</div>
				</div>
				<div class="item option">
					<div class="header">操作</div>
					<div class="value edit-btns z-flex justify-center">
						<div>
							<v-btn text @click="onEditItem(item)"> 编辑 </v-btn>
						</div>
						<div>
							<v-btn text color="red" @click="onDeleteItem(item)"> 删除 </v-btn>
						</div>
					</div>
				</div>
			</div>
		</div>

		<!-- 新增/编辑 BEGIN -->
		<z-dynamic-form
			ref="dynamic"
			:formId="formId"
			:detail="detailInfo"
			:fieldList="cells.fields"
			@confirm="handleConfirm"
		></z-dynamic-form>
		<!-- 新增/编辑 END -->
	</div>
</template>
<script>
import { mapGetters } from "vuex";
import cells from "./cells";
import TableMixins from "@/mixins/TableMixins";

export default {
	name: "SettingsContract",
	mixins: [TableMixins],
	data() {
		return {
			cells,
			formId: "SettingsContract",
			datSource: []
		};
	},
	mounted() {
		this.getDataList();
	},
	computed: {
		...mapGetters(["project", "config"])
	},
	methods: {
		computedType(type) {
			let label = "";
			cells.claimTypes.forEach(item => {
				if (item.value === type) {
					label = item.label;
				}
			});
			return label;
		},
		async getDataList() {
			const body = {
				code: this.project.code,
			};
			const result = await this.$store.dispatch("GET_CONFIG_CONTRACT_LIST", body);
			this.toasted.dynamic(result.msg, result.code);
			if (result.code === 200) {
				this.datSource = result.data.list;
			} else {
				this.datSource = [];
			}
		},
		handleConfirm(effect, form) {
			form = {
				proId: this.config.proId,
				id: this.detailInfo.ID,
				...form,
				code: this.project.code
			};
			if (effect.status === -1) {
				this.addListItem(form);
			}
			if (effect.status === 1) {
				this.updateListItem(form);
			}
		},
		async addListItem(body) {
			const result = await this.$store.dispatch("ADD_CONFIG_CONTRACT", body);
			this.toasted.dynamic(result.msg, result.code);
			if (result.code === 200) {
				this.getDataList();
				this.$refs.dynamic.close();
			}
		},
		onEditItem(item) {
			this.getDetail(item);
			this.$refs.dynamic.open({ title: "编辑", status: 1 });
		},
		async updateListItem(form) {
			const result = await this.$store.dispatch("EDIT_CONFIG_CONTRACT", form);
			this.toasted.dynamic("修改成功", result.code);
			if (result.code === 200) {
				this.getDataList();
				this.$refs.dynamic.close();
			}
		},
		async onDeleteItem(item) {
			const body = {
				ids: [item.ID]
			};
			const result = await this.$store.dispatch("DELETE_CONFIG_CONTRACT", body);
			this.toasted.dynamic("删除成功", result.code);
			if (result.code === 200) {
				this.getDataList();
			}
		},
		onAdd() {
			this.getDetail({});
			this.$refs.dynamic.open({ title: "新增", status: -1 });
		}
	}
};
</script>
<style scoped lang="scss">
.list {
	margin-top: 20px;
	border: 1px solid #eee;
	.item {
		flex: 1;
		text-align: center;
		min-width: 120px;
	}
	.header {
		padding: 10px 0;
		background: #f5f5f5;
	}
	.value {
		padding: 10px 0;
	}
}
</style>
