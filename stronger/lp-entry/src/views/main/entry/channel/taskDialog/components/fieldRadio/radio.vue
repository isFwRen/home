<template>
	<div class="radio">
		<!-- 字段名称 BEGIN -->
		<div class="title">{{ label }}</div>
		<!-- 字段名称 END -->
		<div>
			<v-radio-group :id="`input_${id}`" :ref="`input_${id}`" v-model="value" @change="onChange">
				<div class="box">
					<v-radio
						v-for="(el, index) in options"
						:key="index"
						:label="el"
						:value="el"
						class="smallBox"
						color="primary"
					></v-radio>
				</div>
			</v-radio-group>
		</div>

		<div class="slot">
			<slot></slot>
		</div>
	</div>
</template>

<script>
export default {
	name: "FieldRadio",
	mixins: [],

	props: {
		autofocus: {
			type: Boolean,
			default: false
		},

		radioOptions: {
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
			value: "",
			focus: false,
			options: this.radioOptions
		};
	},


	methods: {
		onChange(value) {
			this.$emit("change", [value]);
		}
	}
};
</script>

<style scoped lang="scss">
.radio {
	.title {
		font-size: 16px !important;
		font-weight: bold;
		color: black;
	}

	.box {
		width: 100%;
		display: flex;
		justify-content: flex-start;
		flex-wrap: wrap;

		.smallBox {
			width: 49%;
		}
	}
}

.v-input--selection-controls {
	margin-top: 10px;
}

::v-deep .v-label {
	color: black !important;
}
</style>