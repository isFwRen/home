import LoginRoutes from "./login";
import MainRoutes from "./main";
import UsageRoutes from "./usage";
import NormalRoutes from "./normal";
import transitRoutes from "./transit";
import DownloadRoutes from "./download";

const router = [
	{
		path: "",
		redirect: "/login"
	},
	DownloadRoutes,
	LoginRoutes,
	transitRoutes,
	MainRoutes,
	UsageRoutes,
	NormalRoutes
];

export default router;
