<template>
	<div class="task">
		<div class="z-flex align-center">
			<z-btn
				class="mr-3"
				color="primary"
				small
				unlocked
				@click="$refs.projects.onSelectAll()"
			>
				{{ projectValues.length === auth.proItems.length ? "全不选" : "全选" }}
			</z-btn>

			<z-checkboxs
				:formId="formId"
				formKey="items"
				ref="projects"
				:options="auth.proItems"
				:defaultValue="projectValues"
				@change="selectProjects"
			></z-checkboxs>
		</div>

		<v-divider></v-divider>

		<v-subheader class="pl-0">所选项目({{ projectValues.length }}个)：</v-subheader>

		<div v-for="(project, index) in projects" :key="project.proCode">
			<!-- 所选项目 BEGIN -->
			<task-details :proCode="project.proCode" :project="project"></task-details>
			<!-- 所选项目 END -->

			<!-- 分割线 BEGIN -->
			<v-divider v-show="index < projectValues.length - 1" class="mt-8 mb-4"></v-divider>
			<!-- 分割线 END -->
		</div>
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import { tools } from "vue-rocket";
import cells from "./cells";

export default {
	name: "Task",

	data() {
		return {
			formId: "Task",
			cells,
			list: [],
			projectValues: [],
			projects: []
		};
	},

	created() {
		this.getList();
	},

	methods: {
		async getList() {
			const result = await this.$store.dispatch("PM_TASK_GET_INPUT_CHANNEL_LIST");

			if (result.code === 200) {
				this.list = result.data?.list || [];
			}
		},

		selectProjects(values) {
			this.projects = [];

			values?.map(project => {
				const result = tools.find(this.list, { proCode: project });

				if (!result) {
					this.projects.push({ proCode: project });
				} else {
					this.projects.push(result);
				}
			});

			this.projectValues = values || [];
		}
	},

	computed: {
		...mapGetters(["auth"])
	},

	components: {
		"task-details": () => import("./taskDetails")
	}
};
</script>
