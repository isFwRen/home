<template>
	<div class="customer">
		<div class="mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col cols="2">
					<z-select
						:formId="searchFormId"
						formKey="proCode"
						hideDetails
						label="项目编码"
						:options="auth.proItems"
						:defaultValue="project.code"
						@change="changePro"
					></z-select>
				</v-col>

				<v-col cols="3">
					<z-date-picker
						:formId="searchFormId"
						formKey="time"
						hideDetails
						label="发送时间"
						range
						z-index="10"
						@input="changePro"
					></z-date-picker>
				</v-col>

				<v-col cols="2">
					<z-select
						:formId="searchFormId"
						formKey="msgType"
						clearable
						hideDetails
						label="消息类型"
						:options="typeItems"
						:defaultValue="project.code"
						@change="changePro"
					></z-select>
				</v-col>

				<v-col cols="2">
					<z-select
						:formId="searchFormId"
						formKey="status"
						clearable
						hideDetails
						label="消息状态"
						:options="statusItems"
						:defaultValue="project.code"
						@change="changePro"
					></z-select>
				</v-col>

				<z-btn class="pb-3 px-3" color="primary" @click="onSearch">
					<v-icon size="20">mdi-magnify</v-icon>
					查询
				</z-btn>
			</v-row>
		</div>

		<div class="table">
			<vxe-table
				:data="desserts"
				:border="tableBorder"
				:size="tableSize"
				@checkbox-all="handleSelectAll"
				@checkbox-change="handleSelectChange"
			>
				<vxe-column type="seq" title="序号" width="60"></vxe-column>

				<template v-for="item in cells.headers">
					<vxe-column
						v-if="item.value === 'options'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default="{ row }">
							<div class="py-2 z-flex">
								<z-btn
									v-if="row.status == 1"
									color="primary"
									depressed
									small
									@click="handleReply(row)"
								>
									回复
								</z-btn>

								<z-btn
									v-else
									color="primary"
									depressed
									outlined
									small
									@click="handleView(row)"
								>
									查看
								</z-btn>
							</div>
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'replyTime'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default="{ row }">
							<span>{{ row[item.value] | dateFormat("YYYY-MM-DD HH:mm:ss") }}</span>
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'status'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default="{ row }">
							<span>{{ row[item.value] | chineseStatus }}</span>
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'msgType'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default="{ row }">
							<span>{{ row[item.value] | chineseType }}</span>
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'sendTime'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default="{ row }">
							<span>{{ row[item.value] | dateFormat("YYYY-MM-DD HH:mm:ss") }}</span>
						</template>
					</vxe-column>

					<vxe-column
						v-else
						:key="item.value"
						:field="item.value"
						:title="item.text"
						:width="item.width"
					>
					</vxe-column>
				</template>
			</vxe-table>

			<z-pagination
				:options="pageSizes"
				@page="handlePage"
				:total="pagination.total"
			></z-pagination>
		</div>

		<!-- 回复 BEGIN -->
		<z-dynamic-form
			ref="dynamic"
			:detail="detailInfo"
			:fieldList="status === 'reply' ? cells.replyFields : cells.viewFields"
			:config="{
				isReply: {
					disabled: status !== 'reply',
					items: cells.isReplyItems,
					mutex: [
						{
							formKey: 'expectNum',
							always: false,
							includes: [false]
						}
					]
				},

				expectNum: {
					disabled: status !== 'reply'
				}
			}"
			:cancelProps="{
				text: status === 'view' ? '关闭' : '取消'
			}"
			:confirmProps="{
				visible: status === 'reply',
				text: '回复'
			}"
			@confirm="handleConfirm"
		></z-dynamic-form>
		<!-- 回复 END -->
	</div>
</template>

<script>
import moment from "moment";
import { mapGetters } from "vuex";
import { tools } from "vue-rocket";
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";

let [typeItems, statusItems] = [[], []];

export default {
	name: "Customer",
	mixins: [TableMixins],

	data() {
		return {
			formId: "customer",
			cells,
			manual: true,
			dispatchList: "CUSTOMER_GET_LIST",
			status: void 0,
			statusItems: [],
			typeItems: []
		};
	},

	computed: {
		...mapGetters(["auth", "project"])
	},

	watch: {
		sabayon: {
			handler(sabayon) {
				this.$store.commit("GLOBAL_NOTIFICATION_UPDATE_ITEM", {
					count: sabayon.data?.maxOrder
				});
			},
			deep: true
		}
	},

	created() {
		this.getConstants();
	},

	methods: {
		changePro(value) {
			console.log(value);
		},

		handleReply(row) {
			this.status = "reply";
			const sendTime = moment(row.sendTime).format("YYYY-MM-DD hh:mm:ss");
			const replyTime = moment(row.sendTime).format("YYYY-MM-DD hh:mm:ss");
			this.detailInfo = { ...row, sendTime, replyTime };

			this.$refs.dynamic.open({ title: "消息回复", id: row.ID });
		},

		handleView(row) {
			this.status = "view";
			const sendTime = moment(row.sendTime).format("YYYY-MM-DD hh:mm:ss");
			const replyTime = moment(row.sendTime).format("YYYY-MM-DD hh:mm:ss");
			this.detailInfo = { ...row, sendTime, replyTime };

			this.$refs.dynamic.open({ title: "消息查看" });
		},

		async handleConfirm({ id }, form) {
			const { proCode } = this.forms[this.searchFormId];

			const body = {
				proCode,
				id,
				isReply: form.isReply,
				expectNum: form.expectNum
			};

			const result = await this.$store.dispatch("CUSTOMER_POST_REPLY", body);

			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				this.$refs.dynamic.close();
				this.getList();
			}
		},

		async getConstants() {
			const result = await this.$store.dispatch("MESSAGE_GET_CONSTANTS");

			if (result.code === 200) {
				// 消息类型
				for (let key in result.data.type) {
					const text = result.data.type[key];

					this.typeItems.push({
						label: text,
						value: +key
					});
				}

				// 消息状态
				for (let key in result.data.status) {
					const text = result.data.status[key];

					this.statusItems.push({
						label: text,
						value: +key
					});
				}

				typeItems = this.typeItems;
				statusItems = this.statusItems;
			}
		}
	},

	filters: {
		chineseType: value => {
			const result = tools.find(typeItems, { value });
			return result?.label || "-";
		},

		chineseStatus: value => {
			const result = tools.find(statusItems, { value });
			return result?.label || "-";
		}
	}
};
</script>
