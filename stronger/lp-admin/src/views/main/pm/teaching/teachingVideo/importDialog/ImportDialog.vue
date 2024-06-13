<template>
	<div class="import-dialog">
		<lp-dialog ref="dialog" :title="title" width="700" @dialog="handleDialog">
			<div class="pt-6" slot="main">
				<z-file-input
					formId="files"
					formKey="file"
					accept="video/*"
					:auto-upload="false"
					:multiple="!id"
					prependIcon="mdi-file-excel-outline"
					:defaultValue="videos"
					label="视频上传"
					@click="getUserInfo"
					@change="onSelectVideos"
				>
				</z-file-input>
				<div class="z-flex align-center">
					<div class="label">是否必学</div>
					<div class="radio" style="padding-left: 20px">
						<v-radio-group v-model="isRequired" row>
							<v-radio label="是" :value="1"></v-radio>
							<v-radio label="否" :value="0"></v-radio>
						</v-radio-group>
					</div>
				</div>
			</div>

			<div class="z-flex" slot="actions">
				<z-btn class="mr-3" color="primary" outlined @click="onOK">确定</z-btn>
				<z-btn class="mr-3" color="primary" outlined @click="onClose">关闭</z-btn>
			</div>
		</lp-dialog>

		<lp-spinners :overlay="overlay"></lp-spinners>
	</div>
</template>

<script>
import { localStorage } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";
import DialogMixins from "@/mixins/DialogMixins";

const { baseURLApi } = lpTools.baseURL();

export default {
	name: "TeachingFieldRulesImportDialog",
	mixins: [DialogMixins],

	props: {
		id: {
			type: String,
			required: false
		},

		limit: {
			type: [Number, String],
			default: 1
		},

		proCode: {
			type: String,
			required: false
		},

		sysBlockName: {
			type: String,
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
			fileHeaders: {},
			url: "",
			videos: [],
			overlay: false,
			isRequired: 1,
			file: null
		};
	},

	created() {
		this.getUserInfo();
	},

	methods: {
		getUserInfo() {
			const { token, user } = localStorage.get(["token", "user"]);

			this.fileHeaders = {
				"x-token": token,
				"x-user-id": user.id
			};
		},

		async onOK() {
			const body = {
				isRequired: this.isRequired,
				id: this.id || "",
				file: this.file,
				proCode: this.forms['TeachingVideoSearch'].proCode
			};
	
			
			const uploadVideos = async () => {
				const result = await this.$store.dispatch(
					"UPLOAD_PM_TEACHING_TEACHING_VIDEO_VIDEOS",
					body
				);
				this.toasted.dynamic(result.msg, result.code);
				this.onClose();
			};

			let name = this.file.name?.split(".")[0];
			if(!name){
				name = this.file[0].name?.split(".")[0];
			}

			if (name === this.sysBlockName) {
				await uploadVideos();
			} else {
				this.toasted.warning("分块名称不正确!");
			}
		},
		async onSelectVideos(files) {
			this.overlay = true;

			// this.onClose();

			// const uploadVideos = async () => {
			// 	const body = {
			// 		proCode: this.proCode || "",
			// 		id: this.id || "",
			// 		file: files
			// 	};

			// 	const result = await this.$store.dispatch(
			// 		"UPLOAD_PM_TEACHING_TEACHING_VIDEO_VIDEOS",
			// 		body
			// 	);

			// 	this.toasted.dynamic(result.msg, result.code);
			// };

			// if (this.id) {
			// 	const name = files.name?.split(".")[0];

			// 	if (name === this.sysBlockName) {
			// 		await uploadVideos();

			// 	} else {
			// 		this.toasted.warning("分块名称不正确!");
			// 	}
			// } else {
			// 	await uploadVideos();
			// }

			this.overlay = false;

			this.file = files;
		}
	},

	watch: {
		dialog: {
			handler(dialog) {
				if (dialog) {
					this.videos = [];

					this.video?.map(item => {
						this.videos.push({ label: item.name, url: `${baseURLApi}${item.path}` });
					});
				}
			},
			immediate: true
		}
	},

	components: {
		"lp-spinners": () => import("@/components/lp-spinners")
	}
};
</script>
