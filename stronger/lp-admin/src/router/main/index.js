import EntryRoutes from "./entry";
import HomeRoutes from "./home";
import LogRoutes from "./log";
import NoticesRoutes from "./notices";
import PDRoutes from "./pd";
import PMRoutes from "./pm";
import ReportRoutes from "./report";

const MainRoutes = {
	path: "/main",
	name: "Main",
	meta: {
		key: "main"
	},
	component: () => import("@/views/main"),

	children: [
		{
			path: "/main",
			redirect: "home"
		},

		EntryRoutes,
		HomeRoutes,
		LogRoutes,
		NoticesRoutes,
		PDRoutes,
		PMRoutes,
		ReportRoutes
	]
};

export default MainRoutes;
