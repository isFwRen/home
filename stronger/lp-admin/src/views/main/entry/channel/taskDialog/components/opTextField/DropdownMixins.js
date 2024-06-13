import { tools } from "vue-rocket";

export default {
	props: {
		items: {
			type: Array,
			required: false
		}
	},

	data() {
		return {
			typing: false,
			dropdownItems: []
		};
	},

	methods: {
		handleListItemInput({ item, index }) {
			if (this.typing) {
				this.$emit("dropdownEnter", {
					index: index,
					value: item
				});

				this.$refs[`input_${this.id}`].focus();
				this.typing = false;
			}
		}
	},

	watch: {
		items: {
			handler(items) {
				if (items.length) {
					this.dropdownItems = tools.deepClone(items);
				} else {
					this.dropdownItems = [];
				}
			},
			immediate: true
		}
	}
};
