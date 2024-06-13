<template>
	<lp-dialog ref="dialog" :title="title" fullscreen persistent @dialog="handleDialog">
		<div class="pt-6" slot="main">
			<div class="z-flex">
				<lp-tabs class="mb-6" :options="cells.tabsOptions" @change="onTab"></lp-tabs>

				<v-spacer></v-spacer>

				<z-btn class="mt-5 mr-n3" color="primary" small @click="onSave"> 保存 </z-btn>
			</div>

			<field-config v-if="currentTab === 1" ref="field"></field-config>

			<front-chunk v-else-if="currentTab === 2" ref="front"></front-chunk>

			<one-code-front-chunk v-else-if="currentTab === 3" ref="onCode"></one-code-front-chunk>

			<suggest-chunk v-else ref="suggest"></suggest-chunk>
		</div>
	</lp-dialog>
</template>

<script>
import { mapGetters } from "vuex";
import DialogMixins from "@/mixins/DialogMixins";
import ConfigMixins from "../../../ConfigMixins";
import cells from "./cells";

const refs = new Map([
	[1, "field"],
	[2, "front"],
	[3, "onCode"],
	[4, "suggest"]
]);

export default {
	name: "ExchangeFieldChunk",
	mixins: [DialogMixins, ConfigMixins],

	data() {
		return {
			cells,
			currentTab: 1
		};
	},

	methods: {
		onTab(item) {
			this.currentTab = item.value;
			this.getAllFieldList();
			this.getChunkFieldList();
		},

		onSave() {
			this.$refs[refs.get(this.currentTab)].saveExchange();
		},

		// 项目字段
		async getAllFieldList() {
			const result = await this.$store.dispatch("GET_CONFIG_ALL_FIELD_LIST", {
				proId: this.proId
			});

			if (result.code === 200) {
				this.$store.commit("UPDATE_EDIT_TEMP", {
					fieldList: result.data
				});
			}
		},

		// 模板字段
		async getChunkFieldList() {
			const { chunkId } = this.storage.get("config");

			const result = await this.$store.dispatch("GET_CONFIG_TEMP_CHUNK_FIELD_LIST", {
				id: chunkId
			});

			if (result.code === 200) {
				const { tempBFRelation, oneBlockRelation, twoBlockRelation, threeBlockRelation } =
					result.data;

				// 字段配置
				const tempFieldList = [];

				if (tempBFRelation.length) {
					for (let item of tempBFRelation) {
						tempFieldList.push({
							ID: item.fId,
							code: item.fCode,
							name: item.fName
						});
					}
				}

				// 前置分块
				const frontChunkList = [];

				if (oneBlockRelation.length) {
					for (let item of oneBlockRelation) {
						frontChunkList.push({
							ID: item.preBId,
							code: item.preBCode,
							name: item.preBName
						});
					}
				}

				// 一码前置分块
				const onCodeFrontChunkList = [];

				if (twoBlockRelation.length) {
					for (let item of twoBlockRelation) {
						onCodeFrontChunkList.push({
							ID: item.preBId,
							code: item.preBCode,
							name: item.preBName
						});
					}
				}

				// 参考分块
				const suggestChunkList = [];

				if (threeBlockRelation.length) {
					for (let item of threeBlockRelation) {
						suggestChunkList.push({
							ID: item.preBId,
							code: item.preBCode,
							name: item.preBName
						});
					}
				}

				this.$store.commit("UPDATE_EDIT_TEMP", {
					tempFieldList,
					frontChunkList,
					onCodeFrontChunkList,
					suggestChunkList
				});
			}
		}
	},

	computed: {
		...mapGetters(["config"])
	},

	watch: {
		dialog: {
			handler(dialog) {
				if (!dialog) {
					this.rememberIds({ chunkId: "" });
				} else {
					this.getAllFieldList();
					this.getChunkFieldList();
				}
			},
			immediate: true
		}
	},

	components: {
		"field-config": () => import("./fieldConfig"),
		"front-chunk": () => import("./frontChunk"),
		"one-code-front-chunk": () => import("./oneCodeFrontChunk"),
		"suggest-chunk": () => import("./suggestChunk")
	}
};
</script>

<style lang="scss"></style>
