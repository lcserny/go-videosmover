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
        this._moving = false;
        this._grouping = false;
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

    get moving() {
        return this._moving;
    }

    set moving(value) {
        this._moving = value;
    }

    get grouping() {
        return this._grouping;
    }

    set grouping(value) {
        this._grouping = value;
    }
}

// data layer
class InMemoryVideoDataRepository {
    constructor() {
        this.list = [];
        this.groupVideoData = null;
    }

    add(index, videoData) {
        if (videoData instanceof VideoData) {
            this.list[index] = videoData;
        }
    }

    addGroupVideoData(videoData) {
        if (videoData instanceof VideoData) {
            this.groupVideoData = videoData;
        }
    }

    clearGroupVideoData() {
        this.groupVideoData = null;
    }

    get(index) {
        let videoData = this.list[index];
        if (videoData instanceof VideoData) {
            return videoData;
        }
        return null;
    }

    getGroupVideoData() {
        return this.groupVideoData;
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
        videoData.moving = value !== "unknown";
        this.save(videoData);
        return videoData;
    }

    updateGroupVideoType(value) {
        let videoData = this.repo.getGroupVideoData();
        videoData.type = value;
        videoData.moving = value !== "unknown";
        this.saveGroupVideoData(videoData);
        return videoData;
    }

    updateVideoSkipCache(index, value) {
        let videoData = this.repo.get(index);
        videoData.skipCache = value;
        this.save(videoData);
        return videoData;
    }

    updateGroupVideoSkipCache(value) {
        let videoData = this.repo.getGroupVideoData();
        videoData.skipCache = value;
        this.saveGroupVideoData(videoData);
        return videoData;
    }

    updateVideoSkipOnline(index, value) {
        let videoData = this.repo.get(index);
        videoData.skipOnline = value;
        this.save(videoData);
        return videoData;
    }

    updateGroupVideoSkipOnline(value) {
        let videoData = this.repo.getGroupVideoData();
        videoData.skipOnline = value;
        this.saveGroupVideoData(videoData);
        return videoData;
    }

    updateVideoOutput(index, value) {
        let videoData = this.repo.get(index);
        videoData.output = value;
        this.save(videoData);
        return videoData;
    }

    updateGroupVideoOutput(value) {
        let videoData = this.repo.getGroupVideoData();
        videoData.output = value;
        this.saveGroupVideoData(videoData);
        return videoData;
    }

    updateVideoOutNames(index, values) {
        let videoData = this.repo.get(index);
        videoData.outputNames = values;
        this.save(videoData);
        return videoData;
    }

    updateGroupVideoOutNames(values) {
        let videoData = this.repo.getGroupVideoData();
        videoData.outputNames = values;
        this.saveGroupVideoData(videoData);
        return videoData;
    }

    updateVideoOutOrigin(index, value) {
        let videoData = this.repo.get(index);
        videoData.outputOrigin = value;
        this.save(videoData);
        return videoData;
    }

    updateGroupVideoOutOrigin(value) {
        let videoData = this.repo.getGroupVideoData();
        videoData.outputOrigin = value;
        this.saveGroupVideoData(videoData);
        return videoData;
    }

    updateVideoGrouping(index, value) {
        let videoData = this.repo.get(index);
        videoData.grouping = value;
        this.save(videoData);
        return videoData;
    }

    save(videoData) {
        this.repo.add(videoData.index, videoData);
    }

