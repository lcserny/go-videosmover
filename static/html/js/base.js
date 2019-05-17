class RunHelper {
    setupPing(uri) {
        setInterval(function () {
            fetch(uri).catch(reason => alert("Webview failed to connect to server, reason: " + reason));
        }, 1000);
    }
}

class LoadingHelper {
    static showLoading() {
        let loadingContainer = document.querySelector(".js-loading-container");
        loadingContainer.style.display = 'block';
    }

    static hideLoading() {
        let loadingContainer = document.querySelector(".js-loading-container");
        loadingContainer.style.display = 'none';
    }
}

export default LoadingHelper;

$(document).ready(function () {
    new RunHelper().setupPing("/running");
});