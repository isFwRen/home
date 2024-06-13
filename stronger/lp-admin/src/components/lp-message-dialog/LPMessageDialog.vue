<template>
	<!-- <transition
    enter-active-class="__fadeIn"
    leave-active-class="__fadeOut "
  > -->
	<div class="dialogs" v-if="dialog">
		<header class="dialog-header">
			<div class="img-flex">
				<span class="header-title">{{ title }} </span>
			</div>
			<div class="img-flex">
				<v-spacer></v-spacer>
				<v-btn text @click="onClose" color="#fff" style="font-size: 18px">
					<v-icon>mdi-close</v-icon>
				</v-btn>
			</div>
		</header>
		<div class="dialogs-main">
			<slot name="main"></slot>
		</div>
	</div>
	<!-- </transition> -->
</template>

<script>
export default {
	name: "LPMessageDialog",
	data() {
		return {
			dialog: false
		};
	},
	props: {
		title: {
			type: String,
			default: ""
		}
	},
	methods: {
		onClose() {
			this.dialog = false;
			this.$emit("close");
		},
		onOpen() {
			this.dialog = true;
		}
	}
};
</script>

<style scoped lang="scss">
@keyframes fadeIn {
	0% {
		opacity: 0;
	}

	to {
		opacity: 1;
	}
}
@keyframes fadeOut {
	0% {
		opacity: 1;
	}

	to {
		opacity: 0;
	}
}

.__fadeIn {
	animation: fadeIn 0.45s ease-out;
}

.__fadeOut {
	animation: fadeOut 0.45s ease-out;
}

.dialogs {
	width: 100%;
	height: 100%;
	background-color: #fff;
	position: fixed;
	left: 0;
	top: 0;
	bottom: 0;
	right: 0;
	z-index: 100;
	overflow: auto;

	.dialog-header {
		width: 100%;
		height: 64px;
		background-color: #1976d2 !important;
		display: flex;
		font-size: 19px;
		font-weight: 500;
		color: #fff;
		padding-left: 20px;
		box-shadow: 0 2px 4px -1px rgba(0, 0, 0, 0.2), 0 4px 5px 0 rgba(0, 0, 0, 0.14),
			0 1px 10px 0 rgba(0, 0, 0, 0.12);

		.img-flex {
			width: 50%;
			display: flex;
			align-items: center;
		}
	}

	.dialogs-main {
		padding: 0 25px;
	}
}
</style>
