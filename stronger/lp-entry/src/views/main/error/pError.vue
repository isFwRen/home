<template>
	<div class="lp-error">
		<div class="mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col :cols="2">
					<z-select :formId="searchFormId" formKey="proCode" clearable hideDetails label="项目"
						:options="auth.proItems"></z-select>
				</v-col>

				<v-col :cols="3">
					<z-date-picker :formId="searchFormId" formKey="date" hideDetails label="日期" range
						z-index="10"></z-date-picker>
				</v-col>

				<v-col :cols="2">
					<z-text-field :formId="searchFormId" formKey="name" hideDetails label="字段名称"> </z-text-field>
				</v-col>

				<!-- <v-col :cols="2">
					<z-select
						:formId="searchFormId"
						formKey="complaint"
						clearable
						hideDetails
						label="申诉"
						:options="cells.complaintOptions"
					></z-select>
				</v-col> -->

				<div class="z-flex">
					<z-btn class="pb-3 pl-3" color="primary" @click="onSearch">
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>
				</div>
			</v-row>
		</div>

		<!-- <z-btn class="mb-6" color="primary" :disabled="!ids.length" small outlined @click="onBatchAppeal"> 批量申诉 </z-btn> -->

		<div class="table error-detail-table">
			<vxe-table :data="desserts" :border="tableBorder" :max-height="tableMaxHeight" :size="tableSize"
				:stripe="tableStripe" @checkbox-all="handleSelectAll" @checkbox-change="handleSelectChange">
				<vxe-column type="checkbox" width="40"></vxe-column>

				<template v-for="item in cells.headers">
					<!-- 日期数据 BEGIN -->
					<vxe-column v-if="item.value === 'submitDay'" :field="item.value" :title="item.text" :key="item.value">
						<template #default="{ row }">
							<span>{{ row.submitDay.slice(0, 10) }}</span>
						</template>
					</vxe-column>
					<!-- 日期数据 END -->

					<!-- 错误数据 BEGIN -->
					<vxe-column v-else-if="item.value === 'wrong'" :field="item.value" :title="item.text" :key="item.value">
						<template #default="{ row }">
							<span v-if="_tools.compareString(row.wrong, row.right, 'error--text').targetHtml"
								v-html="_tools.compareString(row.wrong, row.right, 'error--text').targetHtml"></span>
							<span v-else>{{ row.wrong }}</span>
						</template>
					</vxe-column>
					<!-- 错误数据 END -->

					<!-- 正确数据 BEGIN -->
					<vxe-column v-else-if="item.value === 'right'" :field="item.value" :title="item.text" :key="item.value">
						<template #default="{ row }">
							<span v-if="_tools.compareString(row.right, row.wrong, 'error--text').targetHtml"
								v-html="_tools.compareString(row.right, row.wrong, 'error--text').targetHtml"></span>
							<span v-else>{{ row.right }}</span>
						</template>
					</vxe-column>
					<!-- 正确数据 END -->

					<!-- 解析 BEGIN -->
					<vxe-column v-else-if="item.value === 'analysis'" :field="item.value" :title="item.text" :key="item.value">
						<template>
							<span class="primary--text"> 规则解析 </span>
						</template>
					</vxe-column>
					<!-- 解析 END -->

					<!-- 申诉 BEGIN -->
					<vxe-column v-else-if="item.value === 'isComplain'" :field="item.value" :title="item.text" :key="item.value"
						:width="item.width">
						<template #default="{ row, rowIndex }">
							<z-switch :formId="formId" :formKey="`appeal_${rowIndex}_${row.id}`" :disabled="row.isComplain"
								:label="row.isComplain ? '已申诉' : '申诉'" :defaultValue="row.isComplain" @change="onAppeal($event, row)">
							</z-switch>
						</template>
					</vxe-column>
					<!-- 申诉 END -->

					<!-- 差错审核 BEGIN -->
					<vxe-column v-else-if="item.value === 'isWrongConfirm'" :field="item.value" :title="item.text"
						:key="item.value" :width="item.width">
						<template #default="{ row }">
							<span>{{ row.isComplain ? row.isOperationLog : row.isAudit ? row.isOperationLog : "" }}</span>
						</template>
					</vxe-column>
					<!-- 差错审核 END -->

					<vxe-column v-else-if="item.value === 'fieldName'" :field="item.value" :title="item.text" :key="item.value">
						<template #default="{ row }">
							<span @click="onShow(row)">{{ row[item.value] }}</span>
						</template>
					</vxe-column>

					<vxe-column v-else :field="item.value" :title="item.text" :key="item.value"></vxe-column>
				</template>
			</vxe-table>

			<z-pagination class="mt-4" :total="pagination.total" :options="pageSizes" @page="handlePage"></z-pagination>
		</div>

		<lp-dialog ref="showImg" :fullscreen="true">
			<div slot="main">
				<div class="img-wrapper">
					<lp-images class="preview-img" :src="imgUrl" />
				</div>
				<!-- <img :src="imgUrl" alt="" /> -->
			</div>
		</lp-dialog>
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import { tools as lpTools } from "@/libs/util";
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";
import axios from "axios";
import { localStorage } from "vue-rocket";

