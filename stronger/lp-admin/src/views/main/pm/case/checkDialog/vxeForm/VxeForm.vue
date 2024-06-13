<template>
	<vxe-form :data="formData" :rules="formRules" title-align="right" title-width="120">
		<vxe-form-item
			v-for="(el, index) in forms"
			:field="el.field"
			:title="el.title"
			:span="el.span"
			:key="index"
		>
			<template #default="{ data }">
				<vxe-input
					v-if="el.name == 'input'"
					v-model="data[el.field]"
					:placeholder="el.placeholder"
				></vxe-input>

				<vxe-textarea
					v-else-if="el.name == 'text'"
					v-model="data[el.field]"
					:placeholder="el.placeholder"
					:autosize="{ minRows: 2, maxRows: 4 }"
				></vxe-textarea>

				<vxe-select
					v-else-if="el.name == 'select'"
					v-model="data[el.field]"
					transfer
					:placeholder="el.placeholder"
				>
					<vxe-option
						v-for="item in el.options"
						:key="item.value"
						:value="item.value"
						:label="item.label"
					></vxe-option>
				</vxe-select>

				<vxe-input
					v-else-if="el.name == 'date'"
					v-model="data[el.field]"
					type="date"
					:placeholder="el.placeholder"
					transfer
				></vxe-input>

				<vxe-checkbox-group v-else-if="el.name == 'checkbox'" v-model="data[el.field]">
					<vxe-checkbox
						v-for="item in el.rect"
						:label="item.label"
						:content="item.content"
						:key="item.label"
					></vxe-checkbox>
				</vxe-checkbox-group>

				<vxe-radio-group v-else-if="el.name == 'radio'" v-model="data[el.field]">
					<vxe-radio
						v-for="item in el.circle"
						:label="item.label"
						:content="item.content"
						:key="item.label"
					></vxe-radio>
				</vxe-radio-group>
			</template>
		</vxe-form-item>
	</vxe-form>
</template>

<script>
import { sessionStorage } from "vue-rocket";
export default {
	props: {
		belongModule: {
			type: String,
			default: () => {}
		},
		belongModuleForm: {
			type: String,
			default: () => {}
		},
		forms: {
			type: Array,
			default: () => {}
		},
		formRules: {
			type: Object,
			default: () => {}
		},
		cloudData: {
			type: Object,
			default: () => {}
		}
	},
	data() {
		return {
			formData: {
				// caseNum: null,
				// applyTime: null,
				// manageAgency: null,
				// compensationType: [],
				// caseType: null
			}
		};
	},
	methods: {},

	watch: {
		cloudData: {
			handler(newValue) {
				// console.log("cloudData", newValue);
				this.formData = { ...newValue };
			},
			deep: true,
			immediate: true
		},
		formData: {
			handler(newValue) {
				// console.log("formData", newValue);
				let content = sessionStorage.get("checkForm");
				content[this.belongModule][this.belongModuleForm] = this.formData;
				sessionStorage.set("checkForm", content);
			},
			deep: true
		}
	}
};
</script>
