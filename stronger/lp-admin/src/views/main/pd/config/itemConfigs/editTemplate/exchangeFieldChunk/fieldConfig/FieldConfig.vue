<template>
	<div class="pt-2 field-config">
		<v-row class="exchange-container">
			<v-col cols="6">
				<h6 class="text-h6 fw-bold title">模板字段</h6>

				<v-row class="mr-n1 exchange-items">
					<v-col :cols="12">
						<draggable
							class="list-group"
							tag="ul"
							v-model="tempFieldList"
							v-bind="dragOptions"
							@start="drag = true"
							@end="drag = false"
						>
							<transition-group type="transition" :name="!drag ? 'flip-list' : null">
								<div
									class="list-group-item"
									v-for="(item, index) in tempFieldList"
									:key="index"
								>
									<v-tooltip top>
										<template v-slot:activator="{ on, attrs }">
											<v-chip
												class="mr-2 mb-2"
												v-bind="attrs"
												v-on="on"
												color="primary"
												draggable
												label
												outlined
											>
												{{ item.name }}
												<v-icon
													size="16"
													class="ml-1 mr-n1"
													color="grey"
													@click="swap(item, index, false)"
													>mdi-minus-circle</v-icon
												>
											</v-chip>
										</template>
										<span>{{ item.code }}</span>
									</v-tooltip>
								</div>
							</transition-group>
						</draggable>
					</v-col>
				</v-row>
			</v-col>

			<v-col cols="6">
				<h6 class="ml-2 text-h6 fw-bold title">项目字段</h6>

				<v-row class="ml-n1 exchange-items">
					<v-col :cols="12">
						<draggable
							class="list-group"
							tag="ul"
							v-model="fieldList"
							v-bind="dragOptions"
							@start="drag = true"
							@end="drag = false"
						>
							<transition-group type="transition" :name="!drag ? 'flip-list' : null">
								<div
									class="list-group-item"
									v-for="(item, index) in fieldList"
									:key="index"
								>
									<v-tooltip top>
										<template v-slot:activator="{ on, attrs }">
											<v-chip
												class="mr-2 mb-2"
												v-bind="attrs"
												v-on="on"
												color="primary"
												draggable
												label
												outlined
											>
												{{ item.name }}
												<v-icon
													size="16"
													class="ml-1 mr-n1"
													@click="swap(item, index, true)"
													>mdi-plus-circle</v-icon
												>
											</v-chip>
										</template>
										<span>{{ item.code }}</span>
									</v-tooltip>
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
	name: "FieldConfig",

	data() {
		return {
			drag: false,
			tempFieldList: [],
			fieldList: []
		};
	},

	methods: {
		async saveExchange() {
			const { chunkId } = this.storage.get("config");

			const form = {
				blockId: chunkId,
				tempBFRelationArr: []
			};

			if (this.tempFieldList.length) {
				for (let item of this.tempFieldList) {
					form.tempBFRelationArr.push({
						fCode: item.code,
						fId: item.ID,
						fName: item.name,
						tempBId: chunkId
					});
				}
			}

			const result = await this.$store.dispatch("UPDATE_CHUNK_FIELD_RELATION", form);
			this.toasted.dynamic(result.msg, result.code);
		},

		// isAdd === true; 项目字段 → 模板字段
		// isAdd === false; 模板字段 → 项目字段
		swap(item, index, isAdd) {
			let arrName = ["fieldList", "tempFieldList"];
			if (isAdd) {
				arrName = ["tempFieldList", "fieldList"];
			}

			this[arrName[0]].push({
				code: item.code,
				ID: item.ID,
				name: item.name
			});
			this[arrName[1]].splice(index, 1);
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
		"editTemp.tempFieldList": {
			handler(list) {
				this.tempFieldList = [...list];
			},
			immediate: true
		},

		"editTemp.fieldList": {
			handler(list) {
				this.fieldList = [...list];
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
