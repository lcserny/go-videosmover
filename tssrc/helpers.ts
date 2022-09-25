export class LoadingHelper {
    static showLoading() {
        let loadingContainer = document.querySelector<HTMLDivElement>(".js-loading-container");
        loadingContainer.style.display = 'initial';
    }

    static hideLoading() {
        let loadingContainer = document.querySelector<HTMLDivElement>(".js-loading-container");
        loadingContainer.style.display = 'none';
    }
}

export class RunHelper {
    setupPing(uri: string, intervalMs: number) {
        setInterval(function () {
            fetch(uri).catch(reason => alert("Webview failed to connect to server, reason: " + reason));
        }, intervalMs);
    }
}