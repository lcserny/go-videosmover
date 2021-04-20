import {RunHelper} from "./helpers";

$(function() {
    new RunHelper().setupPing("/running", 1000);
});