<template>
	<div class="insuran">
		<div class="count_header">诊疗信息</div>
		<div class="grid">
			<div class="cell" v-for="(value, key) in formsData" :key="key">
				<input type="text" v-model="formsData[key]" />
			</div>
		</div>

		<div class="forms_wrap">
			<vxe-form :data="formStore" title-width="115" :rules="cells.formRule1">
				<vxe-form-item
					:title="el.title"
					:field="el.field"
					v-for="(el, index) in cells.formList1"
					:key="index"
					:span="8"
				>
					<template #default="{ data }">
						<vxe-input
							v-if="el.name === 'input'"
							v-model="data[el.field]"
							:disabled="el.disabled"
						></vxe-input>
						<vxe-input
							v-else-if="el.name == 'date'"
							v-model="data[el.field]"
							type="date"
							transfer
						></vxe-input>
						<vxe-select
							v-else-if="el.name == 'select'"
							v-model="data[el.field]"
							transfer
						>
							<vxe-option
								v-for="item in el.options"
								:key="item.value"
								:value="item.value"
								:label="item.label"
							></vxe-option>
						</vxe-select>
					</template>
				</vxe-form-item>
			</vxe-form>
		</div>

		<div class="tables_wrap">
			<vxe-table
				show-overflow
				ref="tableRef"
				:edit-config="{ trigger: 'dblclick', mode: 'row' }"
				:data="tableData"
			>
				<vxe-column type="seq" title="序号"></vxe-column>
				<vxe-column
					field="name1"
					title="诊断类型"
					:edit-render="{ autofocus: '.vxe-input--inner' }"
				>
					<template #edit="{ row }">
						<vxe-input v-model="row.name1"></vxe-input>
					</template>
				</vxe-column>
				<vxe-column field="name2" title="西医诊断" :edit-render="{}">
					<template #edit="{ row }">
						<vxe-input v-model="row.name2"></vxe-input>
					</template>
				</vxe-column>
				<vxe-column field="name3" title="疾病代码" :edit-render="{}">
					<template #edit="{ row }">
						<vxe-input v-model="row.name3"></vxe-input>
					</template>
				</vxe-column>
				<vxe-column field="name4" title="中医诊断" :edit-render="{}">
					<template #edit="{ row }">
						<vxe-input v-model="row.name4"></vxe-input>
					</template>
				</vxe-column>
				<vxe-column field="name5" title="疾病代码" :edit-render="{}">
					<template #edit="{ row }">
						<vxe-input v-model="row.name5"></vxe-input>
					</template>
				</vxe-column>
				<vxe-column field="name6" title="手术及操作名称" :edit-render="{}">
					<template #edit="{ row }">
						<vxe-input v-model="row.name6"></vxe-input>
					</template>
				</vxe-column>
				<vxe-column field="name7" title="手术及操作代码" :edit-render="{}">
					<template #edit="{ row }">
						<vxe-input v-model="row.name7"></vxe-input>
					</template>
				</vxe-column>
				<vxe-column field="name8" title="手术及操作日期" :edit-render="{}">
					<template #edit="{ row }">
						<vxe-input v-model="row.name8"></vxe-input>
					</template>
				</vxe-column>
			</vxe-table>
			<div class="plus_icons">
				<v-icon size="28" color="#007AFF" class="curpor" @click="onAdd">
					mdi-plus-box</v-icon
				>
			</div>
		</div>

		<div class="content_list">
			<div class="count_header">伤残鉴定信息</div>
			<div class="grid">
				<div class="cell" v-for="(value, key) in formsData1" :key="key">
					<input type="text" v-model="formsData1[key]" />
				</div>
			</div>

			<div class="forms_wrap">
				<vxe-form :data="forms" title-width="115" :rules="cells.formRule2">
					<vxe-form-item
						:title="el.title"
						:field="el.field"
						v-for="(el, index) in cells.formList2"
						:key="index"
						:span="8"
					>
						<template #default="{ data }">
							<vxe-input
								v-if="el.name === 'input'"
								v-model="data[el.field]"
								:disabled="el.disabled"
							></vxe-input>
							<vxe-input
								v-else-if="el.name == 'date'"
								v-model="data[el.field]"
								type="date"
								transfer
							></vxe-input>
							<vxe-select
								v-else-if="el.name == 'select'"
								v-model="data[el.field]"
								transfer
							>
								<vxe-option
									v-for="item in el.options"
									:key="item.value"
									:value="item.value"
									:label="item.label"
								></vxe-option>
							</vxe-select>
						</template>
					</vxe-form-item>
				</vxe-form>
			</div>
		</div>

		<div class="content_info">
			<div class="count_header">身故信息</div>
			<div class="forms_wrap">
				<vxe-form :data="forms3" title-width="115">
					<vxe-form-item
						:title="el.title"
						:field="el.field"
						v-for="(el, index) in cells.formList3"
						:key="index"
						:span="el.span"
					>
						<template #default="{ data }">
							<vxe-input
								v-if="el.name === 'input'"
								v-model="data[el.field]"
								:disabled="el.disabled"
							></vxe-input>
							<vxe-input
								v-else-if="el.name == 'date'"
								v-model="data[el.field]"
								type="date"
								transfer
							></vxe-input>
							<vxe-select
								v-else-if="el.name == 'select'"
								v-model="data[el.field]"
								transfer
							>
								<vxe-option
									v-for="item in el.options"
									:key="item.value"
									:value="item.value"
									:label="item.label"
								></vxe-option>
							</vxe-select>
						</template>
					</vxe-form-item>
				</vxe-form>
			</div>
		</div>
	</div>
