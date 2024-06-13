<template>
	<div class="lp-tabs">
		<v-tabs :value="defaultValue">
			<v-tab
				v-for="(item, index) in options"
				:key="`lp_tab_${index}`"
				@change="onChange(item)"
				@click="onClick(item)"
				>{{ item.label }}</v-tab
			>
		</v-tabs>
	</div>
</template>

<script>
export default {
	name: "LPTabs",

	props: {
		defaultValue: {
			type: [String, Number],
			default: 0
		},

		options: {
			type: Array,
			required: true
		}
	},

	methods: {
		onChange(item) {
			this.$emit("change", item);
		},

		onClick(item) {
			this.$emit("click", item);
		}
	}
};
</script>
