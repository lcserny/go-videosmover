import LoadingHelper from "./base.js";

class VideoData {
    constructor(index) {
        this.index = index;
        this.name = "";
        this.fileName = "";
        this.path = "";
        this.subs = [];
        this.type = "unknown";
        this.skipCache = false;
        this.skipOnline = false;
        this.output = "";
        this.outputNames = [];
        this.outputOrigin = "";
    }

    setName(name) {
        this.name = name;
    }

    setFileName(fileName) {
        this.fileName = fileName;
    }

    setPath(path) {
        this.path = path;
    }

    setSubs(subs) {
        this.subs = subs;
    }

    setType(type) {
        this.type = type;
    }

    setSkipCache(skipCache) {
        this.skipCache = skipCache;
    }

    setSkipOnline(skipOnline) {
        this.skipOnline = skipOnline;
    }

    setOutput(output) {
        this.output = output;
    }

    setOutputNames(outNames) {
        this.outputNames = outNames;
    }

    setOutputOrigin(origin) {
        this.outputOrigin = origin;
    }

    getIndex() {
        return this.index;
    }

    getName() {
        return this.name;
    }

    getFileName() {
        return this.fileName;
    }

    getPath() {
        return this.path;
    }

    getSubs() {
        return this.subs;
    }

    getType() {
        return this.type;
    }

    getSkipCache() {
        return this.skipCache;
    }

    getSkipOnline() {
        return this.skipOnline;
    }

    getOutput() {
        return this.output;
    }

    getOutputNames() {
        return this.outputNames;
    }

    getOutputOrigin() {
        return this.outputOrigin;
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
        videoData.setType(value);
        this.save(videoData);
        return videoData;
    }

    updateVideoSkipCache(index, value) {
        let videoData = this.repo.get(index);
        videoData.setSkipCache(value);
        this.save(videoData);
        return videoData;
    }

    updateVideoSkipOnline(index, value) {
        let videoData = this.repo.get(index);
        videoData.setSkipOnline(value);
        this.save(videoData);
        return videoData;
    }

    updateVideoOutput(index, value) {
        let videoData = this.repo.get(index);
        videoData.setOutput(value);
        this.save(videoData);
        return videoData;
    }

    updateVideoOutNames(index, values) {
        let videoData = this.repo.get(index);
        videoData.setOutputNames(values);
        this.save(videoData);
        return videoData;
    }

    updateVideoOutOrigin(index, value) {
        let videoData = this.repo.get(index);
        videoData.setOutputOrigin(value);
        this.save(videoData);
        return videoData;
    }

    save(videoData) {
        this.repo.add(videoData.getIndex(), videoData);
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
                name: videoData.getName(),
                type: videoData.getType(),
                skipcache: videoData.getSkipCache(),
                skiponlinesearch: videoData.getSkipOnline(),
            })
        });
        return response.json();
    }

    async requestMoveVideosAsync() {
        let moveDataList = [];
        for (let videodata of this.repo.getAll()) {
            if (videodata.getType() === "unknown") {
                continue;
            }

            moveDataList.push({
                video: videodata.getPath(),
                subs: videodata.getSubs(),
                type: videodata.getType(),
                outName: videodata.getOutput()
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
            if (videoData.getType() === "movie" || videoData.getType() === "tv") {
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
        let videoData = new VideoData(initialData.Index);
        videoData.setName(initialData.Name);
        videoData.setFileName(initialData.FileName);
        videoData.setPath(initialData.VideoPath);
        videoData.setSubs(initialData.Subtitles);
        return videoData;
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
        if (this.service.shouldShowMoveButton()) {
            moveVideosButton.style.display = 'initial';
        } else {
            moveVideosButton.style.display = 'none';
        }
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
        } else {
            this.highlightRow(row, true);
            LoadingHelper.showLoading();
            this.service.requestOutputDataAsync(videoData)
                .then(outData => {
                    this.service.updateVideoOutNames(videoData.getIndex(), outData["names"]);
                    this.service.updateVideoOutOrigin(videoData.getIndex(), outData["origin"]);
                    this.triggerChangeOutputTextBox(outputTextBox, outData["names"][0]);
                })
                .finally(() => {
                    LoadingHelper.hideLoading();
                    this.checkShowMoveVideosButton();
                });
        }
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
