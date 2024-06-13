<template>
	<div class="top-bar">
		<div class="bar">
			<v-tooltip v-for="item in topBarList" :key="item.value" bottom>
				<template v-slot:activator="{ on, attrs }">
					<v-icon
						v-bind="attrs"
						v-on="on"
						:class="['icon', item.icon]"
						@click="handleEvent(item.value)"
					></v-icon>
				</template>
				<span>{{ item.text }}</span>
			</v-tooltip>
		</div>
	</div>
</template>

<script>
import { tools } from "@/libs/util";
import { topBarList } from "./cells";

export default {
	name: "TopBar",

	data() {
		return {
			topBarList
		};
	},

	methods: {
		handleEvent(eventName) {
			tools.debounce(() => {
				this.$emit("topBarEvent", eventName);
			}, 150);
		}
	}
};
</script>

<style scoped lang="scss">
.top-bar {
	position: absolute;
	width: 100%;
	background-color: rgba(0, 0, 0, 0.4);
	z-index: 2;

	.bar {
		display: flex;
		align-items: center;
		justify-content: flex-end;
		height: 28px;
		margin: 4px 16px;
	}
}

.icon {
	margin-right: 16px;
	color: #fff;
	font-size: 18px;
	opacity: 0.8;
}

.icon:hover {
	opacity: 1;
}

.icon:active {
	color: #39b54a;
}

.icon:last-child {
	margin-right: 0;
}
</style>
