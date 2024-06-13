<template>
	<div class="options">
		<ul class="z-flex align-center">
			<div
				v-for="(item, index) in options"
				:key="item.key"
				:class="{ 'divider z-flex align-center': index === 5 }"
			>
				<li class="mr-2">
					<v-tooltip bottom>
						<template v-slot:activator="{ on, attrs }">
							<v-btn
								:color="item.color"
								icon
								:disabled="item.disabled"
								v-bind="attrs"
								v-on="on"
								@click="onSelect(item)"
							>
								<v-icon :size="item.size" :style="item.style">{{
									item.icon
								}}</v-icon>
							</v-btn>
						</template>
						<span>{{ item.tips }}</span>
					</v-tooltip>
				</li>
			</div>
		</ul>
	</div>
</template>

<script>
import options from "./cells";

export default {
	name: "CanvasToolbar",

	data() {
		return {
			options
		};
	},

	methods: {
		onSelect(item) {
			this.$emit("select", item);
		}
	}
};
</script>

<style scoped lang="scss">
.options {
	position: absolute;
	top: 0;
	left: 50%;
	border: 1px solid rgba(0, 0, 0, 0.06);
	background-color: #fff;
	transform: translateX(-50%);
	z-index: 1;

	ul,
	li {
		padding: 0;
		margin: 0;
		list-style-type: none;
	}

	.divider::before {
		content: "";
		display: inline-block;
		margin-right: 8px;
		width: 1px;
		height: 20px;
		background-color: rgba(0, 0, 0, 0.15);
	}
}
</style>
