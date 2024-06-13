import { mapGetters } from "vuex";

export default {
	methods: {
		async getRoleSysMenu() {
			const result = await this.$store.dispatch("GET_ROLE_SYS_MENU");

			if (result.code === 200) {
				const { Perm } = result.data;
				const [mapPro, proItems] = this.resolvePerm(Perm);
				const res = await this.$store.dispatch("GET_MENUS");
				if (res.code === 200) {
					this.$store.commit("UPDATE_AUTH", {
						menus: res.data.menus,
						perm: Perm,
						mapPro,
						proItems
					});
				}
			}
		},

		resolvePerm(permissions) {
			const mapPro = {};
			const proItems = [];

			permissions.map(permission => {
				if (permission.hasPm) {
					mapPro[permission.proCode] = permission;
					proItems.push({ label: permission.proCode, value: permission.proCode });
				}
			});

			return [mapPro, proItems];
		}
	},

	computed: {
		...mapGetters(["auth"])
	}
};
