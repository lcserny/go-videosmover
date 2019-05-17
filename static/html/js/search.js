import LoadingHelper from "./base.js";

class VideoData {
    constructor(index, name, fileName, path, subs) {
        this._index = index;
        this._name = name;
        this._fileName = fileName;
        this._path = path;
        this._subs = subs;
        this._type = "unknown";
        this._skipCache = false;
        this._skipOnline = false;
        this._output = "";
        this._outputNames = [];
        this._outputOrigin = "";
    }

    get index() {
        return this._index;
    }

    get name() {
        return this._name;
    }

    get fileName() {
        return this._fileName;
    }

    get path() {
        return this._path;
    }

    get subs() {
        return this._subs;
    }

    get type() {
        return this._type;
    }

    set type(value) {
        this._type = value;
    }

    get skipCache() {
        return this._skipCache;
    }

    set skipCache(value) {
        this._skipCache = value;
    }

    get skipOnline() {
        return this._skipOnline;
    }

    set skipOnline(value) {
        this._skipOnline = value;
    }

    get output() {
        return this._output;
    }

    set output(value) {
        this._output = value;
    }

    get outputNames() {
        return this._outputNames;
    }

    set outputNames(value) {
        this._outputNames = value;
    }

    get outputOrigin() {
        return this._outputOrigin;
    }

    set outputOrigin(value) {
        this._outputOrigin = value;
    }
}

// data layer
class InMemoryVideoDataRepository {
    constructor() {
        this.list = [];
    }

    add(index, videoData) {
        if (videoData instanceof VideoData) {
            this.list[index] = videoData;
        }
    }

    get(index) {
        let videoData = this.list[index];
        if (videoData instanceof VideoData) {
            return videoData;
        }
        return null;
    }

    getAll() {
        return this.list;
    }
}

// service layer
class BasicVideoDataService {
    constructor(repo) {
        this.repo = repo;
    }

    updateVideoType(index, value) {
        let videoData = this.repo.get(index);
        videoData.type = value;
        this.save(videoData);
        return videoData;
    }

    updateVideoSkipCache(index, value) {
        let videoData = this.repo.get(index);
        videoData.skipCache = value;
        this.save(videoData);
        return videoData;
    }

    updateVideoSkipOnline(index, value) {
        let videoData = this.repo.get(index);
        videoData.skipOnline = value;
        this.save(videoData);
        return videoData;
    }

    updateVideoOutput(index, value) {
        let videoData = this.repo.get(index);
        videoData.output = value;
        this.save(videoData);
        return videoData;
    }

    updateVideoOutNames(index, values) {
        let videoData = this.repo.get(index);
        videoData.outputNames = values;
        this.save(videoData);
        return videoData;
    }

    updateVideoOutOrigin(index, value) {
        let videoData = this.repo.get(index);
        videoData.outputOrigin = value;
        this.save(videoData);
        return videoData;
    }

    save(videoData) {
        this.repo.add(videoData.index, videoData);
    }

    // TODO: for debugging purposes
    retrieve(index) {
        return this.repo.get(index);
    }

    async requestOutputDataAsync(videoData) {
        const response = await fetch("/ajax/output", {
            method: 'POST',
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({
                name: videoData.name,
                type: videoData.type,
                skipcache: videoData.skipCache,
                skiponlinesearch: videoData.skipOnline,
            })
        });
        return response.json();
    }

    async requestMoveVideosAsync() {
        let moveDataList = [];
        for (let videodata of this.repo.getAll()) {
            if (videodata.type === "unknown") {
                continue;
            }

            moveDataList.push({
                video: videodata.path,
                subs: videodata.subs,
                type: videodata.type,
                outName: videodata.output
            });
        }

        const response = await fetch("/ajax/move", {
            method: 'POST',
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(moveDataList)
        });
        return response.json();
    }

    shouldShowMoveButton() {
        for (let videoData of this.repo.getAll()) {
            if (videoData.type === "movie" || videoData.type === "tv") {
                return true;
            }
        }
        return false;
    }
}

// UI layer
class SearchViewHandler {
    constructor(service) {
        this.service = service;

        const videoRows = document.querySelectorAll('.js-videoRow');
        videoRows.forEach((row) => {
            let videoData = this.convertRowToVideoData(row);
            this.service.save(videoData)
        });
    }

