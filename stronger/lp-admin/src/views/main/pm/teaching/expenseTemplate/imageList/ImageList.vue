<template>
	<div class="image-list">
		<div
			v-for="(image, index) in imageList"
			:class="['image', size]"
			:key="`${image.model.ID}-${Date.now()}`"
			@dblclick="handleViewImage"
		>
			<v-card class="mx-auto" @contextmenu="onShowMenu($event, { ...image, index })">
				<!-- <v-img class="white--text align-end" height="200px" :src="image.base64"> </v-img> -->
				<div style="height: 200px">
					<img
						class="white--text align-end"
						:src="image.base64"
						style="width: 100%; height: 100%; object-fit: contain"
					/>
				</div>

				<div class="pt-4 px-4">
					<v-text-field
						:ref="`input_${index}`"
						dense
						:flat="image.flat"
						:readonly="image.readonly"
						solo
						:value="image.name"
						@keydown="onEnter($event, image)"
					></v-text-field>
				</div>
			</v-card>

			<div v-show="showSelect" class="select">
				<v-checkbox v-model="image.selected" @change="onSelect($event, image)"></v-checkbox>
			</div>
		</div>

		<v-menu
			v-model="showMenu"
			absolute
			class="elevation-0"
			offset-y
			:position-x="x"
			:position-y="y"
		>
			<v-list>
				<v-list-item
					v-for="(item, index) in cells.rightClickOptions"
					:key="index"
					@click="onMenuItem(item)"
				>
					<v-list-item-title>{{ item.label }}</v-list-item-title>
				</v-list-item>
			</v-list>
		</v-menu>

		<attr-dialog ref="attrDialog" :rowInfo="detailInfo"></attr-dialog>
	</div>
</template>

<script>
import { sessionStorage, tools } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";
import cells from "../cells";

const { baseURLApi } = lpTools.baseURL();

export default {
	name: "TeachingExpenseTemplateImageList",

	props: {
		desserts: {
			type: Array,
			default: () => []
		},

		size: {
			validator(value) {
				return ["large", "medium", "small"].includes(value);
			},
			default: "large"
		},

		showSelect: {
			type: Boolean,
			default: false
		}
	},

	data() {
		return {
			cells,
			x: 0,
			y: 0,
			showMenu: false,
			imageList: [],
			imageIndex: 0,
			choiceItems: [],
			baseURLApi,
			detailInfo: {}
		};
	},

	methods: {
		handleViewImage() {
			const thumbs = [];

			this.desserts.map(dessert => {
				thumbs.push({
					thumbPath: `${baseURLApi}${dessert.path}`,
					path: `${baseURLApi}${dessert.path}`
				});
			});

			sessionStorage.set("thumbs", thumbs);

			window.open(
				`${location.origin}/normal/view-images`,
				"_blank",
				"toolbar=yes, scrollbars=yes, resizable=yes"
			);
		},

		onShowMenu(event, info) {
			event.preventDefault();
			this.showMenu = false;
			this.x = event.clientX;
			this.y = event.clientY;

			this.$nextTick(() => {
				this.showMenu = true;
			});

			this.detailInfo = info;
			this.imageIndex = info.index;
		},

		onMenuItem(info) {
			switch (info.value) {
				case "rename":
					this.imageList.map((image, imageIndex) => {
						image.readonly = true;
						image.flat = true;

						if (imageIndex === this.imageIndex) {
							image.readonly = false;
							image.flat = false;
						}
					});

					this.$refs[`input_${this.imageIndex}`][0].focus();

					break;

				case "delete":
					this.$modal({
						visible: true,
						title: "删除提示",
						content: `请确认是否要删除？`,
						confirm: async () => {
							const body = {
								proCode: this.detailInfo.proCode,
								ids: [this.detailInfo.ID]
							};

							const result = await this.$store.dispatch(
								"DELETE_PM_TEACHING_EXPENSE_TEMPLATE_ITEM",
								body
							);

							this.toasted.dynamic(result.msg, result.code);

							if (result.code === 200) {
								this.$emit("deleted");
							}
						}
					});
					break;

				case "attr":
					this.$refs.attrDialog.onOpen(0);
					break;
			}
		},

		async onEnter(event, info) {
			if (event.keyCode !== 13) {
				return;
			}

			const value = event.target.value;

			if (!value) {
				this.toasted.warning("不能为空命名!");
				this.$emit("renamed");
				return;
			}

			this.imageList.map(image => {
				image.autofocus = false;
				image.flat = true;
			});

			const body = {
				proCode: info.proCode,
				id: info.model.ID,
				name: value
			};

			const result = await this.$store.dispatch(
				"RENAME_PM_TEACHING_EXPENSE_TEMPLATE_ITEM",
				body
			);

			this.toasted.dynamic(result.msg, result.code);

			this.$emit("renamed");
		},

		onSelect(value, { ID: id }) {
			if (value) {
				this.choiceItems.push(id);
			} else {
				const index = tools.findIndex(this.choiceItems, id);
				this.choiceItems.splice(index, 1);
			}

			this.$emit("select", this.choiceItems);
		}
	},

	watch: {
		desserts: {
			handler(desserts) {
				this.imageList = [];

				desserts.map(dessert => {
					this.imageList = [
						...this.imageList,
						{ ...dessert, base64: "", readonly: true, flat: true }
					];
				});
				this.imageList.forEach(async ele => {
					const newBase64 = await lpTools.getTokenImg(`${baseURLApi}${ele.path}`);
					if (newBase64) {
						lpTools.getBase64(newBase64).then(base64String => {
							ele.base64 = base64String;
						});
					}
				});
			},
			immediate: true
		}
	},

	components: {
		"attr-dialog": () => import("./attrDialog")
	}
};
</script>

<style scoped lang="scss">
.image-list {
	display: flex;
	flex-wrap: wrap;

	.image {
		position: relative;
		padding: 12px;
		box-sizing: border-box;

		&.large {
			width: 33.33333333%;
		}

		&.medium {
			width: 25%;
		}

		&.small {
			width: 20%;
		}

		img {
			width: 100%;
		}

		.select {
			position: absolute;
			top: 0;
			left: 20px;
		}
	}
}
</style>
