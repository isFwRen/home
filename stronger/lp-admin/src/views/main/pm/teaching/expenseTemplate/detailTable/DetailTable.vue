<template>
	<div class="detail-table">
		<vxe-table :data="desserts" border size="mini" @cell-click="onCell">
			<template v-for="item in cells.headers">
				<vxe-column
					v-if="item.value === 'size'"
					:field="item.value"
					:title="item.text"
					:key="item.value"
					:width="item.width"
				>
					<template #default="{ row }">
						<span> {{ row.size }} KB</span>
					</template>
				</vxe-column>

				<vxe-column
					v-else-if="item.value === 'name'"
					:field="item.value"
					:title="item.text"
					:key="item.value"
					:width="item.width"
				>
					<template #default="{ row }">
						<span style="cursor: pointer" @click="onImage(row)"> {{ row.name }}</span>
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
</template>

<script>
import { sessionStorage } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";
import cells from "./cells";

const { baseURLApi } = lpTools.baseURL();

export default {
	name: "TeachingExpenseTemplateDetailTable",

	props: {
		desserts: {
			type: Array,
			default: () => []
		},

		filters: {
			type: Object,
			default: () => {}
		}
	},

	data() {
		return {
			cells,
			cellInfo: {}
		};
	},

	methods: {
		onCell(cell) {
			this.cellInfo = {
				columnIndex: cell.$columnIndex,
				rowIndex: cell.$rowIndex,
				row: this.desserts[cell.$rowIndex]
			};
		},

		onImage() {
			const thumbs = [];

			this.desserts.map(dessert => {
				thumbs.push({
					thumbPath: `${baseURLApi}${dessert.path}`,
					path: `${baseURLApi}${dessert.path}`
				});
			});

			sessionStorage.set("thumbs", thumbs);

			window.open(
				`${location.origin}/normal/view-images`,
				"_blank",
				"toolbar=yes, scrollbars=yes, resizable=yes"
			);
		}
	}
};
</script>
