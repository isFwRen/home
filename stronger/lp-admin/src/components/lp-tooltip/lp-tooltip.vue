<template>
	<transition name="fade" v-if="content">
		<div
			@click="copyToClipboard"
			v-show="visible"
			ref="tooltipRef"
			class="tooltip"
			@mouseover="onShow"
			@mouseleave="close"
		>
			<div :class="[dir === 'left' ? 'trancle--left' : 'trancle--right']"></div>
			{{ content }}
		</div>
	</transition>
</template>

<script>
export default {
	data() {
		return {
			visible: false,
			content: "",
			dir: "",
			isHover: false,
			stack: []
		};
	},
	beforeDestroy() {
		this.visible = false;
		this.isHover = false;
	},
	methods: {
		onShow() {
			this.isHover = true;
			this.visible = true;
		},
		copyToClipboard(event) {
			navigator.clipboard
				.writeText(this.content)
				.then(() => {
					console.log("Text copied to clipboard");
				})
				.catch(error => {
					console.error("Failed to copy text: ", error);
				});
			event.stopPropagation();
			event.preventDefault();
			return false;
		},
		show(options) {
			this.content = options.content;
			let { styles } = options;
			this.visible = true;
			this.dir = styles.dir;
			if (this.$refs.tooltipRef) {
				this.$refs.tooltipRef.style.left = styles.left + "px";
				this.$refs.tooltipRef.style.top = styles.top + "px";
			}
		},
		close() {
			this.isHover = false;
			this.visible = false;
		},
		hide() {
			if (!this.isHover) {
				this.visible = false;
			}
		}
	}
};
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
	transition: opacity 0.5s;
}
.fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */ {
	opacity: 0;
}
.tooltip {
	position: fixed;
	width: 450px;
	padding: 5px 10px;
	background-color: hsl(0, 0%, 20%);
	color: #fff;
	border-radius: 4px;
	z-index: 9999;
}
.trancle--left {
	position: absolute;
	width: 0;
	height: 0;
	border: 10px solid;
	left: -20px;
	top: 50%;
	transform: translateY(-50%);
	border-color: transparent hsl(0, 0%, 20%) transparent transparent;
}
.trancle--right {
	position: absolute;
	width: 0;
	height: 0;
	border: 10px solid;
	right: -20px;
	top: 50%;
	transform: translateY(-50%);
	border-color: transparent transparent transparent hsl(0, 0%, 20%);
}
</style>
