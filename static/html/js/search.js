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
            this.list.splice(index, 0, videoData);
        }
    }

    get(index) {
        let videoData = this.list[index];
        if (videoData instanceof VideoData) {
            return videoData;
        }
        return null;
    }
}

// service layer
class VideoDataService {
    constructor(repo) {
        this.repo = repo;
    }

    updateVideoType(index, value) {
        let videoData = this.repo.get(index);
        videoData.setType(value);
        this.repo.add(index, videoData);
        return videoData;
    }

    updateVideoSkipCache(index, value) {
        let videoData = this.repo.get(index);
        videoData.setSkipCache(value);
        this.repo.add(index, videoData);
        return videoData;
    }

    updateVideoSkipOnline(index, value) {
        let videoData = this.repo.get(index);
        videoData.setSkipOnline(value);
        this.repo.add(index, videoData);
        return videoData;
    }

    updateVideoOutput(index, value) {
        let videoData = this.repo.get(index);
        videoData.setOutput(value);
        this.repo.add(index, videoData);
        return videoData;
    }

    updateVideoOutNames(index, values) {
        let videoData = this.repo.get(index);
        videoData.setOutputNames(values);
        this.repo.add(index, videoData);
        return videoData;
    }

    updateVideoOutOrigin(index, value) {
        let videoData = this.repo.get(index);
        videoData.setOutputOrigin(value);
        this.repo.add(index, videoData);
        return videoData;
    }

    save(videoData) {
        this.repo.add(videoData.getIndex(), videoData);
    }

    retrieve(index) {
        return this.repo.get(index);
    }

    async requestOutputDataAsync(videoData) {
        const jsonResp = await fetch("/ajax/output", {
            method: 'POST',
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({
                name: videoData.getName(),
                type: videoData.getType(),
                skipcache: videoData.getSkipCache(),
                skiponlinesearch: videoData.getSkipOnline(),
            })
        });
        return jsonResp.json();
    }
}

// UI layer
class SearchHandler {
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

    handleVideoTypeChange(radio, event) {
        let videoData = this.service.updateVideoType(this.findIndex(radio), radio.value);

        let row = this.findRow(radio);
        let outputTextBox = row.querySelector(".js-videoOutputInput");
        if (radio.value === "unknown") {
            this.highlightRow(row, false);
            outputTextBox.value = "";
        } else {
            this.highlightRow(row, true);
            LoadingHelper.showLoading();
            this.service.requestOutputDataAsync(videoData)
                .then(outData => {
                    this.service.updateVideoOutNames(videoData.getIndex(), outData["names"]);
                    this.service.updateVideoOutOrigin(videoData.getIndex(), outData["origin"]);
                    outputTextBox.value = outData["names"][0];
                })
                .finally(() => LoadingHelper.hideLoading());
        }
        outputTextBox.dispatchEvent(new Event("keyup"));
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
    }
}

// init
$(document).ready(function () {
    new SearchHandler(new VideoDataService(new InMemoryVideoDataRepository())).register();
});
