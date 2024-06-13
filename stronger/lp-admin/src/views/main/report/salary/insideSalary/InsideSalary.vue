<template>
	<div class="inside-salary">
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
							:range="item.range"
							:suffix="item.suffix"
							z-index="10"
							:defaultValue="item.defaultValue"
						></z-date-picker>
					</template>
				</v-col>

				<div class="z-flex">
					<z-btn class="pb-3 pl-3" color="primary" @click="onSearch">
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
				@response="handleResponse"
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
				@click="onBack1(item)"
			>
				{{ item.text }}
			</z-btn>
		</div>

		<div class="table inside-salary-table">
			<vxe-table :border="tableBorder" :data="desserts" :size="tableSize">
				<vxe-column type="seq" title="序号"></vxe-column>
				<template v-for="item in headers">
					<vxe-colgroup
						v-if="item.children != undefined && !keyList[item.value]"
						align="center"
						:title="item.text"
						:key="item.value"
					>
						<template v-for="record in item.children">
							<vxe-column
								:field="record.value"
								:title="record.text"
								:key="record.value"
								align="center"
								:width="record.width ? record.width : '70px'"
							></vxe-column>
						</template>
					</vxe-colgroup>

					<vxe-colgroup
						v-else-if="keyList[item.value]"
						align="center"
						:title="item.text"
						:key="item.value"
					>
						<template v-for="code in item.children">
							<vxe-column
								#default="{ row }"
								:title="code.text"
								:key="code.value"
								align="center"
								:width="code.width ? code.width : '70px'"
							>
								{{ row[item.value][code.value] }}
							</vxe-column>
						</template>
					</vxe-colgroup>

					<vxe-column
						v-else
						:field="item.value"
						:title="item.text"
						:key="item.value"
						align="center"
						:width="item.width ? item.width : '50px'"
					></vxe-column>
				</template>
			</vxe-table>

			<z-pagination
				:options="pageSizes"
				:total="pagination.total"
				@page="handlePage"
			></z-pagination>
		</div>
	</div>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";
import moment from "moment";
import { localStorage, tools } from "vue-rocket";
import io from "socket.io-client";
import { tools as lpTools } from "@/libs/util";

const { baseURL, baseURLApi } = lpTools.baseURL();

const action = lpTools.isIntranet()
	? `${baseURL}report-management/pt/salary/upload`
	: `${baseURLApi}report-management/internal/salary/upload`;

export default {
	name: "InsideSalary",
	mixins: [TableMixins],

	data() {
		return {
			formId: "InsideSalary",
			cells,
			dispatchList: "GET_SALARY_INSIDE_LIST",
			dispatchUpload: "IMPORT_SALARY_INSIDE_UPLOAD",
			dispatchDownLoad: "EXPORT_SALARY_INSIDE_DOWNLOAD",
			manual: true,
			ym: "",
			headers: [],
			socketPath: "global-send-message",
			action,
			//需要动态处理的key值 散列表的形式便于查找
			keyList: {
				workLoad: 1,
				reducedProportion: 1,
				externalQuality: 1,
				internalQuality: 1,
				totalBusinessVolume: 1
			}
		};
	},

	created() {
		this.getUserInfo();
		const userId = localStorage.get("user").ID;
		var room = "sendMessage";

		if (!userId) {
			this.toasted.dynamic(result.msg, 400);
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

	watch: {
		"sabayon.data.list": {
			handler(list) {
				if (list) {
					this.buildHeaders(list);
				}
			},
			immediate: true,
			deep: true
		}
	},
	methods: {
		handleResponse({ result }) {
			this.toasted.dynamic(result.msg, result.code);
		},
		// 获取用户信息
		getUserInfo() {
			const user = this.localStorage.get("user");
			const token = this.localStorage.get("token");
			const secret = this.storage.get("secret");
			const code = lpTools.GetCode(secret);

			this.user = user;

			this.fileHeaders = {
				"x-token": token,
				"x-user-id": user.id,
				"x-code": String(code)
			};

		},

		async onBack1() {
			const ym = moment(this.params.date).format("YYYYMM");
			await this.$store.dispatch("EXPORT_SALARY_INSIDE_DOWNLOAD", ym);
		},

		handlePage({ pageSize, pageNum: pageIndex }) {
			this.params = {
				...this.params,
				pageSize,
				pageIndex
			};
		},

		buildHeaders(list) {
			if (!list.length) {
				return;
			}
			//获取表的标题栏目
			this.headers = tools.deepClone(this.cells.headers);
			//添加动态key
			for (let i in this.headers) {
				if (this.keyList[this.headers[i].value] === 1) {
					//遍历属性转换key
					for (let key in list[0][this.headers[i].value]) {
						//设置标题和对应字段
						this.headers[i].children.push({
							text: key,
							value: key
						});
					}
				}
			}
		}
	}
};
</script>
