<template>
	<div class="one-code-front-chunk">
		<v-row class="exchange-container">
			<v-col cols="6">
				<h6 class="text-h6 fw-bold title">一码前置分块</h6>

				<v-row class="mr-n1 exchange-items">
					<v-col :cols="12">
						<draggable
							class="list-group"
							tag="ul"
							v-model="onCodeFrontChunkList"
							v-bind="dragOptions"
							@start="drag = true"
							@end="drag = false"
						>
							<transition-group type="transition" :name="!drag ? 'flip-list' : null">
								<div
									class="list-group-item"
									v-for="(item, index) in onCodeFrontChunkList"
									:key="index"
								>
									<v-chip
										v-show="editTemp.chunkId !== item.ID"
										class="mr-2 mb-2"
										color="primary"
										draggable
										label
										outlined
									>
										{{ item.name }}
										<v-icon size="16" class="ml-1 mr-n1" color="grey"
											>mdi-minus-circle</v-icon
										>
									</v-chip>
								</div>
							</transition-group>
						</draggable>
					</v-col>
				</v-row>
			</v-col>

			<v-col cols="6">
				<h6 class="ml-2 text-h6 fw-bold title">项目分块</h6>

				<v-row class="ml-n1 exchange-items">
					<v-col :cols="12">
						<draggable
							class="list-group"
							tag="ul"
							v-model="chunkList"
							v-bind="dragOptions"
							@start="drag = true"
							@end="drag = false"
						>
							<transition-group type="transition" :name="!drag ? 'flip-list' : null">
								<div
									class="list-group-item"
									v-for="(item, index) in chunkList"
									:key="index"
								>
									<v-chip
										v-show="editTemp.chunkId !== item.ID"
										class="mr-2 mb-2"
										color="primary"
										draggable
										label
										outlined
									>
										{{ item.name }}
										<v-icon size="16" class="ml-1 mr-n1"
											>mdi-plus-circle</v-icon
										>
									</v-chip>
								</div>
							</transition-group>
						</draggable>
					</v-col>
				</v-row>
			</v-col>
		</v-row>
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import draggable from "vuedraggable";

export default {
	name: "OneCodeFrontChunk",

	data() {
		return {
			drag: false,
			onCodeFrontChunkList: [],
			chunkList: []
		};
	},

	methods: {
		async saveExchange() {
			const { chunkId } = this.storage.get("config");

			const form = {
				blockId: chunkId,
				myType: 1,
				tempBlockRelationArr: []
			};

			if (this.onCodeFrontChunkList.length) {
				for (let item of this.onCodeFrontChunkList) {
					form.tempBlockRelationArr.push({
						myType: 1,
						preBCode: item.code,
						preBId: item.ID,
						preBName: item.name,
						tempBId: chunkId
					});
				}
			}

			const result = await this.$store.dispatch("UPDATE_CHUNK_CHUNK_RELATION", form);

			this.toasted.dynamic(result.msg, result.code);
		}
	},

	computed: {
		...mapGetters(["editTemp"]),

		dragOptions() {
			return {
				animation: 200,
				group: "description",
				disabled: false,
				ghostClass: "ghost"
			};
		}
	},

	watch: {
		"editTemp.onCodeFrontChunkList": {
			handler(list) {
				this.onCodeFrontChunkList = [...list];
			},
			immediate: true
		},

		"editTemp.chunkList": {
			handler(list) {
				this.chunkList = [...list];
			},
			immediate: true
		}
	},

	components: {
		draggable
	}
};
</script>

<style scoped lang="scss">
@import "../exchange.scss";
</style>
