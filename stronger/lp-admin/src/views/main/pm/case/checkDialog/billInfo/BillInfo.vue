<template>
	<div>
		<v-expansion-panels multiple v-model="flag">
			<v-expansion-panel v-for="(item, i) in title" :key="i">
				<v-expansion-panel-header color="#e9f6ff" class="header">
					{{ item }}
				</v-expansion-panel-header>
				<v-expansion-panel-content v-if="i == 0">
					<vxe-table
						ref="xTable1"
						:data="desserts1"
						border
						stripe
						align="center"
						style="margin-top: 15px"
						min-height="90"
						class="mytable-scrollbar"
						:edit-config="{ trigger: 'dblclick', mode: 'cell' }"
						@edit-closed="editClosedEvent1"
					>
						<vxe-column type="seq" width="60" title="序号"></vxe-column>
						<template v-for="item in cells.headers1">
							<!-- 时间 BEGIN -->
							<vxe-column
								v-if="item.value === 'CreatedAt'"
								:field="item.value"
								:title="item.text"
								:key="item.value"
								:width="item.width"
								:sortable="item.sortable"
								:fixed="item.fixed"
							>
								<template #default="{ row }">
									{{ row.CreatedAt | dateFormat("YYYY-MM-DD") }}
								</template>
							</vxe-column>
							<!-- 时间 END -->
							<vxe-column
								v-else
								:field="item.value"
								:fixed="item.fixed"
								:title="item.text"
								:key="item.value"
								:width="item.width"
								:sortable="item.sortable"
								:edit-render="{ autofocus: '.vxe-input--inner' }"
							>
								<template #edit="{ row }">
									<vxe-input v-model="row[item.value]" type="text"></vxe-input>
								</template>
							</vxe-column>
						</template>
					</vxe-table>
				</v-expansion-panel-content>
				<v-expansion-panel-content v-else>
					<vxe-table
						ref="xTable2"
						:data="desserts2"
						border
						stripe
						align="center"
						style="margin-top: 15px; margin-bottom: 15px"
						:radio-config="{ highlight: true }"
						@radio-change="radioChangeEvent"
						class="mytable-scrollbar"
						:edit-config="{ trigger: 'dblclick', mode: 'cell' }"
						@edit-closed="editClosedEvent2"
					>
						<vxe-column type="radio" width="60"> </vxe-column>
						<vxe-column type="seq" width="60" title="序号"></vxe-column>
						<template v-for="item in cells.headers2">
							<!-- 时间 BEGIN -->
							<vxe-column
								v-if="item.value === 'CreatedAt'"
								:field="item.value"
								:title="item.text"
								:key="item.value"
								:width="item.width"
								:sortable="item.sortable"
								:fixed="item.fixed"
							>
								<template #default="{ row }">
									{{ row.CreatedAt | dateFormat("YYYY-MM-DD") }}
								</template>
							</vxe-column>
							<!-- 时间 END -->
							<vxe-column
								v-else
								:field="item.value"
								:fixed="item.fixed"
								:title="item.text"
								:key="item.value"
								:width="item.width"
								:sortable="item.sortable"
								:edit-render="{ autofocus: '.vxe-input--inner' }"
							>
								<template #edit="{ row }">
									<vxe-input v-model="row[item.value]" type="text"></vxe-input>
								</template>
							</vxe-column>
						</template>
					</vxe-table>
					<vxe-form
						:forms="cells.form2"
						:formRules="cells.formRule"
						:belongModule="'billInfo'"
						:belongModuleForm="'info3'"
						:cloudData="billInfo.info3"
					></vxe-form>
					<vxe-table
						ref="xTable3"
						:data="desserts3"
						border
						stripe
						align="center"
						style="margin-top: 15px; margin-bottom: 15px"
						class="mytable-scrollbar"
						:edit-config="{ trigger: 'dblclick', mode: 'cell' }"
						@edit-closed="editClosedEvent3"
					>
						<vxe-column type="seq" width="60" title="序号"></vxe-column>
						<vxe-column type="expand" width="80" title="解析">
							<template #content="{ rowIndex }">
								<div class="expand-wrapper">
									<vxe-form
										:forms="cells.form3"
										:belongModule="'billInfo'"
										:belongModuleForm="'info' + 4 + rowIndex"
										:cloudData="billInfo['info' + 4 + rowIndex]"
									></vxe-form>
								</div>
							</template>
						</vxe-column>
						<template v-for="item in cells.headers3">
							<!-- 时间 BEGIN -->
							<vxe-column
								v-if="item.value === 'CreatedAt'"
								:field="item.value"
								:title="item.text"
								:key="item.value"
								:width="item.width"
								:fixed="item.fixed"
							>
								<template #default="{ row }">
									{{ row.CreatedAt | dateFormat("YYYY-MM-DD") }}
								</template>
							</vxe-column>
							<!-- 时间 END -->
							<vxe-column
								v-else
								:field="item.value"
								:fixed="item.fixed"
								:title="item.text"
								:key="item.value"
								:width="item.width"
								:edit-render="{ autofocus: '.vxe-input--inner' }"
							>
								<template #edit="{ row }">
									<vxe-input v-model="row[item.value]" type="text"></vxe-input>
								</template>
							</vxe-column>
						</template>
					</vxe-table>
				</v-expansion-panel-content>
			</v-expansion-panel>
		</v-expansion-panels>
	</div>
