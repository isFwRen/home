<template>
	<div class="op-fields-title">
		<h3 class="mr-1">{{ computedTitle }}</h3>
		<v-icon @click="handleVideo">mdi-youtube</v-icon>

		<teaching-video-dialog
			ref="dialog"
			:blockCode="blockCode"
			:proCode="proCode"
			:videoSource="isOpenSource"
		></teaching-video-dialog>
	</div>
</template>

<script>
export default {
	name: "OpFieldsTitle",

	props: {
		blockCode: {
			type: String,
			required: false
		},

		isLoop: {
			type: Boolean,
			default: false
		},

		proCode: {
			type: String,
			required: false
		},

		title: {
			type: String,
			required: false
		},

		isOpenSource: {
			type: Array,
			required: false
		}
	},

	methods: {
		handleVideo() {
			this.$refs.dialog.onOpen();
		}
	},

	computed: {
		computedTitle() {
			return `(${this.blockCode}${this.isLoop ? " - 循环分块" : ""})${this.title}`;
		}
	},

	watch: {
		isOpenSource: {
			handler(newValue) {
				if (newValue.length != 0) {
					this.$refs.dialog.onOpen();
				}
			}
		}
	},

	components: {
		"teaching-video-dialog": () => import("../teachingVideoDialog")
	}
};
</script>

<style scoped lang="scss">
.op-fields-title {
	display: flex;
	align-items: center;
}
</style>