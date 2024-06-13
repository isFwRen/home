<template>
	<div class="op-text-field" :id="id">
		<div class="z-flex align-center op-text-field-label">
			<!-- 字段名称 BEGIN -->
			<v-tooltip bottom>
				<template v-slot:activator="{ on, attrs }">
					<label
						:class="{
							'label-accent': accent,
							'text-red': computedTempField
						}"
						v-bind="attrs"
						v-on="on"
						>{{ label }}
					</label>
				</template>
				<span>{{ labelTip }}</span>
			</v-tooltip>
			<!-- 字段名称 END -->

			<!-- 字段规则 BEGIN -->
			<v-icon class="ml-1" @click="onRuleIconClick">{{ appendIcon }} </v-icon>
			<!-- 字段规则 END -->

			<span
				v-if="showLength && value.length"
				:style="{
					color: value.length < 8 ? '#ff7f50' : value.length < 15 ? '#ff4500' : 'red',
					fontSize: value.length < 8 ? '15px' : value.length < 15 ? '18px' : '20px'
				}"
				>{{ value.length }}</span
			>
		</div>

		<v-text-field
			:id="`input_${id}`"
			:ref="`input_${id}`"
			v-model="value"
			autocomplete="off"
			:class="disabled ? 'disabled-field' : ''"
			:autofocus="autofocus"
			:disabled="disabled"
			:error="error"
			:error-messages="errorMessages"
			:placeholder="computedIncludes"
			:rules="rules"
			@blur="onBlur"
			@change="onChange"
			@keydown="onKeydown"
			@keyup.enter="onEnter"
			@focus="onFocus"
			@input="onInput"
			@keyup="onKeyup"
			@click="onClear"
		></v-text-field>

		<div class="slot">
			<slot></slot>
		</div>

		<!-- 动态信息提示 BEGIN -->
		<div v-show="hint" v-html="hint"></div>
		<!-- 动态信息提示 END -->

		<!-- 静态信息提示 BEGIN -->
		<div v-for="(message, index) of hints" :key="index" v-html="message"></div>
		<!-- 静态信息提示 END -->

		<!-- 下拉 BEGIN -->
		<div v-if="typing && dropdownItems.length" class="dropdown">
			<z-list
				:dataSource="dropdownItems"
				:max-height="304"
				:defaultValue="value"
				@hover="handleHoverListItem"
				@select="handleSelectListItem"
			>
				<template v-slot:default="{ item }">
					<div class="list-item">
						{{ item }}
					</div>
				</template>
			</z-list>
		</div>
		<!-- 下拉 END -->
	</div>
</template>

<script>
import OpqMixins from "./OpqMixins";
import ValidationsMixins from "./ValidationsMixins";
import DropdownMixins from "./DropdownMixins";

const debounce = (() => {
	let timer = null;

	return (fn, interval = 60) => {
		timer = setTimeout(() => {
			clearTimeout(timer);
			fn();
		}, interval);
	};
})();

export default {
	name: "OpTextField",
	mixins: [OpqMixins, ValidationsMixins, DropdownMixins],

	props: {
		accent: {
			type: Boolean,
			default: false
		},

		autofocus: {
			type: Boolean,
			default: false
		},

		defaultValue: {
			type: [String, Number, Boolean, Array, Object],
			default: undefined
		},

		disabled: {
			type: Boolean,
			default: false
		},

		field: {
			type: Object,
			default: () => ({})
		},

		fieldsIndex: {
			type: [Number, String],
			default: 0
		},

		fieldsList: {
			type: Array,
			default: () => []
		},

		fieldsObject: {
			type: Object,
			default: () => ({})
		},

		focusFieldsIndex: {
			type: Number,
			default: -1
		},

		hint: {
			type: [String, Number],
			required: false
		},

		id: {
			type: String,
			required: false
		},

		label: {
			type: String,
			required: false
		},

		labelTip: {
			type: String,
			required: false
		},

		op: {
			type: String,
			required: false
		},

		thumbIndex: {
			type: Number,
			default: 0
		},

		appendIcon: {
			type: String,
			default: "mdi-help-circle-outline"
		},

		appendTip: {
			type: String,
			required: false
		},

		validations: {
			type: Array,
			required: false
		}
	},

	data() {
		return {
			value: null,

			showLength: false,
			allowSpace: false,

			focus: false,

			comRules: []
		};
	},

	computed: {
		computedTempField() {
			if (this.accent && this.fieldsIndex === this.focusFieldsIndex) {
				return true;
			}

			return false;
		},

		computedIncludes() {
			return `${this.includes && this.includes.toString()}`;
		}
	},

	watch: {
		defaultValue: {
			handler(value) {
				this.value = value;
			},
			immediate: true
		},
		field: {
			handler(newVal) {
				if (this.op == "op0") {
					// if (newVal.name == '图片页码') return
					this.value = newVal.op0Value;
				}
				if (this.op == "op1") {
					this.value = newVal.op1Value;
				}
				if (this.op == "op2") {
					this.value = newVal.op2Value;
				}
				if (this.op == "opq") {
					this.value = newVal.opqValue;
				}
			},
			immediate: true,
			deep: true
		}
	},

	methods: {
		onBlur(event) {
			this.focus = false;
			event.customValue = this.value;
			this.$emit("blur", event);
		},

		onChange(value) {
			this.$emit("change", value);
		},

		onKeydown(event) {
			if (event.keyCode === 38 && this.typing && this.items.length) {
				return;
			}

			debounce(() => {
				event.customValue = this.value;
				this.$emit("keydown", event);
			});
		},

		onEnter(event) {
			event.customValue = this.value;
			event.customItems = this.items;

			this.$emit("enter", event);

			this.validateAfterEnterInput();
		},

		onFocus(event) {
			this.focus = true;
			event.customValue = this.value;
			this.comRules = [];
			for (let rule of this.validations) {
				this.comRules.push(rule["key"]);
			}

			this.$emit("focus", event);
		},

		onClear() {
			this.$emit("focusClear", -1);
		},

		onInput(value) {
			this.typing = true;

			this.$emit("input", value);

			this.validateAfterEnterInput();
		},

		onKeyup(event) {
			if (event.keyCode === 38 && this.typing && this.items.length) {
				return;
			}
			if (!this.comRules.includes("spaced")) {
				this.value = this.allowSpace ? this.value : this.value.replace(/\s+/g, "");
			}
			if (this.comRules.includes("show_length")) {
				this.showLength = true;
			}
			// this.value = this.value.trim();
			event.customValue = this.value;
			this.$emit("keyup", event);
		},

		onRuleIconClick() {
			this.$emit("ruleClick", { name: this.field.name });
		}
	}
};
</script>

<style scoped lang="scss">
.op-text-field {
	position: relative;
	min-height: 94px;

	label {
		color: #000;
	}

	.op-text-field-label {
		position: relative;
		top: 18px;
		z-index: 0;

		label.label-accent {
			color: #0610f7;
			font-size: 15px;
			font-weight: bolder;
		}
	}

	.disabled-field::v-deep {
		cursor: not-allowed;

		.v-text-field__slot {
			background-color: #d9d9d9;
		}
	}
}

.dropdown {
	position: absolute;
	top: 90px;
	left: 0;
	padding: 4px 0;
	width: 100%;
	max-height: 304px;
	background-color: #fff;
	box-shadow: 0 4px 6px 0 rgb(32 33 36 / 28%);
	z-index: 1;

	.list-item {
		padding: 4px 8px;
		color: rgba(0, 0, 0, 0.87);
		cursor: pointer;
	}
}
</style>