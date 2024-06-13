<template>
	<div style="border: 1px solid #ccc">
		<Toolbar
			style="border-bottom: 1px solid #ccc"
			:editor="editor"
			:defaultConfig="toolbarConfig"
			:mode="mode"
		/>
		<Editor
			style="height: 500px; overflow-y: hidden"
			v-model="html"
			:defaultConfig="editorConfig"
			:mode="mode"
			@onCreated="onCreated"
		/>
	</div>
</template>
<script>
import Vue from "vue";
import { Editor, Toolbar } from "@wangeditor/editor-for-vue";

export default Vue.extend({
	components: { Editor, Toolbar },
	data() {
		return {
			editor: null,
			html: this.content,
			toolbarConfig: {},
			editorConfig: {
				placeholder: "请输入内容..."
			},
			mode: "default" // or 'simple'
		};
	},
	props: ["content"],
	methods: {
		onCreated(editor) {
			var that = this;
			editor.getMenuConfig("uploadImage").customUpload =
				// 自定义上传
				async function (file, insertFn) {
					// file 即选中的文件
					// 自己实现上传，并得到图片 url alt href
					// 最后插入图片
					const result = await that.$store.dispatch("NOTIC_EDIT_FILE_UPLOAD", {
						file: file
					});
					if (result.code == 200) {
						const file = result.data.file;
						insertFn(file.url, file.name, "");
					}
				};
			this.editor = Object.seal(editor); // 一定要用 Object.seal() ，否则会报错
			console.log(this.editor.getMenuConfig("uploadImage"));
		}
	},
	watch: {
		html: function (newHtml) {
			this.$emit("getMsg", newHtml);
		},
		content: function (newContent) {
			this.html = newContent;
		}
	},
	beforeDestroy() {
		const editor = this.editor;
		if (editor == null) return;
		editor.destroy(); // 组件销毁时，及时销毁编辑器
	}
});
</script>

<style scoped>
@import "~@wangeditor/editor/dist/css/style.css";
</style>
