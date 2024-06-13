export default {
	data() {
		return {
			dialog: false
		};
	},

	methods: {
		onClose() {
			this.$refs.dialog.onClose();
		},

		onOpen() {
			this.$refs.dialog.onOpen();
		}
	}
};
