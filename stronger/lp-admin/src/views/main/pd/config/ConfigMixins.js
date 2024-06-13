import storage from "@/libs/util.storage";

const defaultConfig = {
	proId: "",
	tempId: "",
	chunkId: ""
};

export default {
	data() {
		return {
			proId: ""
		};
	},

	methods: {
		rememberIds(config) {
			let storageConfig = storage.get("config") || defaultConfig;

			if (config.proId !== undefined) {
				this.proId = config.proId;
			}

			storageConfig = Object.assign(storageConfig, config);

			storage.set("config", storageConfig);
		}
	},

	watch: {
		$route: {
			handler() {
				const storageConfig = this.storage.get("config");
				const storageProject = this.storage.get("project");

				if (!storageConfig) {
					this.storage.set("config", defaultConfig);
				} else {
					this.proId = storageConfig.proId;

					this.effectParams = {
						...this.effectParams,
						proId: this.proId,
						proCode: storageProject.code
					};
				}
			},
			immediate: true
		}
	}
};
