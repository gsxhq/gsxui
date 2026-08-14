// Chart bootstrap: the renderer body (chart.render.js, ~90KB — adapted
// from templui, see NOTICE.md) loads on the first chart in the DOM, so
// pages without charts never parse it.
import { init } from "./gsxui.js";

let bodyPromise;
function initRoot(el) {
	bodyPromise ??= import("./chart.render.js");
	bodyPromise.then(({ renderChart }) => renderChart(el));
}
init("[data-gsxui-slot-chart]", initRoot);
