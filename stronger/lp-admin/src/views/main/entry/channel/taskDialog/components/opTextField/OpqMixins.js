export default {
	props: {
		firstDiffIndex: {
			type: Number,
			default: -1
		}
	},

	watch: {
		// 找到问题件第一个?
		firstDiffIndex: {
			handler(index) {
				this.$nextTick(() => {
					if (this.autofocus && this.id && index > -1) {
						const el = document.getElementById(this.id);
						this._tools.setCursorPosition(el, index);
					}
				});
			},
			immediate: true
		}
	}
};
