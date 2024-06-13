<template>
	<v-expansion-panel>
		<v-expansion-panel-header>路径配置</v-expansion-panel-header>
		<v-expansion-panel-content>
			<div class="table path-table">
				<vxe-table
					:data="desserts"
					:border="tableBorder"
					:loading="loading"
					:merge-cells="mergeCells"
					:size="tableSize"
					@checkbox-all="handleSelectAll"
					@checkbox-change="handleSelectChange"
				>
					<template v-for="item in cells.headers">
						<vxe-column
							v-if="item.value === 'currentPath'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							:width="item.width"
						>
							<template #default="{ row }">
								<div class="py-2 z-flex">
									<z-switch
										formId="actively"
										:formKey="row.formKey"
										:defaultValue="row.isCurrentPath"
										@change="onSwitch($event, row)"
									>
									</z-switch>
								</div>
							</template>
						</vxe-column>

						<vxe-column
							v-else
							:field="item.value"
							:title="item.text"
							:key="item.value"
							:width="item.width"
						>
						</vxe-column>
					</template>
				</vxe-table>
			</div>
		</v-expansion-panel-content>
	</v-expansion-panel>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";
import ConfigMixins from "@/views/main/pd/config/ConfigMixins";
import { R } from "vue-rocket";
import cells from "./cells";

const reverse = new Map([
	["insideDownload", "outsideDownload"],
	["outsideDownload", "insideDownload"],
	["insideUpload", "outsideUpload"],
	["outsideUpload", "insideUpload"]
]);

const brother = new Map([
	["insideDownload", "insideUpload"],
	["insideUpload", "insideDownload"],
	["outsideDownload", "outsideUpload"],
	["outsideUpload", "outsideDownload"]
]);

export default {
	name: "PathTable",
	mixins: [TableMixins, ConfigMixins],

	data() {
		return {
			dispatchList: "GET_CONFIG_PATH_PATH_LIST",
			cells,
			mergeCells: [
				{ row: 0, col: 0, rowspan: 2, colspan: 0 },
				{ row: 2, col: 0, rowspan: 2, colspan: 0 }
			]
		};
	},

	methods: {
		async onSwitch(value, row) {
			row.isCurrentPath = value;

			const { formKey, downloadType } = row;
			let [insideItem, outsideItem, items] = [{}, {}, []];

			const brotherItem = R.find(this.desserts, brother.get(formKey));
			const reverseItem = R.find(this.desserts, reverse.get(formKey));

			// 取反
			if (reverseItem) {
				reverseItem.isCurrentPath = !row.isCurrentPath;
			}

			// 测试(内部)
			if (downloadType === 1) {
				insideItem = {
					id: row.id,
					isDownload:
						formKey === "insideDownload"
							? row.isCurrentPath
							: brotherItem.isCurrentPath,
					isUpload:
						formKey === "insideUpload" ? row.isCurrentPath : brotherItem.isCurrentPath
				};

				if (reverseItem) {
					const reverseBrotherItem = R.find(
						this.desserts,
						brother.get(reverseItem.formKey)
					);

					outsideItem = {
						id: reverseItem.id,
						isDownload:
							reverseItem.formKey === "outsideDownload"
								? reverseItem.isCurrentPath
								: reverseBrotherItem.isCurrentPath,
						isUpload:
							reverseItem.formKey === "outsideUpload"
								? reverseItem.isCurrentPath
								: reverseBrotherItem.isCurrentPath
					};
				}
			}
			// 测试(客户)
			else {
				outsideItem = {
					id: row.id,
					isDownload:
						formKey === "outsideDownload"
							? row.isCurrentPath
							: brotherItem.isCurrentPath,
					isUpload:
						formKey === "outsideUpload" ? row.isCurrentPath : brotherItem.isCurrentPath
				};

				if (reverseItem) {
					const reverseBrotherItem = R.find(
						this.desserts,
						brother.get(reverseItem.formKey)
					);

					insideItem = {
						id: reverseItem.id,
						isDownload:
							reverseItem.formKey === "insideDownload"
								? reverseItem.isCurrentPath
								: reverseBrotherItem.isCurrentPath,
						isUpload:
							reverseItem.formKey === "insideUpload"
								? reverseItem.isCurrentPath
								: reverseBrotherItem.isCurrentPath
					};
				}
			}

			if (R.isYummy(outsideItem)) {
				items = items.concat([outsideItem]);
			}

			if (R.isYummy(insideItem)) {
				items = items.concat([insideItem]);
			}

			const result = await this.$store.dispatch("UPDATE_CONFIG_PATH_ITEM", items);

			this.toasted.dynamic(result.msg, result.code);
		}
	},

	watch: {
		"sabayon.data.list": {
			handler(list) {
				const [inside, outside] = [[], []];

				for (let item of list) {
					const template = {
						proId: item.proId,
						id: item.ID,
						env: item.downloadType === 1 ? "测试(内部)" : "测试(客户)",
						downloadType: item.downloadType
					};

					if (item.downloadType === 1) {
						inside.push(
							{
								...template,
								name: "下载路径",
								path: item.scan,
								formKey: "insideDownload",
								isCurrentPath: item.isDownload
							},
							{
								...template,
								name: "回传路径",
								path: item.upload,
								formKey: "insideUpload",
								isCurrentPath: item.isUpload
							}
						);
					} else {
						outside.push(
							{
								...template,
								name: "下载路径",
								path: item.scan,
								formKey: "outsideDownload",
								isCurrentPath: item.isDownload
							},
							{
								...template,
								name: "回传路径",
								path: item.upload,
								formKey: "outsideUpload",
								isCurrentPath: item.isUpload
							}
						);
					}
				}

				this.desserts = [...inside, ...outside];
			}
		}
	}
};
</script>

<style lang="scss">
.path-table {
	width: 100%;
}
</style>
