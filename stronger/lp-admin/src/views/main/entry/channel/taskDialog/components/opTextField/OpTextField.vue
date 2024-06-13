<template>
	<div class="op-text-field" :id="id">
		<div class="z-flex align-center op-text-field-label">
			<v-tooltip bottom>
				<template v-slot:activator="{ on, attrs }">
					<label :class="{ 'label-accent': accent }" v-bind="attrs" v-on="on"
						>{{ label }}
					</label>
				</template>
				<span>{{ labelTip }}</span>
			</v-tooltip>

			<v-tooltip right>
				<template v-slot:activator="{ on, attrs }">
					<v-icon class="ml-1" v-bind="attrs" v-on="on">{{ appendIcon }} </v-icon>
				</template>
				<span>{{ appendTip }}</span>
			</v-tooltip>

			<span v-show="showLength && value.length" class="pl-1 fw-bold primary--text">{{
				value.length
			}}</span>
		</div>

		<v-text-field
			:ref="`input_${id}`"
			v-model="value"
			autocomplete="off"
			:class="disabled ? 'disabled-field' : ''"
			:autofocus="autofocus"
			:disabled="disabled"
			:error="error"
			:error-messages="errorMessages"
			:rules="rules"
			@blur="onBlur"
			@change="onChange"
			@keydown="onKeydown"
			@keyup.enter="onEnter"
			@focus="onFocus"
			@input="onInput"
			@keyup="onKeyup"
		></v-text-field>

		<div class="slot">
			<slot></slot>
		</div>

		<div v-for="(message, index) of svHintList" :key="index" v-html="message"></div>

		<!-- 只能为 BEGIN -->
		<div v-if="includes.length">
			{{ includes }}
		</div>
		<!-- 只能为 END -->

		<!-- 下拉 BEGIN -->
		<div v-if="typing && items.length" class="dropdown">
			<z-list
				:dataSource="dropdownItems"
				:max-height="304"
				:defaultValue="value"
				@input="handleListItemInput"
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
import { tools } from "vue-rocket";
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

		appendIcon: {
			type: String,
			default: "mdi-help-circle-outline"
		},

		appendTip: {
			type: String,
			required: false
		}
	},

	data() {
		return {
			value: null,

			showLength: false,
			allowSpace: false,

			focus: false
		};
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
			for (let i = 0; i < this.validations.length; i += 1) {
				const result =
					tools.isYummy(this.rules[i]) &&
					this.rules[i]({
						value: this.value,
						rule: this.validations[i].rule,
						message: this.validations[i].message
					});

				console.log(result);

				if (result !== true) {
					this.error = true;
				}
			}

			event.customValue = this.value;
			event.customItems = this.items;

			this.$emit("enter", event);
		},

		onFocus(event) {
			this.focus = true;
			event.customValue = this.value;
			this.$emit("focus", event);
		},

		onInput(value) {
			this.typing = true;
			this.$emit("input", value);
		},

		onKeyup(event) {
			if (event.keyCode === 38 && this.typing && this.items.length) {
				return;
			}

			this.value = this.allowSpace ? this.value : this.value.replace(/\s+/g, "");
			event.customValue = this.value;
			this.$emit("keyup", event);
		}
	},

	watch: {
		defaultValue: {
			handler(value) {
				this.value = value;
			},
			immediate: true
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
		cursor: pointer;
	}
}
</style>
