<template>
	<div class="expense-template">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col :cols="2">
					<z-select
						:formId="searchFormId"
						formKey="proCode"
						hideDetails
						label="项目编码"
						:options="auth.proItems"
					></z-select>
				</v-col>

				<v-col :cols="2">
					<z-text-field
						:formId="searchFormId"
						formKey="name"
						hideDetails
						label="影像名称"
					></z-text-field>
				</v-col>

				<div class="z-flex pb-2">
					<z-btn class="mr-3" color="primary" @click="getList">
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>

					<z-btn class="mr-3" color="primary" @click="onAddItem">
						<v-icon>mdi-plus</v-icon>
						新增
					</z-btn>

					<z-btn class="mr-3" color="primary" @click="onChoice">
						<v-icon size="19" class="mr-1">mdi-selection-ellipse-arrow-inside</v-icon>
						{{ showSelect ? "取消" : "选择" }}
					</z-btn>

					<z-btn
						v-if="choiceItems.length"
						class="mr-3"
						color="error"
						@click="onBatchDelete"
					>
						<v-icon>mdi-trash-can-outline</v-icon>
						批量删除
					</z-btn>

					<z-select
						:formId="formId"
						formKey="code"
						defaultValue="large"
						dense
						hideDetails
						:options="cells.moreOptions"
						outlined
						width="125"
						@change="changeItem"
					></z-select>
				</div>
			</v-row>
		</div>

		<detail-table
			v-if="mode === 'detail'"
			ref="detailTable"
			:desserts="desserts"
			:filters="filters"
		></detail-table>

		<image-list
			v-else
			ref="imageList"
			:desserts="desserts"
			:size="mode"
			:showSelect="showSelect"
			@select="onSelectImages"
			@deleted="getList"
			@renamed="getList"
		></image-list>

		<z-dynamic-form
			ref="dynamic"
			:formId="formId"
			:title="title"
			:fieldList="cells.fields"
			:config="config"
			:width="600"
			@change:file="handleFileUpload"
			@confirm="handleConfirm"
		>
		</z-dynamic-form>
	</div>
</template>

<script>
import { mapState, mapGetters } from "vuex";
import { rocket } from "vue-rocket";
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";

export default {
	name: "TeachingExpenseTemplate",
	mixins: [TableMixins],

	data() {
		return {
			formId: "TeachingExpenseTemplate",
			cells,
			dispatchList: "GET_PM_TEACHING_EXPENSE_TEMPLATE_LIST",
			mode: "large",
			showSelect: false,
			choiceItems: [],
			manual: true,
			title: "新增",
			config: {},
			files: [],
			filters: {}
		};
	},

	created() {
		this.config = {
			proCode: {
				items: this.auth.proItems
			}
		};
	},

	methods: {
		changeCode(value) {
			this.proCode = value;
		},

		onAddItem() {
			this.$refs.dynamic.open();
		},

		changeItem(value) {
			this.filters = this.forms[this.searchFormId];

			this.mode = value;
		},

		onChoice() {
			this.showSelect = !this.showSelect;
		},

		onSelectImages(images) {
			this.choiceItems = images;
		},

		async onBatchDelete() {
			const form = this.forms[this.searchFormId];

			const body = {
				proCode: form.proCode,
				ids: this.choiceItems
			};

			const result = await this.$store.dispatch(
				"DELETE_PM_TEACHING_EXPENSE_TEMPLATE_ITEM",
				body
			);

			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				this.getList();
			}
		},

		handleFileUpload(files) {
			const maxSize = 2 * 1024;

			for (let file of files) {
				const size = file.size / 1024;

				if (size > maxSize) {
					this.toasted.warning(`${file.name}超过2M!`);
					return;
				}
			}

			this.files = files;
		},

		async handleConfirm({}, form) {
			const body = {
				...form,
				files: this.files,
			};

			const result = await this.$store.dispatch(
				"POST_PM_TEACHING_EXPENSE_TEMPLATE_ITEM",
				body
			);

			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				rocket.emit("ZHT_RESET_FORM", this.formId);
				this.$refs.dynamic.close();
				this.getList();
			}
		}
	},

	computed: {
		...mapState(["forms"]),
		...mapGetters(["auth"])
	},

	components: {
		"detail-table": () => import("./detailTable"),
		"image-list": () => import("./imageList")
	}
};
</script>
