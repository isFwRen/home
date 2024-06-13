<template>
	<div class="lp-case-fields">
		<v-row class="z-flex">
			<template v-for="(item, index) in fieldList">
				<v-col
					:key="`${formId}${item.formKey}${index}`"
					class="z-flex align-center"
					:cols="cols"
				>
					<label class="mr-3">{{ item.label }}</label>

					<z-text-field
						v-if="item.inputType === 'text'"
						:formId="formId"
						:formKey="item.formKey"
						:class="item.class"
						:readonly="readonly"
						width="250"
						:defaultValue="item.defaultValue"
					>
					</z-text-field>

					<z-select
						v-if="item.inputType === 'select'"
						:formId="formId"
						:formKey="item.formKey"
						:class="item.class"
						:readonly="readonly"
						width="250"
						:options="[{ label: item.defaultValue, value: item.defaultValue }]"
						:defaultValue="item.defaultValue"
					></z-select>

					<div v-if="item.inputType === 'const'" :class="['z-flex', item.class]">
						<z-select
							:formId="formId"
							:formKey="'waiting'"
							:readonly="readonly"
							width="50"
							:options="item.options"
							:defaultValue="item.defaultValue"
						></z-select>

						<z-text-field
							:formId="formId"
							:formKey="item.formKey"
							:readonly="readonly"
							width="250"
							:defaultValue="item.defaultValue"
						>
						</z-text-field>
					</div>
				</v-col>
			</template>
		</v-row>
	</div>
</template>

<script>
export default {
	name: "LPCaseFields",

	props: {
		cols: {
			type: [Number, String],
			default: 12
		},

		fieldList: {
			type: Array,
			required: true
		},

		readonly: {
			type: Boolean,
			default: true
		}
	},

	data() {
		return {
			formId: "LPCaseFields"
		};
	}
};
</script>

<style scoped lang="scss">
.lp-case-fields {
	label {
		display: block;
		min-width: 57px;
		text-align: right;
	}
}
</style>