</template>

<script>
import cells from "./cells";
export default {
	props: {
		BillInfo: {
			type: Object,
			default: () => {}
		}
	},
	data() {
		return {
			flag: [0, 1],
			title: ["票据信息汇总", "账单详细信息"],
			cells,
			desserts1: [
				{
					invoiceCount: "",
					originalAmount: "",
					adjustAmount: "",
					deductionAmount: "",
					reimbursementAmount: ""
				}
			],
			desserts2: cells.desserts2,
			desserts3: cells.desserts3,
			selectRow: {},
			billInfo: {}
		};
	},

	mounted() {
		this.$refs.xTable2[0].setRadioRow(this.desserts2[0]);
	},
	methods: {
		radioChangeEvent({ row }) {
			this.selectRow = row;
			console.log("单选事件", this.selectRow);
		},
		editClosedEvent1() {
			const $table = this.$refs.xTable1;
			let content = sessionStorage.get("checkForm");
			content.billInfo.desserts1 = $table.getData();
			sessionStorage.set("checkForm", content);
		},
		editClosedEvent2() {
			const $table = this.$refs.xTable2;
			let content = sessionStorage.get("checkForm");
			content.billInfo.desserts2 = $table.getData();
			sessionStorage.set("checkForm", content);
		},
		editClosedEvent3() {
			const $table = this.$refs.xTable3;
			let content = sessionStorage.get("checkForm");
			content.billInfo.desserts3 = $table.getData();
			sessionStorage.set("checkForm", content);
		}
	},

	watch: {
		BillInfo: {
			handler(newValue) {
				this.billInfo = JSON.parse(JSON.stringify(newValue));
				if (this.billInfo.desserts1 && this.billInfo.desserts1.length != 0) {
					this.desserts1 = this.billInfo.desserts1;
				}
				if (this.billInfo.desserts2 && this.billInfo.desserts2.length != 0) {
					this.desserts2 = this.billInfo.desserts2;
				}
				if (this.billInfo.desserts3 && this.billInfo.desserts3.length != 0) {
					this.desserts3 = this.billInfo.desserts3;
				}
				// console.log("billInfo", newValue);
			},
			immediate: true,
			deep: true
		}
	},

	components: {
		"vxe-form": () => import("../vxeForm")
	}
};
</script>

<style lang="scss">
.header {
	color: #007aff;
	font-weight: bolder;
	font-size: 15px;
}
.v-expansion-panel--active > .v-expansion-panel-header {
	min-height: 0px;
}

.expand-wrapper {
	padding: 20px;
}
</style>