    convertRowToVideoData(row) {
        let initialData = JSON.parse(row.dataset.init);
        return new VideoData(
            initialData.Index,
            initialData.Name,
            initialData.FileName,
            initialData.VideoPath,
            initialData.Subtitles,
        );
    }

    findRow(rowChild) {
        return rowChild.closest(".js-videoRow")
    }

    findIndex(rowChild) {
        try {
            return JSON.parse(this.findRow(rowChild).dataset.init).Index;
        } catch (e) {
            throw "Couldn't find parent video row of element " + rowChild;
        }
    }

    highlightRow(row, show) {
        let className = "highlight-row";
        if (show) {
            row.classList.add(className);
        } else {
            row.classList.remove(className);
        }
    }

    checkShowMoveVideosButton() {
        let moveVideosButton = document.querySelector("#js-moveVideosButton");
        moveVideosButton.style.display = this.service.shouldShowMoveButton() ? "initial" : "none";
    }

    triggerChangeOutputTextBox(textBox, value) {
        textBox.value = value;
        textBox.dispatchEvent(new Event("keyup"));
    }

    handleVideoTypeChange(radio, event) {
        let videoData = this.service.updateVideoType(this.findIndex(radio), radio.value);

        let row = this.findRow(radio);
        let outputTextBox = row.querySelector(".js-videoOutputInput");
        if (radio.value === "unknown") {
            this.highlightRow(row, false);
            this.checkShowMoveVideosButton();
            this.triggerChangeOutputTextBox(outputTextBox, "");
            return;
        }

        this.highlightRow(row, true);
        LoadingHelper.showLoading();
        this.service.requestOutputDataAsync(videoData)
            .then(outData => {
                this.service.updateVideoOutNames(videoData.index, outData["names"]);
                this.service.updateVideoOutOrigin(videoData.index, outData["origin"]);
                this.triggerChangeOutputTextBox(outputTextBox, outData["names"][0]);
            })
            .finally(() => {
                LoadingHelper.hideLoading();
                this.checkShowMoveVideosButton();
            });
    }

    triggerSearchVideosButton() {
        const searchVideosForm = document.querySelector("#js-searchVideosForm");
        searchVideosForm.submit();
    }

    showJQueryMoveIssuesModalWith(modalBody) {
        const moveIssuesModalBody = document.querySelector("#js-moveIssuesModalBody");
        moveIssuesModalBody.innerHTML = modalBody;

        // TODO: remove JQuery...
        $("#js-moveIssuesModal").modal("show");
    }

    handleMoveVideosButtonClick(button, event) {
        LoadingHelper.showLoading();
        this.service.requestMoveVideosAsync()
            .then(response => {
                if (response.length === 0) {
                    this.triggerSearchVideosButton();
                    return;
                }

                this.showJQueryMoveIssuesModalWith(JSON.stringify(response, undefined, 2));
                LoadingHelper.hideLoading();
            });
    }

    register() {
        // register event handlers on document
        const videoTypeRadios = document.querySelectorAll('.js-videoRow .js-videoTypeInput');
        videoTypeRadios.forEach((radio) => {
            radio.addEventListener('change', (event) => {
                this.handleVideoTypeChange(radio, event);
            });
        });

        const skipCacheCheckboxes = document.querySelectorAll('.js-videoRow .js-videoSkipCacheInput');
        skipCacheCheckboxes.forEach((checkBox) => {
            checkBox.addEventListener('change', (event) => {
                this.service.updateVideoSkipCache(this.findIndex(checkBox), checkBox.checked);
            });
        });

        const skipOnlineCheckboxes = document.querySelectorAll('.js-videoRow .js-videoSkipOnlineSearchInput');
        skipOnlineCheckboxes.forEach((checkBox) => {
            checkBox.addEventListener('change', (event) => {
                this.service.updateVideoSkipOnline(this.findIndex(checkBox), checkBox.checked);
            });
        });

        const outputTextBoxes = document.querySelectorAll('.js-videoRow .js-videoOutputInput');
        outputTextBoxes.forEach((textBox) => {
            textBox.addEventListener("keyup", (event) => {
                this.service.updateVideoOutput(this.findIndex(textBox), textBox.value);
            });
        });

        const moveVideosButton = document.querySelector("#js-moveVideosButton");
        moveVideosButton.addEventListener("click", (event) => {
            this.handleMoveVideosButtonClick(moveVideosButton, event);
        });
    }
}

// init
$(document).ready(function () {
    new SearchViewHandler(new BasicVideoDataService(new InMemoryVideoDataRepository())).register();
});
