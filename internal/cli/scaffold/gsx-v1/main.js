// gsx dev panel: Cmd-D / Ctrl-D
import "virtual:gsx-devpanel";
import "./style.css";
import { setupCounter } from "./counter.js";

setupCounter(document.querySelector("#counter"));
