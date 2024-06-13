import Vue from "vue";

import LPCalendar from "./lp-calendar";
import LPDialog from "./lp-dialog";
import LPDropdown from "./lp-dropdown";
import LPMonacoEditor from "./lp-monaco-editor";
import LPNotification from "./lp-notification";
import LPTabs from "./lp-tabs";
import LPSpinners from "./lp-spinners";
import LPTooltipBtn from "./lp-tooltip-btn";
import LPMessageDialog from "./lp-message-dialog";
import LPDrawingBoard from "./lp-drawing-board";
import LPLoading from "./lp-loading";

const components = () => {
	Vue.component("lp-calendar", LPCalendar);
	Vue.component("lp-dialog", LPDialog);
	Vue.component("lp-dropdown", LPDropdown);
	Vue.component("lp-monaco-editor", LPMonacoEditor);
	Vue.component("lp-notification", LPNotification);
	Vue.component("lp-tabs", LPTabs);
	Vue.component("lp-spinners", LPSpinners);
	Vue.component("lp-tooltip-btn", LPTooltipBtn);
	Vue.component("lp-message-dialog", LPMessageDialog);
	Vue.component("lp-drawing-board", LPDrawingBoard);
	Vue.component("lp-loading", LPLoading);
};

Vue.use(components);
