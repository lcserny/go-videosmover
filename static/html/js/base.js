class RunHelper {
    setupPing(uri) {
        setInterval(function () {
            fetch(uri).catch(reason => alert("Webview failed to connect to server, reason: " + reason));
        }, 1000);
    }
}

$(document).ready(function () {
    new RunHelper().setupPing("/running");
});