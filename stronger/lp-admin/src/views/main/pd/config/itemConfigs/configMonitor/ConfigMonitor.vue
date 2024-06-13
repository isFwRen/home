<template>
	<div class="config-monitor">
		<div class="pt-2 pb-2">项目下载监控配置（针对采用FTP进行数据传输的项目）</div>
		<div class="table_wrap">
			<vxe-table :data="desserts" :edit-config="{ trigger: 'dblclick', mode: 'cell' }">
				<vxe-column type="seq" title="序号" width="60"></vxe-column>
				<template v-for="item in cells.headers">
					<vxe-column
						v-if="item.value === 'frequency'"
						:min-width="item.width"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default="{ row }">
							<vxe-select v-model="row.frequency" style="width: 100px">
								<vxe-option
									v-for="item in cells.frequencyOptions"
									:key="item.value"
									:value="item.value"
									:label="item.label"
								></vxe-option>
							</vxe-select>
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'CreatedAt'"
						:min-width="item.width"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default="{ row }">
							{{ row.CreatedAt | dateFormat("YYYY-MM-DD") }}
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'UpdatedAt'"
						:min-width="item.width"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default="{ row }">
							{{ row.UpdatedAt | dateFormat("YYYY-MM-DD") }}
						</template>
					</vxe-column>

					<vxe-column
						v-else
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
					</vxe-column>
				</template>
			</vxe-table>
		</div>

		<div class="download-explain">
			<h6 class="explain_title">下载说明</h6>
			<div class="content">
				<textarea
					ref="textareaRef"
					@dblclick="onDblclick"
					@blur="onBlur"
					class="textarea_desc"
					cols="30"
					rows="10"
					v-model="desc"
				></textarea>
			</div>
		</div>
	</div>
</template>

<script>
import { mapGetters, mapState } from "vuex";
import TableMixins from "@/mixins/TableMixins";
import ConfigMixins from "@/views/main/pd/config/ConfigMixins";
import cells from "./cells";

export default {
	name: "ConfigMonitor",
	mixins: [TableMixins, ConfigMixins],

	data() {
		return {
			showBorder: false,
			desc: "this is a desc",
			dispatchList: "GET_FTP_MONITOR",
			dispatchListEdit: "EDIT_FTP_MONITOR",
			cells
		};
	},

	mounted() {
		this.$refs.textareaRef.addEventListener("mousedown", function (event) {
			event.preventDefault();
		});
	},
	watch: {
		desserts(val) {
			this.desc = val[0].desc;
		}
	},
	methods: {
		onMouseLeave() {
			console.log("onMouseLeave");
		},
		onDblclick() {
			this.$refs.textareaRef.focus();
			this.showBorder = true;
		},
		async onBlur() {
			this.showBorder = false;
			const data = {
				createdCode: this.desserts[0].createdCode,
				createdName: this.desserts[0].createdName,
				desc: this.desc,
				frequency: this.desserts[0].frequency,
				id: this.desserts[0].ID,
				proCode: this.desserts[0].proCode,
				wrongMsg: this.desserts[0].wrongMsg
			};
			console.log(data, "data");
			const result = await this.$store.dispatch(this.dispatchListEdit, data);
			if (result.code === 200) {
				this.onSearch();
			}
		}
	},
	computed: {
		...mapGetters(["config"])
	},
	unmounted() {
		this.$refs.textareaRef.removeEventListener("mousedown", function (event) {});
	}
};
</script>
<style lang="scss">
.config-monitor {
	min-height: 100vh;
}

.download-explain {
	.explain_title {
		font-size: 16px;
		padding-bottom: 10px;
	}
}
.content {
	padding: 2px 5px;
	color: #3c3c43c7;
	.textarea_desc {
		border: 1px solid #eee;
		transition: 0.2s;
		resize: none;
		outline: none;
		width: 100%;
		padding: 10px;
		border-radius: 10px;
	}
}

.table_wrap {
	.vxe-table--body-wrapper {
		overflow: visible;
	}
}
</style>
