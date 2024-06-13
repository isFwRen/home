<template>
	<div class="top-bar">
		<div class="bar">
			<v-tooltip v-for="item in topBarList" :key="item.text" bottom>
				<template v-slot:activator="{ on, attrs }">
					<span :class="[
						'icon',
						item.icon,
						item.value === 'cut' ? isCut && 'active' : void 0,
						item.value === 'rect' ? isRect && 'active' : void 0,
						item.value === 'text' ? isText && 'active' : void 0,
						item.value === 'link' && 'link',
						degArr.includes(item.name) ? 'fs' : ''
					]" v-bind="attrs" v-on="on" @click="handleEvent(item.value, item.deg)">{{ item.name }}</span>
				</template>

				<span>{{ item.text }}</span>
			</v-tooltip>
		</div>
	</div>
</template>

<script>
import tools from "../../libs/tools";
import { topBarList } from "./cells";
import { localStorage } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";
export default {
	name: "TopBar",

	props: {
		isCut: {
			type: Boolean,
			default: false
		},

		isRect: {
			type: Boolean,
			default: false
		},

		isText: {
			type: Boolean,
			default: false
		},

		src: {
			type: String,
			required: false
		}
	},

	data() {
		return {
			topBarList,
			degArr: ['-5°', '5°', '-10°', '10°', '-15°', '15°']
		};
	},

	methods: {
		handleEvent(eventName, deg) {
			// 查看原图
			if (eventName === "link") {
				const token = localStorage.get("token");
				const user = localStorage.get("user");
				const secret = localStorage.get("secret") || "";
				let code = "";
				if (secret) {
					code = lpTools.GetCode(secret);
				}
				var xhr = new XMLHttpRequest();
				xhr.withCredentials = true;
				xhr.open("GET", this.src);
				xhr.responseType = "blob"; // 设置响应类型为二进制数据
				xhr.setRequestHeader("X-Token", token); // 设置自定义请求头
				xhr.setRequestHeader("X-User-Id", user); // 设置自定义请求头
				xhr.setRequestHeader("X-Code", String(code)); // 设置自定义请求头

				xhr.onreadystatechange = function () {
					if (xhr.readyState === 4) {
						var blob = xhr.response;
						var imageUrl = URL.createObjectURL(blob);
						window.open(imageUrl);
					}
				};
				xhr.send();
			}

			tools.debounce(() => {
				this.$emit("topBarEvent", eventName, deg);
			}, 150);
		}
	}
};
</script>

<style scoped lang="scss">
@import "../../icons/iconfont.css";

.top-bar {
	background-color: #282828;
	border-bottom: 1px solid #171717;

	.bar {
		display: flex;
		align-items: center;
		position: relative;
		height: 28px;
		margin: 4px 16px;
	}
}

.icon {
	margin-right: 16px;
	color: #fff;
	font-size: 18px;
	opacity: 0.8;
	cursor: pointer;
	font-size: 18px;

	&.link {
		position: absolute;
		right: 0;
	}

	&.active {
		color: #39b54a;
	}

	&.icon-clear {
		color: #ff5252;
	}

	&.icon-done {
		color: #39b54a;
	}

	&:hover {
		opacity: 1;
	}

	&:active {
		color: #39b54a;
	}

	&.icon-clear:active {
		color: #ff5252;
	}

	&:last-child {
		margin-right: 0;
	}
}

.fs {
	font-size: 16px !important;
}
</style>
