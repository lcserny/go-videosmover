class RunHelper {
    setupPing(uri) {
        setInterval(function () {
            fetch(uri);
        }, 1000);
    }
}

$(document).ready(function () {
    new RunHelper().setupPing("/running");
});