<template>
	<div class="month-report">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col
					v-for="(item, index) in cells.fields"
					:key="`entry_filters_${index}`"
					:cols="item.cols"
				>
					<template v-if="item.inputType === 'date'">
						<z-date-picker
							:formId="formId"
							:formKey="item.formKey"
							:hideDetails="item.hideDetails"
							:hint="item.hint"
							:label="item.label"
							:options="item.options"
							:range="item.range"
							:suffix="item.suffix"
							z-index="10"
							:defaultValue="item.defaultValue"
						></z-date-picker>
					</template>

					<template v-else>
						<z-select
							:formId="formId"
							:formKey="item.formKey"
							:clearable="item.clearable"
							:hideDetails="item.hideDetails"
							:hint="item.hint"
							:label="item.label"
							:options="item.options"
							:suffix="item.suffix"
							:defaultValue="item.defaultValue"
							@change="onChangeType($event, item)"
						></z-select>
					</template>
				</v-col>

				<div class="pb-6 btns">
					<z-btn
						v-for="item of cells.btns"
						:key="item.icon"
						:class="item.class"
						:color="item.color"
						outlined
					>
						{{ item.text }}
					</z-btn>
				</div>
			</v-row>
		</div>
	</div>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";

export default {
	name: "ProjectSettlement",
	mixins: [TableMixins],

	data() {
		return {
			formId: "projectSettlement",
			cells,
			desserts: []
		};
	},

	methods: {
		onChangeType(event, item) {}
	}
};
</script>
