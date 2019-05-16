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

    updateVideoType(index, value) {
        // if value == unknown return
        // get from repo by index
        // if nothing in repo, create a videoData applying videoType and save to repo
        alert("index: " + index + ", value: " + value);
    }
}

class SearchHandler {
    constructor(service) {
        this.service = service;
        // loop over all rows/data if there and use service to convert rows to videoData and add to storage
    }

    findIndex(rowChild) {
        try {
            return rowChild.closest(".js-videoRow").querySelector(".js-videoRowIndex").innerHTML;
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
    }
}

$(document).ready(function () {
    new SearchHandler(new VideoDataService(new VideoDataRepository())).register();
});
