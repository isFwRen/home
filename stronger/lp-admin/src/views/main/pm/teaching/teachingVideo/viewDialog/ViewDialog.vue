<template>
	<div class="view-dialog">
		<lp-dialog
			ref="dialog"
			cardTextClass="pa-0"
			:cardTextStyle="{
				height: '100vh',
				overflow: 'hidden'
			}"
			fullscreen
			:title="title"
			toolbarColor="#272727"
			@dialog="handleDialog"
		>
			<div class="z-flex main" slot="main">
				<div class="flex-grow-1 wrap">
					<video ref="video" width="100%" height="100%" controls>
						<source type="video/mp4" />
						您的浏览器不支持 HTML5 video 标签。
					</video>
				</div>

				<div class="side">
					<v-list dark dense>
						<v-subheader>视频列表</v-subheader>
						<v-list-item-group v-model="selectedItem" color="primary">
							<v-list-item
								v-for="(item, i) in items"
								:key="i"
								@click="onSelectVideo(item)"
							>
								<v-list-item-content>
									<v-list-item-title v-text="item.label"></v-list-item-title>
								</v-list-item-content>

								<v-list-item-icon v-if="i === selectedItem">
									<v-icon>mdi-equalizer</v-icon>
								</v-list-item-icon>
							</v-list-item>
						</v-list-item-group>
					</v-list>
				</div>
			</div>
		</lp-dialog>
	</div>
</template>

<script>
import { tools as lpTools } from "@/libs/util";
import DialogMixins from "@/mixins/DialogMixins";

const { baseURLApi } = lpTools.baseURL();

export default {
	name: "TeachingFieldRulesViewDialog",
	mixins: [DialogMixins],

	props: {
		proCode: {
			tyep: String,
			required: false
		},

		video: {
			type: Array,
			default: () => []
		}
	},

	data() {
		return {
			baseURLApi,
			selectedItem: 1,
			items: [],

			path: ""
		};
	},

	methods: {
		async getMenuList() {
			const body = {
				pageSize: 20,
				pageIndex: 1,
				proCode: this.proCode,
				rule: "有"
			};

			const result = await this.$store.dispatch("GET_PM_TEACHING_TEACHING_VIDEO_LIST", body);

			if (result.code === 200) {
				this.items = [];

				result.data?.list.map((item, index) => {
					this.items.push({
						label: item.video[0].name,
						path: item.video[0].path,
						selected: this.path === item.video[0].path ? true : false
					});

					if (this.path === item.video[0].path) {
						this.selectedItem = index;
					}
				});
			}
		},
		blobToBase64(blob) {
			return new Promise((resolve, reject) => {
				const fileReader = new FileReader();
				fileReader.onload = e => {
					resolve(e.target.result);
				};
				// readAsDataURL
				fileReader.readAsDataURL(blob);
				fileReader.onerror = () => {
					reject(new Error("blobToBase64 error"));
				};
			});
		},
		async onSelectVideo({ path }) {
			this.path = path;
			const blob = await lpTools.getTokenImg(`${this.baseURLApi}${path}`);
			this.blobToBase64(blob).then(res => {
				this.$refs.video.src = res;
			});
		}
	},

	watch: {
		dialog: {
			handler(dialog) {
				if (dialog) {
					this.path = this.video?.[0]?.path;
					this.$nextTick(async () => {
						const blob = await lpTools.getTokenImg(`${this.baseURLApi}${this.path}`);
						this.blobToBase64(blob).then(res => {
							this.$refs.video.src = res;
						});
					});
					this.getMenuList();
				}
			},
			immediate: true
		}
	}
};
</script>

<style scoped lang="scss">
.main {
	padding-top: 64px;
	height: 100vh;
	box-sizing: border-box;
	border-bottom: 1px solid rgba(0, 0, 0, 0.2);
	background-color: #121212;

	.side {
		overflow-y: auto;
		border-left: 1px solid rgba(0, 0, 0, 0.2);
		background-color: #1e1e24;
	}
}

@media screen and (min-width: 1280px) {
	.main {
		width: 100vw;

		.side {
			width: 400px;
		}
	}
}

@media screen and (max-width: 1280px) {
	.main {
		width: 1280px;

		.side {
			width: 340px;
		}
	}
}
</style>
