class VideoData {
    constructor(name, value) {
        this.name = name;
        this.value = value;
    }

    getName() {
        return this.name;
    }

    getValue() {
        return this.value;
    }
}

class VideoDataRepository {
    constructor() {
        // init in memory list storage
        this.list = [];
    }

    add(videoData) {
        if (videoData instanceof VideoData) {
            this.list.push(videoData);
        }
    }

    get(index) {
        let videoData = this.list.get(index);
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
}

class SearchHandler {
    constructor(service) {
        this.service = service;
        // loop over all rows/data if there and use service to convert rows to videoData and add to storage
    }

    register() {
        // register event handlers on document
        const videoTypeRadios = document.querySelectorAll('.js-videoRow .js-videoTypeInput');
        videoTypeRadios.forEach((radio) => {
            radio.addEventListener('change', (event) => {
                alert('Thanks for clicking! ' + radio.value);
            });
        });
    }
}

$(document).ready(function () {
    new SearchHandler(new VideoDataService(new VideoDataRepository())).register();
});
