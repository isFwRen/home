<template>
	<div class="contents" ref="dragEle" :style="moveStyle">
		<div class="contents-icon"><v-icon>mdi-menu</v-icon></div>
		<div class="contents-container">
			<h6 class="contents-title" ref="dragTarget">目录</h6>
			<div class="contents-wrap">
				<div
					class="item"
					v-for="(item, index) in data"
					:class="[activeIndex === index ? 'li--active' : '']"
					:key="index"
					@click="handleContents(item.index)"
				>
					{{ item.title }}
				</div>
			</div>
		</div>
	</div>
</template>
<script>
import dragMixin from "@/mixins/dragMixin";
export default {
	name: "contents",
	mixins: [dragMixin],
	data() {
		return {
		
		};
	},
	props: ["data", "activeIndex"],
	methods: {
		handleContents(index) {
			this.$emit("toContents", index);
		}
	}
};
</script>
<style scoped lang="scss">
.li--active {
	color: #1e80ff !important;
}

@keyframes move {
	0% {
		transform: translate(-5px);
	}
	50% {
		transform: translate(15px);
	}
	100% {
		transform: translate(-5px);
	}
}
.contents {
	position: fixed;
	z-index: 100;
	left: 0;
	top: 200px;
	&:hover {
		.contents-icon {
			opacity: 0;
			visibility: hidden;
		}
		.contents-container {
			opacity: 1;
			visibility: visible;
		}
	}

	.contents-icon {
		position: absolute;
		left: 0;
		top: 0;
		animation-name: move;
		animation-duration: 0.3s;
		animation-delay: 1s;
		animation-iteration-count: 4;
		animation-timing-function: ease-out;
		transition: opacity 0.45s, visibility 0.45s ease-out;
		visibility: visible;
		opacity: 1;
		width: 50px;
		height: 50px;
		border-radius: 50%;
		background-color: #fff;
		border: 1px solid #eee;
		line-height: 45px;
		font-weight: 600;
		text-align: center;
	}
	.contents-container {
		transition: opacity 0.45s, visibility 0.45s ease-out;
		opacity: 0;
		visibility: hidden;
		position: absolute;
		left: 0;
		top: 0;
		user-select: none;
		border-radius: 5px;
		background-color: #fff;
		padding: 20px;
		box-sizing: border-box;
		border: 1px solid #eee;
		box-shadow: 0px 0px 4px rgba(0, 0, 0, 0.12);
		.contents-title {
			cursor: move;
			font-weight: 600;
			padding-bottom: 10px;
			border-bottom: 1px solid #e4e6eb;
		}
		.contents-wrap {
			margin-top: 10px;
			padding: 0 20px;
			padding-left: 0 !important;
			max-height: 300px;
			min-width: 140px;
			overflow-y: auto;
			.item {
				color: #252933;
				cursor: pointer;
				padding: 3px 0;
				&:hover {
					color: #1e80ff;
				}
			}
		}
	}
}
/* 修改滚动条的样式 */
::-webkit-scrollbar {
	width: 8px; /* 设置滚动条宽度 */
}

::-webkit-scrollbar-thumb {
	background-color: #888; /* 设置滚动条滑块颜色 */
	border-radius: 4px; /* 设置滚动条滑块圆角 */
}

::-webkit-scrollbar-track {
	background-color: #f1f1f1; /* 设置滚动条背景色 */
}

/* 鼠标悬停在滚动条上时的样式 */
::-webkit-scrollbar-thumb:hover {
	background-color: #555;
}
</style>
