<template>
	<div class="pt-salary">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col
					v-for="(item, index) in cells.fields"
					:key="`entry_filters_${index}`"
					:cols="item.cols"
				>
					<template v-if="item.inputType === 'input'">
						<z-text-field
							:formId="searchFormId"
							:formKey="item.formKey"
							:clearable="item.clearable"
							:hideDetails="item.hideDetails"
							:hint="item.hint"
							:label="item.label"
							:suffix="item.suffix"
							:defaultValue="item.defaultValue"
						>
						</z-text-field>
					</template>
					<template v-else-if="item.inputType === 'date'">
						<z-date-picker
							:formId="searchFormId"
							:formKey="item.formKey"
							:hideDetails="item.hideDetails"
							:hint="item.hint"
							:label="item.label"
							:pickerType="item.pickerType"
							:options="item.options"
							:range="item.range"
							:suffix="item.suffix"
							z-index="10"
							:defaultValue="item.defaultValue"
						></z-date-picker>
					</template>
				</v-col>
				<div class="z-flex">
					<z-btn
						:formId="formId"
						btnType="validate"
						class="pb-3 pl-3"
						color="primary"
						@click="onSearch"
					>
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>
				</div>
			</v-row>
		</div>
		<div class="z-flex pb-6 btns">
			<z-file-input
				:formId="formId"
				formKey="file"
				label="工资表"
				accept=".xlsx"
				:action="action"
				hide-details
				placeholder="点击文件或将文件拖拽到这里"
				class="mt-n3 mr-3"
				:deleteIcon="false"
				:headers="fileHeaders"
				parcel
				multiple
				width="120"
				@click="getUserInfo"
			>
			</z-file-input>
			<z-btn
				v-for="item of cells.btns"
				:key="item.icon"
				:class="item.class"
				:color="item.color"
				small
				outlined
				@click="onBack2(item)"
			>
				<v-icon class="text-h6">{{ item.icon }}</v-icon>
				{{ item.text }}
			</z-btn>
		</div>
		<div class="table pt-salary-table">
			<vxe-table :border="tableBorder" :data="desserts" :size="tableSize">
				<template v-for="item in cells.headers">
					<vxe-column
						v-if="item.value === 'payDay'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default="{ row }">
							{{ row.payDay | dateFormat("YYYY-MM") }}
						</template>
					</vxe-column>

					<vxe-column v-else :field="item.value" :title="item.text" :key="item.value">
					</vxe-column>
				</template>
			</vxe-table>
		</div>
		<z-pagination
			:options="pageSizes"
			:total="pagination.total"
			@page="handlePage"
		></z-pagination>
	</div>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";
import moment from "moment";
import io from "socket.io-client";
import { localStorage, tools as lpTools } from "@/libs/util";

const isIntranet = lpTools.isIntranet();
const { baseURL, baseURLApi } = lpTools.baseURL();
const action = lpTools.isIntranet()
	? `${baseURL}report-management/pt/salary/upload`
	: `${baseURLApi}/report-management/pt/salary/upload`;

export default {
	name: "PtSalary",
	mixins: [TableMixins],

	data() {
		return {
			formId: "ptSalary",
			cells,
			dispatchList: "GET_SALARY_PT_LIST",
			dispatchDownLoad: "EXPORT_SALARY_PT_DOWNLOAD",
			socketPath: "global-send-message",
			action
		};
	},

	created() {
		this.getUserInfo();
		const userId = localStorage.get("user").ID;
		var room = "sendMessage";

		if (!userId) {
			this.toasted.error(result.msg);
		}

		// const socket = io(`${baseURL}global-sendMessage`, {
		// 	query: {
		// 		room,
		// 		userId
		// 	},
		// 	transports: ["websocket"]
		// });

		// socket.on(this.socketPath, result => {
		// 	this.toasted.dynamic(result.msg, result.code);
		// });
	},

	methods: {
		uploadedPtSalary({ result }) {
			this.toasted.dynamic(result.msg, result.code);
		},
		// 获取用户信息
		getUserInfo() {
			const user = localStorage.get("user");
			const token = localStorage.get("token");
			const secret = this.storage.get("secret");
			let code = "";
			if (secret) {
				code = lpTools.GetCode(secret);
			}
			this.user = user;

			this.fileHeaders = {
				"x-token": token,
				"x-user-id": user.id,
				"x-code": String(code)
			};
		},
		async onBack2() {
			const ym = moment(this.params.date).format("YYYYMM");
			await this.$store.dispatch("EXPORT_SALARY_PT_DOWNLOAD", ym);
		},
		handlePage({ pageSize, pageNum: pageIndex }) {
			this.params = {
				...this.params,
				pageSize,
				pageIndex
			};
		}
	}
};
</script>