</template>

<script>
import { sessionStorage } from "vue-rocket";
import cells from "./cells";
export default {
	props: {
		InsuranceInfo: {
			type: Object,
			default: () => {}
		}
	},

	data() {
		return {
			content: {},
			forms: {},
			forms3: {},
			formsData: {
				name1: "序号",
				name2: "理赔类型",
				name3: "医疗类型",
				name4: "就诊日期",
				name5: "就诊医院",
				name6: "1",
				name7: "医疗",
				name8: "住院",
				name9: "2323-10-10",
				name10: "茂名市茂南区人民医院"
			},
			formsData1: {
				name1: "序号",
				name2: "伤残鉴定书编号",
				name3: "鉴定日期",
				name4: "鉴定标准",
				name5: "伤残比例",
				name6: "1",
				name7: "1111",
				name8: "2023-10-10",
				name9: "蔷薇色的花",
				name10: "20%"
			},
			tableData: [
				{
					name1: "西医",
					name2: "下肢的其他损失，水平未特指",
					name3: "T13",
					name4: 1,
					name5: 1,
					name6: 1,
					name7: 1,
					name8: "2024-03-04",
					operation: ""
				}
			],
			cells,
			formStore: {
				name2: "珠海市妇幼保健院"
			}
		};
	},
	created() {
		this.initData();

		this.content = sessionStorage.get("checkForm");
		this.content.insuranceInfo.formsData = this.formsData;
		this.content.insuranceInfo.formsData1 = this.formsData1;
		this.content.insuranceInfo.tableData = this.tableData;
		this.content.insuranceInfo.forms = this.forms;
		this.content.insuranceInfo.forms3 = this.forms3;
		this.content.insuranceInfo.formStore = this.formStore;
		sessionStorage.set("checkForm", this.content);
	},
	methods: {
		initData() {

			if (this.InsuranceInfo.tableData && this.InsuranceInfo.tableData.length > 0) {
				this.tableData = this.InsuranceInfo.tableData;
			}
			if (
				this.InsuranceInfo.formsData &&
				Object.keys(this.InsuranceInfo.formsData).length > 0
			) {
				this.formsData = this.InsuranceInfo.formsData;
			}
			if (
				this.InsuranceInfo.formsData1 &&
				Object.keys(this.InsuranceInfo.formsData1).length > 0
			) {
				this.formsData1 = this.InsuranceInfo.formsData1;
			}
			if (this.InsuranceInfo.forms && Object.keys(this.InsuranceInfo.forms).length > 0) {
				this.forms = this.InsuranceInfo.forms;
			}
			if (this.InsuranceInfo.forms3 && Object.keys(this.InsuranceInfo.forms3).length > 0) {
				this.forms3 = this.InsuranceInfo.forms3;
			}
			if (
				this.InsuranceInfo.formStore &&
				Object.keys(this.InsuranceInfo.formStore).length > 0
			) {
				this.formStore = this.InsuranceInfo.formStore;
			}
		},
		onAdd() {
			this.tableData.push({
				name1: 1,
				name2: 1,
				name3: 1,
				name4: 1,
				name5: 1,
				name6: 1,
				name7: 1,
				name8: 1,
				operation: ""
			});
		},
		updateSesstion() {
			sessionStorage.set("checkForm", this.content);
		}
	},

	watch: {
		formsData: {
			handler(value) {
				this.content.insuranceInfo.formsData = value;
				this.updateSesstion();
			},
			deep: true
		},
		formsData1: {
			handler(value) {
				this.content.insuranceInfo.formsData1 = value;
				this.updateSesstion();
			},
			deep: true
		},
		tableData: {
			handler(value) {
				this.content.insuranceInfo.tableData = value;
				this.updateSesstion();
			},
			deep: true
		},
		forms: {
			handler(value) {
				this.content.insuranceInfo.forms = value;
				this.updateSesstion();
			},
			deep: true
		},
		forms3: {
			handler(value) {
				this.content.insuranceInfo.forms3 = value;
				this.updateSesstion();
			},
			deep: true
		},
		formStore: {
			handler(value) {
				this.content.formStore = value;
				this.updateSesstion();
			},
			deep: true
		}
	}
};
</script>

<style lang="scss">
.blue_color {
	color: #007aff;
	font-weight: 600;
	cursor: pointer;
}
.insuran {
	padding: 10px;
	background-color: #fff;
	.count_header {
		background-color: #e9f6ff;
		padding: 10px 10px;
		color: #007aff;
		font-weight: 600;
		font-size: 15px;
	}

	.grid {
		display: grid;
		margin-top: 10px;
		grid-template-columns: repeat(5, 1fr);
		font-size: 14px;
		border-bottom: 1px solid #eee;
		border-right: 1px solid #eee;
		.cell {
			border-top: 1px solid #eee;
			border-left: 1px solid #eee;
			padding: 5px 5px;
			input {
				display: block;
				width: 100%;
				height: 100%;
				padding: 5px;
				border: 1px solid transparent;
				border-radius: 10px;
			}
			input:focus {
				outline: #eee;
				border: 1px solid #eee;
			}
		}
	}

	.forms_wrap {
		padding: 10px;
	}
	.tables_wrap {
		font-size: 14px;
	}
}
.plus_icons {
	padding: 10px 0 20px 0;
}

.curpor {
	cursor: pointer;
}
</style>
