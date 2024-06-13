<template>
	<div class="update-menu-dialog">
		<lp-dialog
			ref="dialog"
			:title="title"
			:width="650"
			transition="dialog-bottom-transition"
			@dialog="handleDialog"
		>
			<div slot="main">
				<z-btn-toggle
					:formId="formId"
					formKey="menuType"
					color="primary"
					mandatory
					:options="cells.menuTypeOptions"
					:defaultValue="selectedMenuType"
					@change="changeMenuType"
				></z-btn-toggle>

				<v-row class="z-flex align-end">
					<v-col
						v-for="(item, index) in fields"
						:key="`case_filters_${index}`"
						:class="item.colsClass"
						:cols="item.cols"
					>
						<template v-if="item.inputType === 'text' && item.show">
							<z-text-field
								:formId="formId"
								:formKey="item.formKey"
								:disabled="item.disabled"
								:hideDetails="item.hideDetails"
								:hint="item.hint"
								:label="item.label"
								:suffix="item.suffix"
								:prependOuterClass="item.prependOuterClass"
								:validation="item.validation"
								:defaultValue="detail[item.formKey]"
							>
								<span :class="item.prependOuterClass" slot="prepend-outer">{{
									item.prependOuter
								}}</span>
							</z-text-field>
						</template>

						<template v-else-if="item.inputType === 'select' && item.show">
							<z-select
								:formId="formId"
								:formKey="item.formKey"
								:disabled="item.disabled"
								:hideDetails="item.hideDetails"
								:hint="item.hint"
								:label="item.label"
								:clearable="item.clearable"
								:options="item.options"
								:suffix="item.suffix"
								:prependOuterClass="item.prependOuterClass"
								:validation="item.validation"
								:defaultValue="detail[item.formKey]"
							>
								<span :class="item.prependOuterClass" slot="prepend-outer">{{
									item.prependOuter
								}}</span>
							</z-select>
						</template>

						<template v-else-if="item.inputType === 'autocomplete' && item.show">
							<z-autocomplete
								:formId="formId"
								:formKey="item.formKey"
								:disabled="item.disabled"
								:hideDetails="item.hideDetails"
								:hint="item.hint"
								:label="item.label"
								:clearable="item.clearable"
								:options="item.options"
								:suffix="item.suffix"
								:prependOuterClass="item.prependOuterClass"
								:validation="item.validation"
								:defaultValue="detail[item.formKey]"
								@change="changeAutocomplete($event, item)"
							>
								<span :class="item.prependOuterClass" slot="prepend-outer">{{
									item.prependOuter
								}}</span>
							</z-autocomplete>
						</template>

						<template v-else-if="item.inputType === 'switch' && item.show">
							<z-switch
								:formId="formId"
								:formKey="item.formKey"
								:disabled="item.disabled"
								:label="item.label"
								:defaultValue="detail[item.formKey] || item.defaultValue"
							></z-switch>
						</template>
					</v-col>
				</v-row>
			</div>

			<div class="z-flex" slot="actions">
				<z-btn class="mr-3" color="normal" @click="onClose">取消</z-btn>

				<z-btn
					btnType="validate"
					:formId="formId"
					class="mr-2"
					color="primary"
					@click="onConfirm"
					>确认</z-btn
				>
			</div>
		</lp-dialog>
	</div>
</template>

<script>
import { rocket, R } from "vue-rocket";
import { mapGetters } from "vuex";
import DialogMixins from "@/mixins/DialogMixins";
import cells, { fields } from "./cells";

const type_menu_hidden = ["action", "api", "apiId"];
const type_menu_text = ["title"];
const type_btn_hidden = ["component", "path"];
const type_btn_autocomplete = ["title"];

export default {
	name: "UpdateMenuDialog",
	mixins: [DialogMixins],

	data() {
		return {
			formId: "UpdateMenuDialog",
			dispatchForm: "UPDATE_STAFF_MENU_MANAGE_TREE_LEAF",
			cells,
			fields,

			selectedMenuType: cells.MENU_TYPE
		};
	},

	methods: {
		// 菜单类型
		changeMenuType(value) {
			switch (value) {
				case cells.MENU_TYPE:
					for (let item of this.fields) {
						item.show = type_menu_hidden.includes(item.formKey) ? false : true;

						if (type_menu_text.includes(item.formKey)) {
							item.inputType = "text";
						}
					}
					break;

				case cells.BUTTON_TYPE:
					for (let item of this.fields) {
						item.show = type_btn_hidden.includes(item.formKey) ? false : true;

						if (type_btn_autocomplete.includes(item.formKey)) {
							item.inputType = "autocomplete";
						}
					}
					break;
			}

			this.selectedMenuType = value;
		},

		// 标题
		changeAutocomplete(value) {
			const result = R.find(this.staff.titleOptions, value);

			if (R.isYummy(result)) {
				this.detail = {
					...this.detail,
					action: result.action,
					api: result.path,
					apiId: result.id
				};
			}
		},

		// 确认
		async onConfirm() {
			const data = {
				status: this.status,
				...this.detail,
				...this.forms[this.formId]
			};

			if (data.menuType === "1") {
				if (this.storage.get("Leslie")) {
					this.submit(data);
				} else {
					this.toasted.warning("小样~ 你无权限使用菜单！");
				}
			} else {
				this.submit(data);
			}
		}
	},

	computed: {
		...mapGetters(["staff"])
	},

	watch: {
		dialog(dialog) {
			if (dialog) {
				this.$nextTick(() => {
					this.fields[1].options = this.staff.titleOptions;
					console.log(this.detail);
					console.log(this.fields);
				});
			}
		}
	}
};
</script>
