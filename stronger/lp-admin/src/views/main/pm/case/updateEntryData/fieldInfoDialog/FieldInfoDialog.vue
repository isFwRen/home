<template>
	<div class="field-info-dialog">
		<lp-dialog
			ref="dialog"
			title="详情"
			transition="dialog-bottom-transition"
			width="1080"
			@dialog="handleDialog"
		>
			<div slot="main" class="pt-8 pb-4 main">
				<div class="z-table">
					<div class="row">
						<div class="col z-title">字段编码</div>
						<div class="col">{{ field.code }}</div>
						<div class="col z-title">字段名称</div>
						<div class="col">{{ field.name }}</div>
						<div class="col z-title">客户投诉</div>
					</div>

					<div class="row">
						<div class="col z-title">最终数据</div>
						<div class="col">{{ field.finalValue }}</div>
						<div class="text-center col">
							<z-btn icon @click="onCorrect">
								<v-icon>mdi-circle-edit-outline</v-icon>
							</z-btn>
						</div>
					</div>

					<div class="row">
						<div class="col z-title">初审录入</div>
						<div class="col">{{ block.op0Code }}</div>
						<div class="col z-title">初审录入数据</div>
						<div class="col">{{ field.op0Value }}</div>
						<div class="col z-title">初审录入时间</div>
						<div class="col">
							{{ block.op0SubmitAt | dateFormat("YYYY-MM-DD HH:mm:ss") }}
						</div>
						<div class="col"></div>
					</div>

					<div class="row">
						<div class="col z-title">一码录入</div>
						<div class="col">{{ block.op1Code }}</div>
						<div class="col z-title">一码录入数据</div>
						<div class="col">{{ field.op1Value }}</div>
						<div class="col z-title">一码录入时间</div>
						<div class="col">
							{{ block.op1SubmitAt | dateFormat("YYYY-MM-DD HH:mm:ss") }}
						</div>
						<div class="col"></div>
					</div>

					<div class="row">
						<div class="col z-title">二码录入</div>
						<div class="col">{{ block.op2Code }}</div>
						<div class="col z-title">二码录入数据</div>
						<div class="col">{{ field.op2Value }}</div>
						<div class="col z-title">二码录入时间</div>
						<div class="col">
							{{ block.op2SubmitAt | dateFormat("YYYY-MM-DD HH:mm:ss") }}
						</div>
						<div class="col"></div>
					</div>

					<div class="row">
						<div class="col z-title">问题件录入</div>
						<div class="col">{{ block.opqCode }}</div>
						<div class="col z-title">问题件录入数据</div>
						<div class="col">{{ field.opqValue }}</div>
						<div class="col z-title">问题件录入时间</div>
						<div class="col">
							{{ block.opqSubmitAt | dateFormat("YYYY-MM-DD HH:mm:ss") }}
						</div>
						<div class="col"></div>
					</div>

					<div class="row">
						<div class="col z-title">录入提示</div>
						<div class="col">{{ fieldConf.prompt }}</div>
					</div>
				</div>
			</div>
		</lp-dialog>

		<z-dynamic-form
			ref="dynamic"
			:formId="formId"
			title="填写正确数据"
			:fieldList="cells.fields"
			:config="config"
			:detail="detailInfo"
			:width="600"
			@search:responsibleCode="handleCodeSearch"
			@change:responsibleCode="handleCodeChange"
			@confirm="handleConfirm"
		>
		</z-dynamic-form>
	</div>
</template>

<script>
import moment from "moment";
import { tools } from "vue-rocket";
import DialogMixins from "@/mixins/DialogMixins";
import cells from "./cells";

export default {
	name: "FieldInfoDialog",
	mixins: [DialogMixins],

	props: {
		info: {
			type: Object,
			default: () => {}
		}
	},

	data() {
		return {
			formId: "FieldInfoDialog",
			cells,
			block: {},
			field: {},
			fieldConf: {},
			config: {
				responsibleCode: {
					items: []
				}
			},
			detailInfo: {
				editDate: moment().format("YYYY-MM-DD"),
				month: moment().format("YYYY-MM")
			}
		};
	},

	methods: {
		async getCaseFieldInfo() {
			let [block, field, fieldConf] = [{}, {}, {}];

			const result = await this.$store.dispatch("GET_CASE_FIELD_INFO", this.info);

			if (result.code === 200) {
				block = result.data.block || {};
				field = result.data.field || {};
				fieldConf = result.data.fieldConf || {};
			} else {
				this.toasted.error(result.msg);
			}

			[this.block, this.field, this.fieldConf] = [block, field, fieldConf];
		},

		handleDialog(dialog) {
			dialog && this.getCaseFieldInfo();
		},

		onCorrect() {
			this.$refs.dynamic.open();
		},

		// 责任人工号
		async handleCodeSearch(value) {
			const params = {
				pageIndex: 1,
				pageSize: 1000,
				code: value,
				status: true
			};

			const result = await this.$store.dispatch("STAFF_GET_USER_LIST", params);

			if (result.code === 200) {
				const items = [];

				result.data.list.map(item => {
					items.push({ label: item.code, value: item.code, name: item.nickName });
				});

				this.config = {
					...this.config,
					responsibleCode: {
						items
					}
				};
			}
		},

		handleCodeChange(value) {
			const { items } = this.config.responsibleCode;
			const { name } = tools.find(items, { value });

			this.detailInfo = {
				...this.detailInfo,
				responsibleName: name
			};
		},

		async handleConfirm({}, form) {
			const result = await this.$store.dispatch("CASE_POST_FEEDBACK_VALUE", form);

			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				this.$refs.dynamic.close();
			}
		}
	}
};
</script>

<style scoped lang="scss">
.row {
	display: grid;

	&:first-child {
		grid-template-columns: repeat(4, 3fr) 1fr;
		border-top: 1px solid #e8eaec;
	}

	&:nth-child(2) {
		grid-template-columns: 1.5fr 10.5fr 1fr;
	}

	&:nth-child(3),
	&:nth-child(4),
	&:nth-child(5),
	&:nth-child(6) {
		grid-template-columns: 1.5fr 1.5fr 1.5fr 4.5fr 1.5fr 1.5fr 1fr;
	}

	&:last-child {
		grid-template-columns: 1.5fr 11.5fr;
	}

	&:nth-child(odd) {
		background-color: #f8f8f9;
	}

	.col {
		padding: 6px 10px;
		border-right: 1px solid #e8eaec;
		border-bottom: 1px solid #e8eaec;
	}

	.col:first-child {
		border-left: 1px solid #e8eaec;
	}

	.z-title {
		font-weight: bold;
	}
}
</style>
