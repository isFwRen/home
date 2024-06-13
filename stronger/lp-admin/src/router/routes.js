import LoginRoutes from "./login";
import MainRoutes from "./main";
import UsageRoutes from "./usage";
import NormalRoutes from "./normal";

const router = [
	{
		path: "",
		redirect: "/login"
	},

	LoginRoutes,
	MainRoutes,
	UsageRoutes,
	NormalRoutes
];

export default router;