    saveGroupVideoData(videoData) {
        this.repo.addGroupVideoData(videoData);
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
        for (let videoData of this.repo.getAll()) {
            if (!videoData.moving) {
                continue;
            }

            moveDataList.push({
                video: videoData.path,
                subs: videoData.subs,
                type: videoData.type,
                outName: videoData.output
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
        return this.getMovingVideosCount() > 0;
    }

    shouldShowGroupEditButton() {
        return this.getGroupedVideosCount() > 0;
    }

    getGroupedVideosCount() {
        return this.getAllGroupedVideos().length;
    }

    getMovingVideosCount() {
        let count = 0;
        for (let videoData of this.repo.getAll()) {
            if (videoData.moving) {
                count++;
            }
        }
        return count;
    }

    getAllGroupedVideos() {
        const groupedVideoDataList = [];
        for (let videoData of this.repo.getAll()) {
            if (videoData.grouping) {
                groupedVideoDataList.push(videoData);
            }
        }
        return groupedVideoDataList;
    }

    saveVideoDataGroupingLeader() {
        for (let videoData of this.repo.getAll()) {
            if (videoData.grouping) {
                this.repo.addGroupVideoData(videoData);
                return videoData;
            }
        }
        return null;
    }

    resetGroupVideoLeader() {
        this.repo.clearGroupVideoData();
    }

    applyLeaderToGroupedVideos() {
        const leader = this.repo.getGroupVideoData();
        const groupedVideos = this.getAllGroupedVideos();

        for (let video of groupedVideos) {
            video.type = leader.type;
            video.skipCache = leader.skipCache;
            video.skipOnline = leader.skipOnline;
            video.output = leader.output;
            video.outputNames = leader.outputNames;
            video.outputOrigin = leader.outputOrigin;
            video.moving = leader.moving;
            video.grouping = false;
            this.save(video);
        }

        return groupedVideos;
    }
}

// UI layer
class SearchViewHandler {
    constructor(service, modalHandler) {
        this.modalHandler = modalHandler;
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

    toggleClassOnElement(element, show, className) {
        if (show) {
            element.classList.add(className);
        } else {
            element.classList.remove(className);
        }
    }

    checkShowMoveVideosButton() {
        let button = document.querySelector("#js-moveVideosButton");
        document.querySelector("#js-moveVideosCount").innerText = "(" + this.service.getMovingVideosCount() + ")";
        button.style.display = this.service.shouldShowMoveButton() ? "initial" : "none";
    }

    checkShowGroupEditButton() {
        let button = document.querySelector("#js-groupEditButton");
        document.querySelector("#js-groupEditCount").innerText = "(" + this.service.getGroupedVideosCount() + ")";
        button.style.display = this.service.shouldShowGroupEditButton() ? "initial" : "none";
    }

    triggerChangeOutputTextBox(textBox, value) {
        textBox.value = value;
        textBox.dispatchEvent(new Event("keyup"));
    }

    handleVideoTypeChange(radio, event) {
        const index = this.findIndex(radio);
        const videoData = this.service.updateVideoType(index, radio.value);

        const row = this.findRow(radio);
        const outputTextBox = row.querySelector(".js-videoOutputInput");
        const outputDropdown = row.querySelector("#js-videoOutputDropdown" + index);

        if (radio.value === "unknown") {
            this.toggleClassOnElement(row, false, "highlight-row");
            this.checkShowMoveVideosButton();
            this.triggerChangeOutputTextBox(outputTextBox, "");
            return;
        }

        this.toggleClassOnElement(row, true, "highlight-row");
        LoadingHelper.showLoading();
        this.service.requestOutputDataAsync(videoData)
            .then(outData => {
                this.service.updateVideoOutNames(index, outData["names"]);
                this.service.updateVideoOutOrigin(index, outData["origin"]);
                this.triggerChangeOutputTextBox(outputTextBox, outData["names"][0]);
                this.populateOutputDropdownList(outputDropdown, outData["names"]);
            })
            .finally(() => {
                LoadingHelper.hideLoading();
                this.checkShowMoveVideosButton();
            });
    }

    handleVideoGroupTypeChange(radio, event) {
        const videoData = this.service.updateGroupVideoType(radio.value);

        const groupEditModal = document.querySelector("#js-groupEditModal");
        const outputTextBox = groupEditModal.querySelector(".js-videoOutputInput");
        const groupOutputDropdown = groupEditModal.querySelector("#js-videoGroupOutputDropdown");

        if (radio.value === "unknown") {
            this.triggerChangeOutputTextBox(outputTextBox, "");
            return;
        }

        LoadingHelper.showLoading();
        this.service.requestOutputDataAsync(videoData)
            .then(outData => {
                this.service.updateGroupVideoOutNames(outData["names"]);
                this.service.updateGroupVideoOutOrigin(outData["origin"]);
                this.triggerChangeOutputTextBox(outputTextBox, outData["names"][0]);
                this.populateOutputDropdownList(groupOutputDropdown, outData["names"]);
            })
            .finally(() => {
                LoadingHelper.hideLoading();
            });
    }

    populateOutputDropdownList(dropdown, outNames) {
        if (dropdown == null) {
            return;
        }

        let templateHtml = document.querySelector("#js-videoOutputDropdown-item").innerHTML;
        let content = "";
        for (let name of outNames) {
            content += templateHtml.replace(/##outName##/g, name);
        }
        dropdown.innerHTML = content;
    }

    triggerSearchVideosButton() {
        const searchVideosForm = document.querySelector("#js-searchVideosForm");
        searchVideosForm.submit();
    }

    showMoveIssuesModalWith(modalBody) {
        document.querySelector("#js-moveIssuesModalBody").innerHTML = modalBody;
        this.modalHandler.showMoveIssuesModal();
    }

    showGroupEditModal() {
        this.service.saveVideoDataGroupingLeader();
        this.modalHandler.showGroupEditModal();
    }

    handleMoveVideosButtonClick(button, event) {
        LoadingHelper.showLoading();
        this.service.requestMoveVideosAsync()
            .then(response => {
                if (response.length === 0) {
                    this.triggerSearchVideosButton();
                    return;
                }

                this.showMoveIssuesModalWith(JSON.stringify(response, undefined, 2));
                LoadingHelper.hideLoading();
            });
    }

    handleGroupEditCheckBoxChange(row, index, checked) {
        this.service.updateVideoGrouping(index, checked);
        this.toggleClassOnElement(row, checked, "highlight-border");
        this.checkShowGroupEditButton();
    }

    handleGroupEditModalClose() {
        // apply leader to grouped videos
        const changedVideos = this.service.applyLeaderToGroupedVideos();
        for (let video of changedVideos) {
            const row = document.querySelector("#js-videoRow" + video.index);

            const multiEditCheckbox = row.querySelector(".js-videoMultiEdit");
            multiEditCheckbox.checked = false;
            this.toggleClassOnElement(row, false, "highlight-border");

            const videoTypeRadio = row.querySelector(".js-videoTypeInput[value='" + video.type + "']");
            videoTypeRadio.checked = true;
            this.toggleClassOnElement(row, video.type !== "unknown", "highlight-row");

            const skipCacheCheckbox = row.querySelector(".js-videoSkipCacheInput");
            skipCacheCheckbox.checked = video.skipCache;

            const skipOnlineCheckbox = row.querySelector(".js-videoSkipOnlineSearchInput");
            skipOnlineCheckbox.checked = video.skipOnline;

            const outputTextBox = row.querySelector(".js-videoOutputInput");
            const onlineSearchButton = row.querySelector(".js-videoOutputOnlineReSearch");
            this.toggleClassOnElement(onlineSearchButton, video.output !== "", "show-element");
            outputTextBox.value = video.output;

            const dropDown = row.querySelector("#js-videoOutputDropdown" + video.index);
            this.populateOutputDropdownList(dropDown, video.outputNames);
        }

        // reset group UI and repo
        const groupEditModal = document.querySelector("#js-groupEditModal");
        const groupTypeRadios = groupEditModal.querySelectorAll(".js-videoGroupTypeInput");
        for (let radio of groupTypeRadios) {
            radio.checked = false;
        }

        const groupSkipCache = groupEditModal.querySelector(".js-videoGroupSkipCacheInput");
        groupSkipCache.checked = false;

        const groupSkipOnline = groupEditModal.querySelector(".js-videoGroupSkipOnlineSearchInput");
        groupSkipOnline.checked = false;

        const outputTextBox = groupEditModal.querySelector(".js-videoOutputInput");
        outputTextBox.value = "";

        const outputNamesListPopup = groupEditModal.querySelector("#js-videoGroupOutputDropdown");
        outputNamesListPopup.innerHTML = "";

        this.checkShowMoveVideosButton();
        this.checkShowGroupEditButton();

        this.service.resetGroupVideoLeader();
    }

    handleVideoOutputKeyup(row, index, textbox) {
        const btn = row.querySelector(".js-videoOutputOnlineReSearch");
        this.toggleClassOnElement(btn, textbox.value !== "", "show-element");
        this.service.updateVideoOutput(index, textbox.value);
    }

    handleGroupVideoOutputKeyup(textbox) {
        const btn = document.querySelector("#js-videoGroupOutputOnlineReSearch");
        this.toggleClassOnElement(btn, textbox.value !== "", "show-element");
        this.service.updateGroupVideoOutput(textbox.value);
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
                this.handleVideoOutputKeyup(this.findRow(textBox), this.findIndex(textBox), textBox);
            });
        });

        // TODO: register trigger on button click re-search

        const moveVideosButton = document.querySelector("#js-moveVideosButton");
        moveVideosButton.addEventListener("click", (event) => {
            this.handleMoveVideosButtonClick(moveVideosButton, event);
        });

        // grouping listeners
        const groupEditButton = document.querySelector("#js-groupEditButton");
        groupEditButton.addEventListener("click", (event) => {
            this.showGroupEditModal();
        });

        const groupEditCheckboxes = document.querySelectorAll(".js-videoRow .js-videoMultiEdit");
        groupEditCheckboxes.forEach((checkBox) => {
            checkBox.addEventListener("change", (event) => {
                this.handleGroupEditCheckBoxChange(this.findRow(checkBox), this.findIndex(checkBox), checkBox.checked);
            });
        });

        const groupVideoTypeRadios = document.querySelectorAll('#js-groupEditModal .js-videoGroupTypeInput');
        groupVideoTypeRadios.forEach((radio) => {
            radio.addEventListener('change', (event) => {
                this.handleVideoGroupTypeChange(radio, event);
            });
        });

        const groupSkipCacheCheckbox = document.querySelector('#js-groupEditModal .js-videoGroupSkipCacheInput');
        groupSkipCacheCheckbox.addEventListener('change', (event) => {
            this.service.updateGroupVideoSkipCache(groupSkipCacheCheckbox.checked);
        });

        const groupSkipOnlineCheckbox = document.querySelector('#js-groupEditModal .js-videoGroupSkipOnlineSearchInput');
        groupSkipOnlineCheckbox.addEventListener('change', (event) => {
            this.service.updateGroupVideoSkipOnline(groupSkipOnlineCheckbox.checked);
        });

        const groupOutputTextBox = document.querySelector('#js-groupEditModal .js-videoOutputInput');
        groupOutputTextBox.addEventListener("keyup", (event) => {
            this.handleGroupVideoOutputKeyup(groupOutputTextBox);
        });

        // TODO: same for group

        // dynamic event handlers (elements that don't exist yet)
        const body = document.querySelector("body");
        body.addEventListener("click", (event) => {
            const element = event.target;
            if (element.tagName.toLowerCase() === "a" && element.classList.contains("js-dropdown-item")) {
                const textBox = element.closest(".js-outputDropdownContainer").querySelector(".js-videoOutputInput");
                this.triggerChangeOutputTextBox(textBox, element.innerText);
            }
        });

        this.modalHandler.register(this);
    }
}

// TODO: JQuery modals, try to remove later...
class JQueryModalHandler {
    constructor() {
        this.groupEditModal = $("#js-groupEditModal");
        this.moveIssuesModal = $("#js-moveIssuesModal");
    }

    showGroupEditModal() {
        this.groupEditModal.modal("show");
    }

    showMoveIssuesModal() {
        this.moveIssuesModal.modal("show");
    }

    register(viewHandler) {
        this.groupEditModal.on('hidden.bs.modal', function () {
            viewHandler.handleGroupEditModalClose();
        });

        this.moveIssuesModal.on('hidden.bs.modal', function () {
            viewHandler.triggerSearchVideosButton();
        });
    }
}

// init
$(document).ready(function () {
    new SearchViewHandler(new BasicVideoDataService(new InMemoryVideoDataRepository()), new JQueryModalHandler()).register();
});
