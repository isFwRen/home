<template>
	<div class="cost">
		<vxe-form :data="forms" title-width="115">
			<vxe-form-item
				:title="el.title"
				:field="el.field"
				v-for="(el, index) in cells.formList"
				:key="index"
				:span="6"
			>
				<template #default="{ data }">
					<vxe-input
						v-if="el.name === 'input'"
						v-model="data[el.field]"
						:disabled="el.disabled"
					></vxe-input>
				</template>
			</vxe-form-item>

			<vxe-form-item> </vxe-form-item>

			<vxe-form-item title="给付/拒付理由" field="address" :span="24">
				<template #default="{ data }">
					<vxe-textarea
						v-model="data.reason"
						placeholder="请输入地址"
						:autosize="{ minRows: 6, maxRows: 10 }"
						clearable
					></vxe-textarea>
				</template>
			</vxe-form-item>
		</vxe-form>
	</div>
</template>

<script>
import cells from "./cells";
export default {
	props: {
		options: {
			type: Object
		}
	},
	data() {
		return {
			forms: {},
			cells
		};
	},
	watch: {
		options: {
			handler: function (value) {
				this.forms = value;
			},
			immediate: true
		},
		forms: {
			handler: function (value) {
				this.$emit("updateForm", { forms: value, key: "quotaForms" });
			},
			deep: true
		}
	},
	methods: {}
};
</script>

<style lang="scss" scoped></style>
