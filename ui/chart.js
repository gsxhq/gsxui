// Chart bootstrap: the renderer body (chart.render.js, ~90KB — adapted
// from templui, see NOTICE.md) loads on the first chart in the DOM, so
// pages without charts never parse it.
import { init } from "./gsxui.js";

let bodyPromise;
// A failed dynamic import() (a transient network blip, an ad blocker, a
// flaky CDN in front of /ui/) used to cache the REJECTED promise forever:
// bodyPromise ??= ... never re-runs the import once it settles, rejected
// or not, so every chart on the page — including ones added later by a
// swap/morph, which init()'s own re-init contract re-invokes initRoot for
// — silently drew nothing for the rest of the page's life. The .catch
// clears the cache so the NEXT initRoot call (the next mutation, the next
// matched chart) gets a fresh import() attempt, and logs once per failed
// attempt so the failure is not silent.
//
// A cleared bodyPromise alone is not enough to make that next attempt a
// REAL retry, though: per the HTML module-map spec, a failed module fetch
// is cached by the browser keyed on the resolved URL, so a second
// import("./chart.render.js") — same specifier — resolves against that
// same cached failure and rejects again with no new network request ever
// made, cache clear or not. retryCount busts the URL (a query string the
// Go static file servers under /ui/ ignore for path routing, same as any
// static file server) so a retry is a genuinely fresh module-map entry.
let retryCount = 0;
function loadBody() {
	bodyPromise ??= import(retryCount === 0 ? "./chart.render.js" : `./chart.render.js?retry=${retryCount}`).catch(
		(err) => {
			console.error("gsxui: failed to load chart.render.js", err);
			bodyPromise = null;
			retryCount++;
			throw err;
		},
	);
	return bodyPromise;
}
function initRoot(el) {
	loadBody()
		.then(({ renderChart }) => renderChart(el))
		.catch(() => {}); // already logged in loadBody(); do not surface a second, per-caller unhandled rejection
}
init("[data-gsxui-slot-chart]", initRoot);
