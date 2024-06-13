<template>
	<div class="search-wrap" v-if="dialog" v-drag>
		<v-card>
			<v-card-title class="text-h5 grey drag-class lighten-2"> 查询结果 </v-card-title>
			<v-card-text>
				<v-container>
					<v-row>
						<v-col cols="12" sm="12" md="12">
							<v-simple-table fixed-header height="240px">
								<template v-slot:default>
									<thead>
										<tr>
											<th class="text-left">页数</th>
											<th class="text-left">分块名</th>
											<th class="text-left">字段名</th>
											<th class="text-left">字段值</th>
											<th class="text-left">跳转</th>
										</tr>
									</thead>
									<tbody v-if="resultData.length">
										<tr
											v-for="(item, index) in resultData"
											:key="index"
										>
											<td>第{{ item.index + 1 }}页</td>
											<td>{{ item.title }}</td>
											<td>{{ item.code }}{{ item.name }}</td>
											<td>{{ item.resultValue }}</td>
											<td><v-btn elevation="2" @click="handleClick(item)">跳转</v-btn></td>
										</tr>
									</tbody>
									<div v-else>
										<div style="position: relative; left: 150px; top: 100px">
											未搜索到数据
										</div>
									</div>
								</template>
							</v-simple-table>
						</v-col>
					</v-row>
				</v-container>
			</v-card-text>

			<v-divider></v-divider>
			<v-card-actions>
				<v-spacer></v-spacer>
				<v-btn color="primary" text @click="dialog = false">关闭</v-btn>
			</v-card-actions>
		</v-card>
	</div>
</template>
<script>
export default {
	name: "searchResultData",
	data() {
		return {
			dialog: false
		};
	},
	props: {
		resultData: Array
	},
	directives: {
		drag: {
			bind(el, binding) {
				el.onmousedown = e => {
					if (!e.target.className.includes("drag-class")) {
						return;
					}
					e.preventDefault();
					const viewL = 0;
					const viewR = window.innerWidth;
					const viewT = 0;
					const viewB = window.innerHeight;

					e.stopPropagation();
					let defaultX = e.clientX;
					let defaultY = e.clientY;

					let defaultLeft = el.offsetLeft;
					let defaultTop = el.offsetTop;
					el.onmousemove = me => {
						me.stopPropagation();
						el.style.cursor = "move";

						let nowX = me.clientX;
						let nowY = me.clientY;

						let moveX = defaultX - nowX;
						let moveY = defaultY - nowY;

						let nowLeft = defaultLeft - moveX;
						let nowTop = defaultTop - moveY;

						if (nowLeft <= viewL) {
							nowLeft = viewL;
						}
						if (nowTop <= viewT) {
							nowTop = viewT;
						}

						if (viewB - nowTop - el.clientHeight <= 0) {
							nowTop = viewB - el.clientHeight;
						}

						if (viewR - nowLeft - el.clientWidth <= 0) {
							nowLeft = viewR - el.clientWidth;
						}

						el.style.left = `${nowLeft}px`;
						el.style.top = `${nowTop}px`;
					};

					el.onmouseup = event => {
						event.preventDefault();
						el.style.cursor = "default";
						el.onmousemove = null;
						el.onmouseup = null;
					};
				};
			}
		}
	},
	methods: {
		handleClick(item) {
			this.$emit('toLink',item)
		},
		open() {
			this.dialog = true;
		},
		close() {
			this.dialog = true;
		}
	}
};
</script>
<style scoped lang="scss">
.drag-class {
	cursor: move;
}
.search-wrap {
	width: 460px;
	position: fixed;
	right: 30px;
	top: 125px;
	z-index: 300;
}
</style>
