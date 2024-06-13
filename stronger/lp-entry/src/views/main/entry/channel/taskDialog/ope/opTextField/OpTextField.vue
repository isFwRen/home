<template>
	<div class="op-text-field" :id="id">
		<div class="z-flex align-center op-text-field-label">
			<!-- 字段名称 BEGIN -->
			<div class="label">{{ label }}</div>
			<!-- 字段名称 END -->
		</div>

		<v-text-field
			:id="`input_${id}`"
			:ref="`input_${id}`"
			v-model="value"
			autocomplete="off"
			:autofocus="autofocus"
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
	mixins: [DropdownMixins],

	props: {
		autofocus: {
			type: Boolean,
			default: false
		},

		defaultValue: {
			type: [String, Number, Boolean, Array, Object],
			default: undefined
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

		id: {
			type: String,
			required: false
		},

		label: {
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
		}
	},

	data() {
		return {
			value: null,
			focus: false
		};
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
				this.value = newVal.optionList[0];
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
		},

		onFocus(event) {
			this.focus = true;
			event.customValue = this.value;

			this.$emit("focus", event);
		},

		onClear() {
			this.$emit("focusClear", -1);
		},

		onInput(value) {
			this.typing = true;

			this.$emit("input", value);
		},

		onKeyup(event) {
			if (event.keyCode === 38 && this.typing && this.items.length) {
				return;
			}
			// this.value = this.value.trim();
			event.customValue = this.value;
			this.$emit("keyup", event);
		}
	}
};
</script>

<style scoped lang="scss">
.op-text-field {
	position: relative;
	min-height: 94px;

	label {
		font-size: 1.25rem;
		color: #000;
	}

	.op-text-field-label {
		position: relative;
		top: 18px;
		z-index: 0;

		.label {
			font-size: 16px;
			font-weight: bold;
			color: black;
		}

		label.label-accent {
			color: #0610f7;
			font-size: 15px;
			font-weight: bolder;
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