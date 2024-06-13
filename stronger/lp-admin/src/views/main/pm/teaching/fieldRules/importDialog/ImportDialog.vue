<template>
	<div class="import-dialog">
		<lp-dialog ref="dialog" :title="title" width="700" @dialog="handleDialog">
			<div class="pt-6" slot="main">
				<z-upload
					:action="`${baseURLApi}${path}`"
					:defaultValue="[{ url: `${baseURLApi}${this.imagePath}` }]"
					:effectData="{
						proCode,
						id
					}"
					:formId="formId"
					:formKey="id"
					:limit="limit"
					:headers="fileHeaders"
					:max-size="2048"
					:multiple="!id"
					parcel
					@response="handleResponse"
					@click="getUserInfo"
				></z-upload>
			</div>

			<div class="z-flex" slot="actions">
				<z-btn class="mr-3" color="primary" outlined @click="onClose">关闭</z-btn>
			</div>
		</lp-dialog>
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

		imagePath: {
			type: String,
			required: true
		},

		limit: {
			type: [Number, String],
			default: 1
		},

		path: {
			type: String,
			required: true
		},

		proCode: {
			tyep: String,
			required: false
		}
	},

	data() {
		return {
			formId: "TeachingFieldRulesImportDialog",
			baseURLApi,
			fileHeaders: {},
			url: "",
			images: []
		};
	},

	created() {
		this.getUserInfo();
	},

	methods: {
		getUserInfo() {
			const { token, user } = localStorage.get(["token", "user"]);
			const secret = localStorage.get("secret");
			let code = "";
			if (secret) {
				code = lpTools.GetCode(secret);
			}
			this.fileHeaders = {
				"x-token": token,
				"x-user-id": user.id,
				"x-code": String(code)
			};
		},

		handleResponse({ result }) {
			if (result.maxSize) {
				this.toasted.warning("不能超过2M!");
				return;
			}

			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				this.onClose();
			}
		}
	},

	watch: {
		dialog: {
			handler(dialog) {
				if (dialog) {
					this.images = [{ url: `${baseURLApi}${this.imagePath}` }];
				}
			},
			immediate: true
		}
	}
};
</script>
