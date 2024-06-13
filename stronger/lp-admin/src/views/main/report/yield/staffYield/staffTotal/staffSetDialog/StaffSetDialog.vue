<template>
	<div class="view-quality">
		<lp-dialog
			ref="dialog"
			title="设置"
			fullscreen
			persistent
			transition="dialog-bottom-transition"
			@dialog="handleDialog"
		>
			<div class="pt-6 table" slot="main">
				<div class="pb-6 btns">
					<z-btn class="pr-4" color="primary" @click="onNew">
						<v-icon>mdi-plus</v-icon>
						新增
					</z-btn>
				</div>

				<vxe-table :data="desserts" :border="tableBorder" :size="tableSize">
					<vxe-column type="seq" title="序号"></vxe-column>

					<template v-for="item in cells.headers">
						<vxe-column
							v-if="item.value === 'options'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							align="center"
							min-width="150px"
						>
							<template #default="{ row }">
								<z-btn
									class="pr-3"
									color="primary"
									outlined
									small
									@click="onEdit(row)"
								>
									编辑
								</z-btn>
								<z-btn
									class="pr-3"
									color="error"
									outlined
									small
									@click="onDelete(row)"
								>
									删除
								</z-btn>

								<z-btn
									class="pr-3"
									color="primary"
									outlined
									small
									@click="onLog(row)"
								>
									日志
								</z-btn>
							</template>
						</vxe-column>
						<vxe-column
							v-else-if="item.value === 'startTime'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							:width="item.width"
						>
							<template #default="{ row }">
								<div class="py-2 z-flex">
									{{ row[item.value] | dateFormat("YYYY-MM-DD") }}
								</div>
							</template>
						</vxe-column>
						<vxe-column
							v-else-if="item.value === 'UpdatedAt'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							:width="item.width"
						>
							<template #default="{ row }">
								<div class="py-2 z-flex">
									{{ row[item.value] | dateFormat("YYYY-MM-DD HH:mm:ss") }}
								</div>
							</template>
						</vxe-column>
						<vxe-column
							v-else
							:field="item.value"
							:title="item.text"
							:key="item.value"
							align="center"
						>
						</vxe-column>
					</template>
				</vxe-table>
			</div>
		</lp-dialog>

		<!-- 新增/编辑 BEGIN -->
		<!-- :config="{
        op0AsTheBlock: {
          mutex: [
            {
              formKey: 'op0AsTheInvoice',
              always: true
            }
          ]
        },

        op0AsTheInvoice: {
          mutex: [
            {
              formKey: 'op0AsTheBlock',
              always: true
            }
          ]
        }
      }" -->
		<z-dynamic-form
			ref="dynamic"
			:detail="detailInfo"
			:fieldList="cells.fields"
			@confirm="handleConfirm"
		></z-dynamic-form>
		<!-- 新增/编辑 END -->

		<log-dialog ref="log"></log-dialog>
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import TableMixins from "@/mixins/TableMixins";
import DialogMixins from "@/mixins/DialogMixins";
import cells from "./cells";

export default {
	name: "StaffSetDialog",
	mixins: [DialogMixins, TableMixins],

	data() {
		return {
			formId: "staffSetDialog",
			dispatchList: "GET_STAFF_YIELD_LIST",
			dispatchDelete: "DELETE_STAFF_YIELD",
			cells,
			manual: true,
			desserts: [],
			detailRow: {}
		};
	},

	methods: {
		onNew() {
			this.getDetail({});
			this.$refs.dynamic.open({
				title: "新增",
				status: -1
			});
		},

		onEdit(row) {
			this.getDetail(row);
			this.$refs.dynamic.open({
				title: "编辑",
				status: 1
			});
		},

		handleConfirm(effect, form) {
			form = {
				status: effect.status,
				id: this.detailInfo.ID,
				...form
			};

			this.updateListItem(form, "UPDATE_REPORT_YIELD_ITEM");
		},

		async onLog(row) {
			const params = {
				id: row.ID
			};

			const result = await this.$store.dispatch("GET_STAFF_YIELD_LOGS", params);

			this.toasted.dynamic(result.msg, result.code);

			if (result.code == 200) {
				this.$refs.log.desserts = result.data.list;
				this.$refs.log.onOpen();
			}
		},

		onDelete(row) {
			this.$modal({
				visible: true,
				title: "删除提示",
				content: "请确认是否要删除？",
				confirm: () => {
					this.deleteRows([row.ID]);
				}
			});
		},

		async deleteRows(ids) {
			const result = await this.$store.dispatch("DELETE_STAFF_YIELD", ids);

			this.toasted.dynamic(result.msg, result.code);

			if (result.code == 200) {
				this.getList();
			}
		}
	},

	computed: {
		...mapGetters(["auth", "pro"])
	},

	watch: {
		dialog: {
			handler(dialog) {
				if (dialog) {
					this.cells.fields[0].options = this.auth.proItems;
					this.getList();
				}
			},
			immediate: true
		}
	},

	components: {
		"log-dialog": () => import("./logDialog")
	}
};
</script>
