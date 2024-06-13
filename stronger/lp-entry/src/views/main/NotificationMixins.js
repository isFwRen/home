import io from "socket.io-client";
import { tools as lpTools } from "@/libs/util";
import { localStorage, sessionStorage } from "vue-rocket";

export default {
	created() {
		let socketUrl = lpTools.isIntranet() ? process.env.VUE_APP_SOCKET_INNER_URL : process.env.VUE_APP_SOCKET_OUTER_URL;

		const procode = sessionStorage.get("proCode");
		if (procode === "B0108") {
			socketUrl = "https://www.i-confluence.com:41111";
		}
		if (procode === "B0114") {
			socketUrl = "https://www.i-confluence.com:51111";
		}
		const socket = io.connect(`${socketUrl}/global-notice`);
		socket.on("notice", result => {
			const proCode = sessionStorage.get("proCode");;
			if (result.proCode.includes(proCode)) {
				this.messages.push(`通知:${result.msg}`);
			}
		});
		this.$EventBus.$on("clearSocketMessage", () => {
			this.messages.length = 0;
		});
	},
	destroy() {
		this.$EventBus.$off("clearSocketMessage");
	}
};
