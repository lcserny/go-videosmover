class RunHelper {
    setupPing(uri, intervalMs) {
        setInterval(function () {
            fetch(uri).catch(reason => alert("Webview failed to connect to server, reason: " + reason));
        }, intervalMs);
    }
}

class LoadingHelper {
    static showLoading() {
        let loadingContainer = document.querySelector(".js-loading-container");
        loadingContainer.style.display = 'initial';
    }

    static hideLoading() {
        let loadingContainer = document.querySelector(".js-loading-container");
        loadingContainer.style.display = 'none';
    }
}

export default LoadingHelper;

$(document).ready(function () {
    new RunHelper().setupPing("/running", 1000);
});