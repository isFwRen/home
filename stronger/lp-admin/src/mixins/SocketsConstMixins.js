import { localStorage } from "vue-rocket";
import io from "socket.io-client";
import { tools as lpTools } from "@/libs/util";

const isIntranet = lpTools.isIntranet();
const { baseURL, baseURLApi } = lpTools.baseURL();

export default {
	created() {
		// this.subscribeSockets();
	},

	methods: {
		subscribeSockets() {
			if (!this.socketPath) {
				return;
			}

			const userId = localStorage.get("user").ID;
			const room = "const";

			if (!userId) {
				this.toasted.dynamic(result.msg, 400);
			}

			const socket = io(`${baseURL}global-const`, {
				query: {
					room,
					userId
				},
				transports: ["websocket"]
			});

			socket.on("connect", () => {
				const id = socket.id;
				console.log("#connect,", id, socket);
			});

			socket.on(this.socketPath, result => {
				this.toasted.dynamic(result.msg, result.code);

				if (result.code === 200) {
					if (result.msg == "上传常量表成功!") {
						this.toasted.dynamic(result.msg, result.code);
					} else {
						if (isIntranet) {
							location.href = `${baseURL}${result.data}`;
						} else {
							location.href = `${baseURLApi}${result.data}`;
						}
					}
				}
			});
		}
	}
};
