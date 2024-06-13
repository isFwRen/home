<template>
	<v-card>
		<v-toolbar class="z-toolbar" color="primary" dark>
			<z-btn icon dark @click="onBack">
				<v-icon>mdi-arrow-left</v-icon>
			</z-btn>

			<span class="pl-2" v-if="!$route.query.row">修改结果XML</span>
			<span class="pl-2" v-else>查看下载报文</span>

			<v-spacer></v-spacer>
			<span v-if="!$route.query.row">案件号：{{ cases.caseInfo.billNum }}</span>

			<span v-else>案件号：{{ $route.query.row.billNum }}</span>
			<v-tooltip bottom v-if="!$route.query.row">
				<template v-slot:activator="{ on, attrs }">
					<v-icon class="mr-1" color="white" v-bind="attrs" v-on="on" @click="exchange()">
						mdi-swap-horizontal
					</v-icon>
				</template>
				<span> 切换JSON格式 </span>
			</v-tooltip>
			<v-spacer></v-spacer>

			<z-btn icon dark @click="onClose">
				<v-icon>mdi-close</v-icon>
			</z-btn>
		</v-toolbar>

		<v-card-text class="pt-16">
			<div class="xml-content">
				<lp-monaco-editor
					v-if="visible"
					ref="editor"
					:value="xml"
					:format="format"
				></lp-monaco-editor>

				<div class="btns">
					<!-- <z-btn
            class="mr-3"
            color="primary"
            @click="onUpload"
          >
            <v-icon size="16">mdi-arrow-left-top-bold</v-icon>
            回传
          </z-btn> -->

					<z-btn
						class="mr-3"
						color="primary"
						@click="saveLocalXml"
						v-if="!$route.query.row"
					>
						<v-icon size="16">mdi-content-save-outline</v-icon>
						保存
					</z-btn>

					<z-btn color="primary" @click="onBack" v-if="!$route.query.row">
						<v-icon size="16">mdi-arrow-left</v-icon>
						返回
					</z-btn>
				</div>
			</div>
		</v-card-text>
	</v-card>
</template>

<script>
import moment from "moment";
import { mapGetters } from "vuex";
import DialogMixins from "@/mixins/DialogMixins";
import CaseMixins from "../CaseMixins";
import cells from "./cells";

export default {
	name: "UpdateResultXML",
	mixins: [DialogMixins, CaseMixins],

	data() {
		return {
			formId: "updateResultXML",
			cells,
			xml: "",
			format: "",
			flag: true,
			types: "xml",
			visible: false
		};
	},
	watch: {
		"cases.caseInfo": function (val) {
			console.log(val, "valval");
		}
	},
	created() {
		this.getOnlineXml();
	},

	watch: {
		"state.cases": function (val) {
			console.log(valm, "sss");
		}
	},

	computed: {
		...mapGetters(["cases"])
	},

	methods: {
		async getOnlineXml() {
			this.visible = false;
			const { createAt, caseId, proCode } = this.cases.caseInfo;
			// let row =  JSON.parse(window.sessionStorage.getItem('XMLDates'))

			const params1 = {
				year: moment(createAt).format("YYYY"),
				month: moment(createAt).format("M"),
				date: moment(createAt).format("D"),
				// caseId,
				billNum: this.cases.caseInfo.billNum,
				proCode: proCode,
				types: this.types
			};

			let otherInfo = "";
			if (this.$route.query.row) {
				otherInfo = this.$route.query.row.otherInfo;
			}

			if (this.$route.query.row) {
				this.xml = otherInfo;
				this.format = "json";
			} else {
				this.xml = await this.$store.dispatch("GET_CASE_XML", params1);
				if (this.flag) {
					this.format = "xml";
				} else {
					this.format = "json";
				}
			}
			const timeout = setTimeout(() => {
				this.visible = true;
				clearTimeout(timeout);
			}, 100);
		},

		// 保存
		async saveLocalXml() {
			const { createAt, billNum, proCode } = this.cases.caseInfo;

			this.xml = this.$refs.editor.getXML();

			const form = {
				year: moment(createAt).format("YYYY"),
				month: moment(createAt).format("M"),
				date: moment(createAt).format("D"),
				billNum,
				proCode,
				xml: this.xml,
				types: this.types
			};

			const result = await this.$store.dispatch("UPDATE_CASE_XML", form);

			this.toasted.dynamic(result.msg, result.code);
		},

		// 回传
		async onUpload() {
			const data = {
				proCode: this.cases.caseInfo.proCode,
				id: this.cases.caseInfo.caseId
			};

			const result = await this.$store.dispatch("UPLOAD_CASE_ITEM", data);
			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				this.getList();
			}
		},

		// 切换JSON格式
		async exchange() {
			this.flag = !this.flag;
			if (!this.flag) {
				this.types = "json";
			} else {
				this.types = "xml";
			}
			await this.getOnlineXml();
		}
	}
};
</script>

<style lang="scss">
.xml-content {
	height: calc(100vh - 84px);
	overflow: hidden;
	position: relative;

	.btns {
		position: absolute;
		right: 50px;
		bottom: 36px;
	}
}
.line-numbers{
	padding-left: 10px;
	width: auto !important;
}
</style>