const { baseURLApi } = lpTools.baseURL();

export default {
	name: "pError",
	mixins: [TableMixins],

	data() {
		return {
			formId: "pError",
			cells,
			dispatchList: "PRACTICE_ERROR_GET_LIST",

			manual: true,
			imgUrl: "",

			instance: ""
			// list: []
		};
	},

	created() {
		const token = localStorage.get("token");
		const user = localStorage.get("user");

		this.instance = axios.create({
			headers: {
				"x-token": token,
				"x-user-id": user.id
			}
		});
	},

	computed: {
		...mapGetters(["auth"])
	},

	methods: {
		// 批量申诉
		async onBatchAppeal() {
			const list = [];

			this.ids.map(id => {
				list.push({ id });
			});

			const body = {
				proCode: this.forms[this.searchFormId].proCode,
				complainConfirm: true,
				list
			};

			const result = await this.$store.dispatch("ERROR_APPEAL_ITEMS", body);

			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				this.getList();
			}
		},

		async onShow(row) {
			let reg = new RegExp("/files/files/", "g");
			this.imgUrl = baseURLApi + "files/" + row.path + row.picture;
			this.imgUrl = this.imgUrl.replace(reg, "/files/");

			let item = await this.transform(this.imgUrl);

			this.getReader(item).then(res => {
				this.imgUrl = res;
			});

			this.$refs.showImg.onOpen();
		},

		// 申诉
		async onAppeal(value, row) {
			const body = {
				proCode: this.forms[this.searchFormId].proCode,
				complainConfirm: value,
				list: [{ id: row.id || row.ID }]
			};

			const result = await this.$store.dispatch("ERROR_APPEAL_ITEMS", body);

			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				this.getList();
			}
		},

		// 图片转Base64格式
		async transform(el) {
			let code = "";
			const secret = localStorage.get("secret") || "";
			if (secret) {
				code = lpTools.GetCode(secret);
			}
			let res = await this.instance.get(el, {
				responseType: "blob",
				headers: {
					"x-code": String(code)
				}
			});
			return res.data;
		},

		// 读取图片文件
		getReader(blob) {
			return new Promise((resolve, reject) => {
				const reader = new FileReader();
				reader.onloadend = () => {
					const base64String = reader.result;
					resolve(base64String);
				};
				reader.onerror = reject;
				reader.readAsDataURL(blob);
			});
		}
	},
	components: {
		"lp-images": () => import("@/components/lp-images")
	}
	// watch: {
	// 	desserts:{
	// 		handler(newValue){
	// 			this.list = newValue.filter(el => el.name == localStorage.get("user").name);
	// 		},
	// 	}
	// },
};
</script>

<style scoped>
.col-2 {
	flex: 0 0 13.6666666667%;
	max-width: 13.666667%;
}

.col-3 {
	flex: 0 0 21.6666666667%;
	max-width: 21.666667%;
}

.img-wrapper {
	width: 80vw;
	height: 88vh;
	margin: 0 auto;
}

.preview-img {
	width: 100%;
	height: 80%;
}
</style>
