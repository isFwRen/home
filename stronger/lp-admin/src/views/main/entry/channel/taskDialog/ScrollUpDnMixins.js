let fieldsWrapper = null;
let [titleHeight, fieldHeight, deadLine] = [22, 94, 0];

export default {
	created() {
		if (this.op === "opq") {
			fieldHeight = 138;
		}

		deadLine = fieldHeight * 2 + titleHeight;
	},

	methods: {
		// op0(向上向下滚动)
		scrollUpDn({ field }) {
			this.$nextTick(() => {
				const opTextField = document.getElementById(field?.uniqueId);

				if (opTextField?.offsetTop >= deadLine) {
					fieldsWrapper.scrollTop = opTextField.offsetTop - deadLine;
				}
			});
		},

		// 返回顶部
		scrollToTop() {
			const timer = setTimeout(() => {
				fieldsWrapper = document.getElementById("fieldsWrapper");

				if (fieldsWrapper) {
					fieldsWrapper.scrollTop = 0;
				}

				clearTimeout(timer);
			}, 25);
		}
	}
};
