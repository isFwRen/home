<template>
	<div class="view-dialog">
		<lp-dialog ref="dialog" :title="title" width="700" @dialog="handleDialog">
			<div class="pt-6" slot="main">
				<img width="100%" :src="src" />
			</div>

			<div class="z-flex" slot="actions">
				<z-btn class="mr-3" color="primary" outlined @click="onClose">关闭</z-btn>
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
		imagePath: {
			type: String,
			required: true
		}
	},

	data() {
		return {
			baseURLApi,
			src: ""
		};
	},

	watch: {
		imagePath: {
			async handler(src) {
				const blob = await lpTools.getTokenImg(`${this.baseURLApi}${src}`);
				lpTools.getBase64(blob).then(base64String => {
					this.src = base64String;
				});
			}
		},
		immediate: true
	}
};
</script>
