<template>
	<div class="lp-monaco-editor">
		<div ref="editor" style="width: 100%; height: 100%"></div>
	</div>
</template>

<script>
import * as monaco from "monaco-editor";

export default {
	name: "lp-monaco-editor",

	props: {
		value: {
			type: String,
			required: true
		},
		format: {
			type: String,
			default: "xml",
			required: true
		}
	},

	data() {
		return {
			monacoInstance: null
		};
	},
	watch: {
		value: {
			handler(value) {
				if (!value) return;
				this.createXML(value);
			},
			immediate: true
		}
	},

	methods: {
		createXML(value) {
			console.log("createXML1");
			this.$nextTick(() => {
				//this.monacoInstance && this.monacoInstance.dispose();//使用完成销毁实例
				if (!this.monacoInstance) {
					this.monacoInstance = monaco.editor.create(this.$refs.editor, {
						theme: "vs-dark",
						value,
						wordWrap: "on",
						automaticLayout: true,
						language: this.format,
						autoIndent: true,
						formatOnType: true,
						formatOnPaste: true,
						formatDocument: true
					});
				}

				const timer = setTimeout(() => {
					this.monacoInstance
						.getAction("editor.action.formatDocument")
						.run()
						.then(function () {
							console.log("Format completed");
							clearTimeout(timer);
						});
				}, 100);
			});
		},

		getXML() {
			return this.monacoInstance.getValue();
		}
	}
};
</script>

<style lang="scss">
.lp-monaco-editor {
	width: 100%;
	height: 100%;
}
</style>
