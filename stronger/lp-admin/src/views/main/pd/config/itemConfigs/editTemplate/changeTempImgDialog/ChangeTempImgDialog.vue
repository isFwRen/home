<template>
	<lp-dialog ref="dialog" title="修改模板图片" width="700" @dialog="handleDialog">
		<div class="pt-4" slot="main">
			<z-file-input
				:formId="formId"
				formKey="upload"
				:action="`${baseURLApi}pro-config/sys-template/update-images`"
				prepend-icon="mdi-file-image-outline"
				:fileList="list"
				:effectData="{
					sysProTempId: storage.get('config').tempId,
					proCode: config.pro.code
				}"
				:headers="fileHeaders"
				placeholder="点击上传图片"
				:deleteIcon="false"
				parcel
				multiple
				name="tempImages"
				@click="getUserInfo"
				@response="handleResponse"
			></z-file-input>
		</div>
	</lp-dialog>
</template>

<script>
import { mapGetters } from "vuex";
import DialogMixins from "@/mixins/DialogMixins";
import { tools as lpTools } from "@/libs/util";

const { baseURLApi } = lpTools.baseURL();

export default {
	name: "ChangeTempImgDialog",
	mixins: [DialogMixins],

	props: {
		imageList: {
			type: Array,
			default: () => []
		}
	},

	data() {
		return {
			formId: "ChangeTempImgDialog",
			baseURLApi,
			list: [],
			user: {},
			fileHeaders: {},
			token: ""
		};
	},
	created() {
		this.getUserInfo();
	},
	methods: {
		// 获取用户信息
		getUserInfo() {
			const user = this.storage.get("user");
			const token = this.storage.get("token");
			const secret = this.storage.get("secret");
			let code = "";
			if (secret) {
				code = lpTools.GetCode(secret);
			}

			this.user = user;
			this.fileHeaders = {
				"x-token": token,
				"x-user-id": user.id,
				"pro-code": this.project.code,
				"x-code": String(code)
			};
		},
		handleResponse({ result }) {
			console.log(result, "result");
			if (result.code === 200) {
				this.toasted.success(result.msg);
				this.$emit("uploaded");
			} else {
				this.toasted.warning("由于不可抗拒力量，导致图片上传失败!");
			}
		}
	},

	computed: {
		...mapGetters(["config", "project"])
	},

	watch: {
		imageList: {
			handler() {
				this.list = [];
				if (this.imageList.length) {
					for (let image of this.imageList) {
						this.list.push({
							url: `${baseURLApi}${image}`,
							name: image
						});
					}
				}
			},
			immediate: true
		},

		dialog(dialog) {
			if (dialog) {
				this.user = this.storage.get("user");
				this.token = this.storage.get("token");
				console.log(this.project);
				console.log(this.user);
			}
		}
	}
};
</script>
