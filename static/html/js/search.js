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
}

class VideoDataRepository {
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

class VideoDataService {
    constructor(repo) {
        this.repo = repo;
    }

    updateVideoType(index, value) {
        let videoData = this.repo.get(index);
        videoData.setType(value);
        this.repo.add(index, videoData);

        if (value === "unknown") {
            return;
        }

        // TODO: ajax search and populate rest of fields
    }

    updateVideoSkipCache(index, value) {
        let videoData = this.repo.get(index);
        videoData.setSkipCache(value);
        this.repo.add(index, videoData);
    }

    updateVideoSkipOnline(index, value) {
        let videoData = this.repo.get(index);
        videoData.setSkipOnline(value);
        this.repo.add(index, videoData);
    }

    updateVideoOutput(index, value) {
        let videoData = this.repo.get(index);
        videoData.setOutput(value);
        this.repo.add(index, videoData);
    }

    save(videoData) {
        this.repo.add(videoData.getIndex(), videoData);
    }

    retrieve(index) {
        return this.repo.get(index);
    }
}

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

    findIndex(rowChild) {
        try {
            let row = rowChild.closest(".js-videoRow");
            return JSON.parse(row.dataset.init).Index;
        } catch (e) {
            throw "Couldn't find parent video row of element " + rowChild;
        }
    }

    register() {
        // register event handlers on document
        const videoTypeRadios = document.querySelectorAll('.js-videoRow .js-videoTypeInput');
        videoTypeRadios.forEach((radio) => {
            radio.addEventListener('change', (event) => {
                this.service.updateVideoType(this.findIndex(radio), radio.value);
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
            textBox.addEventListener('keyup', (event) => {
                this.service.updateVideoOutput(this.findIndex(textBox), textBox.value);
            });
        });
    }
}

$(document).ready(function () {
    new SearchHandler(new VideoDataService(new VideoDataRepository())).register();
});
