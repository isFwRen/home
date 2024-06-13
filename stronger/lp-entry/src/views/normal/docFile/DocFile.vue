<template>
	<div class="wrap overflow-y-auto" ref="container" id="scroll-target" style="max-height: 100vh">
		<v-card
			style="padding-bottom: 50px; box-shadow: 0 0 0 0 !important; position: relative"
			class="scrollContainer"
			v-scroll:#scroll-target="onScroll"
		>
			<vue-office-pdf :src="path" @error="errorHandler" />
			<div class="btn_box">
				<div class="z-flex justify-center">
					<v-btn color="primary" @click="studyDone"> 完成学习 </v-btn>
				</div>
			</div>
		</v-card>
	</div>
</template>

<script>
import VueOfficePdf from "@vue-office/pdf";
import { mapState } from "vuex";
export default {
	name: "DocFile",
	data() {
		return {
			path: "",
			browseAll: false,
			dispatchList: "UPDATE_GUIDE_STAGE",
			broadListen:null
		};
	},
	created() {
		this.readDocFile();

		this.broadListen = new BroadcastChannel("updateRead");
	},
	computed: {
		...mapState(["train"])
	},
	methods: {
		onScroll(e) {
			// let scrollTop = e.target.scrollTop;
			// if (e.target.scrollHeight - scrollTop - this.$refs.container.clientHeight <= 1000 && !this.browseAll) {
			// 	this.browseAll = true;
			// 	this.$store.commit("SET_TRAINSTEPS", { key: 1, value: true });
			// 	console.log(this.train.trainSteps, this.browseAll);
			// }
		},
		async studyDone() {
			const res = await this.$store.dispatch(this.dispatchList, {});
			this.toasted.dynamic(res.msg, res.code);
			if(res.code === 200){
				this.broadListen.postMessage('updateRead');
				window.close()
			}
		},
		errorHandler() {
			console.log("渲染失败");
		},
		async readDocFile() {
			const res = await this.$store.dispatch("GET_DOC_FILE");
			this.path = res;
		}
	},
	destroy() {
		this.broadListen.close();
	},
	components: {
		VueOfficePdf
	}
};
</script>

<style lang="scss" scoped>
.wrap {
	width: 100%;
	position: relative;
	.btn_box {
		width: 100%;
		position: absolute;
		bottom: 5%;
	}
}
</style>
