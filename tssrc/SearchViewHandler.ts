import {ModalHandler} from "./ModalHandler";
import {VideoDataService} from "./VideoDataService";
import {VideoData} from "./VideoData";
import {LoadingHelper} from "./helpers";
import {OutputResponseData} from "./OutputResponseData";
import {VideoWebResult} from "./VideoWebResult";
import {MoveResponseData} from "./MoveResponseData";

export class SearchViewHandler {

    private _modalHandler: ModalHandler;
    private _service: VideoDataService;

    constructor(service: VideoDataService, modalHandler: ModalHandler) {
        this._modalHandler = modalHandler;
        this._service = service;

        const videoRows = document.querySelectorAll<HTMLTableRowElement>('.js-videoRow');
        videoRows.forEach((row) => {
            let videoData = this.convertRowToVideoData(row);
            this._service.save(videoData)
        });
    }

    convertRowToVideoData(row: HTMLTableRowElement): VideoData {
        let initialData = JSON.parse(row.dataset.init);
        return new VideoData(
            initialData.Index,
            initialData.Name,
            initialData.FileName,
            initialData.VideoPath,
            initialData.Subtitles,
        );
    }

    findRow(rowChild: HTMLElement): HTMLTableRowElement {
        return rowChild.closest(".js-videoRow")
    }

    findIndex(rowChild: HTMLElement): number {
        try {
            return JSON.parse(this.findRow(rowChild).dataset.init).Index;
        } catch (e) {
            throw "Couldn't find parent video row of element " + rowChild;
        }
    }

    toggleClassOnElement(element: HTMLElement, show: boolean, className: string): void {
        if (show) {
            element.classList.add(className);
        } else {
            element.classList.remove(className);
        }
    }

    checkShowMoveVideosButton(): void {
        let button = document.querySelector<HTMLButtonElement>("#js-moveVideosButton");
        document.querySelector<HTMLSpanElement>("#js-moveVideosCount").innerText = "(" + this._service.getMovingVideosCount() + ")";
        button.style.display = this._service.shouldShowMoveButton() ? "initial" : "none";
    }

    checkShowGroupEditButton(): void {
        let button = document.querySelector<HTMLButtonElement>("#js-groupEditButton");
        document.querySelector<HTMLSpanElement>("#js-groupEditCount").innerText = "(" + this._service.getGroupedVideosCount() + ")";
        button.style.display = this._service.shouldShowGroupEditButton() ? "initial" : "none";
    }

    triggerChangeOutputTextBox(textBox: HTMLInputElement, value: string): void {
        textBox.value = value;
        textBox.dispatchEvent(new Event("keyup"));
    }

    handleVideoTypeChange(radio: HTMLInputElement): void {
        const index = this.findIndex(radio);
        const videoData = this._service.updateVideoType(index, radio.value);

        const row = this.findRow(radio);
        const outputTextBox = row.querySelector<HTMLInputElement>(".js-videoOutputInput");
        const outputDropdown = row.querySelector<HTMLDivElement>("#js-videoOutputDropdown" + index);

        if (radio.value === "unknown") {
            this.toggleClassOnElement(row, false, "highlight-row");
            this.checkShowMoveVideosButton();
            this.triggerChangeOutputTextBox(outputTextBox, "");
            return;
        }

        this.toggleClassOnElement(row, true, "highlight-row");
        LoadingHelper.showLoading();
        this._service.requestOutputDataAsync(videoData, false)
            .then(outData => {
                let converted = OutputResponseData.fromJson(outData);
                this._service.updateVideoOutList(index, converted.videos);
                this._service.updateVideoOutOrigin(index, converted.origin);
                this.triggerChangeOutputTextBox(outputTextBox, converted.videos[0].title);
                this.populateOutputDropdownList(outputDropdown, converted.videos);

                LoadingHelper.hideLoading();
                this.checkShowMoveVideosButton();
            }).catch(reason => {
            LoadingHelper.hideLoading();
            this.checkShowMoveVideosButton();
        });
    }

    handleVideoGroupTypeChange(radio: HTMLInputElement): void {
        const videoData = this._service.updateGroupVideoType(radio.value);

        const groupEditModal = document.querySelector<HTMLDivElement>("#js-groupEditModal");
        const outputTextBox = groupEditModal.querySelector<HTMLInputElement>(".js-videoOutputInput");
        const groupOutputDropdown = groupEditModal.querySelector<HTMLDivElement>("#js-videoGroupOutputDropdown");

        if (radio.value === "unknown") {
            this.triggerChangeOutputTextBox(outputTextBox, "");
            return;
        }

        LoadingHelper.showLoading();
        this._service.requestOutputDataAsync(videoData, false)
            .then(outData => {
                let converted = OutputResponseData.fromJson(outData);
                this._service.updateGroupVideoOutList(converted.videos);
                this._service.updateGroupVideoOutOrigin(converted.origin);
                this.triggerChangeOutputTextBox(outputTextBox, converted.videos[0].title);
                this.populateOutputDropdownList(groupOutputDropdown, converted.videos);

                LoadingHelper.hideLoading();
            }).catch(reason => {
            LoadingHelper.hideLoading();
        });
    }

    populateOutputDropdownList(dropdown: HTMLDivElement, outList: Array<VideoWebResult>): void {
        if (dropdown == null) {
            return;
        }
        let templateHtml = document.querySelector<HTMLScriptElement>("#js-videoOutputDropdown-item").innerHTML;
        let content = "";
        for (let vid of outList) {
            content += templateHtml.replace(/##title##/g, vid.title)
                .replace(/##posterURL##/g, vid.posterURL)
                .replace(/##description##/g, vid.description)
                .replace(/##cast##/g, vid.cast.join(", "));
        }
        dropdown.innerHTML = content;
    }

    triggerSearchVideosButton(): void {
        const searchVideosForm = document.querySelector<HTMLFormElement>("#js-searchVideosForm");
        searchVideosForm.submit();
    }

    showMoveIssuesModalWith(modalBody: string): void {
        document.querySelector<HTMLPreElement>("#js-moveIssuesModalBody").innerHTML = modalBody;
        this._modalHandler.showMoveIssuesModal();
    }

    showGroupEditModal(): void {
        this._service.saveVideoDataGroupingLeader();
        this._modalHandler.showGroupEditModal();
    }

    handleMoveVideosButtonClick(): void {
        LoadingHelper.showLoading();
        this._service.requestMoveVideosAsync()
            .then(response => {
                if (response.toString().length === 0) {
                    this.triggerSearchVideosButton();
                    return;
                }
                this.showMoveIssuesModalWith(JSON.stringify(MoveResponseData.fromJson(response), undefined, 2));
                LoadingHelper.hideLoading();
            });
    }

    handleGroupEditCheckBoxChange(row: HTMLTableRowElement, index: number, checked: boolean): void {
        this._service.updateVideoGrouping(index, checked);
        this.toggleClassOnElement(row, checked, "highlight-border");
        this.checkShowGroupEditButton();
    }

    handleGroupEditModalClose(): void {
        // apply leader to grouped videos
        const changedVideos = this._service.applyLeaderToGroupedVideos();
        for (let video of changedVideos) {
            const row = document.querySelector<HTMLTableRowElement>("#js-videoRow" + video.index);

            const multiEditCheckbox = row.querySelector<HTMLInputElement>(".js-videoMultiEdit");
            multiEditCheckbox.checked = false;
            this.toggleClassOnElement(row, false, "highlight-border");

            const videoTypeRadio = row.querySelector<HTMLInputElement>(".js-videoTypeInput[value='" + video.type + "']");
            videoTypeRadio.checked = true;
            this.toggleClassOnElement(row, video.type !== "unknown", "highlight-row");

            const skipCacheCheckbox = row.querySelector<HTMLInputElement>(".js-videoSkipCacheInput");
            skipCacheCheckbox.checked = video.skipCache;

            const skipOnlineCheckbox = row.querySelector<HTMLInputElement>(".js-videoSkipOnlineSearchInput");
            skipOnlineCheckbox.checked = video.skipOnline;

            const videoOutputTextBox = row.querySelector<HTMLInputElement>(".js-videoOutputInput");
            const onlineSearchButton = row.querySelector<HTMLButtonElement>(".js-videoOutputOnlineReSearch");
            this.toggleClassOnElement(onlineSearchButton, video.output !== "" && video.type !== "unknown", "show-element");
            videoOutputTextBox.value = video.output;

            const dropDown = row.querySelector<HTMLDivElement>("#js-videoOutputDropdown" + video.index);
            this.populateOutputDropdownList(dropDown, video.outputList);
        }

        // reset group UI and repo
        const groupEditModal = document.querySelector<HTMLDivElement>("#js-groupEditModal");
        const groupTypeRadios = groupEditModal.querySelectorAll<HTMLInputElement>(".js-videoGroupTypeInput");
        groupTypeRadios.forEach((radio) => {
            radio.checked = false;
        });

        const groupSkipCache = groupEditModal.querySelector<HTMLInputElement>(".js-videoGroupSkipCacheInput");
        groupSkipCache.checked = false;

        const groupSkipOnline = groupEditModal.querySelector<HTMLInputElement>(".js-videoGroupSkipOnlineSearchInput");
        groupSkipOnline.checked = false;

        const outputTextBox = groupEditModal.querySelector<HTMLInputElement>(".js-videoOutputInput");
        outputTextBox.value = "";

        const outputListPopup = groupEditModal.querySelector<HTMLDivElement>("#js-videoGroupOutputDropdown");
        outputListPopup.innerHTML = "";

        this.checkShowMoveVideosButton();
        this.checkShowGroupEditButton();

        this._service.resetGroupVideoLeader();
    }

    handleVideoOutputKeyup(row: HTMLTableRowElement, index: number, textbox: HTMLInputElement): void {
        const videoData = this._service.retrieve(index);
        const btn = row.querySelector<HTMLButtonElement>(".js-videoOutputOnlineReSearch");
        this.toggleClassOnElement(btn, textbox.value !== "" && videoData.type !== "unknown", "show-element");
        this._service.updateVideoOutput(index, textbox.value);
    }

    handleGroupVideoOutputKeyup(textbox: HTMLInputElement): void {
        const videoData = this._service.retrieveGroupVideoData();
        const btn = document.querySelector<HTMLButtonElement>("#js-videoGroupOutputOnlineReSearch");
        this.toggleClassOnElement(btn, textbox.value !== "" && videoData.type !== "unknown", "show-element");
        this._service.updateGroupVideoOutput(textbox.value);
    }

    handleOnlineReSearchButtonClick(row: HTMLTableRowElement, index: number): void {
        const videoData = this._service.retrieve(index);
        const outputTextBox = row.querySelector<HTMLInputElement>(".js-videoOutputInput");
        const outputDropdown = row.querySelector<HTMLDivElement>("#js-videoOutputDropdown" + index);

        LoadingHelper.showLoading();
        this._service.requestOutputDataAsync(videoData, true)
            .then(outData => {
                let converted = OutputResponseData.fromJson(outData);
                this._service.updateVideoOutList(index, converted.videos);
                this._service.updateVideoOutOrigin(index, converted.origin);
                this.triggerChangeOutputTextBox(outputTextBox, converted.videos[0].title);
                this.populateOutputDropdownList(outputDropdown, converted.videos);

                LoadingHelper.hideLoading();
                this.checkShowMoveVideosButton();
            }).catch(reason => {
            LoadingHelper.hideLoading();
            this.checkShowMoveVideosButton();
        });
    }

    handleGroupOnlineReSearchButtonClick(): void {
        const videoData = this._service.retrieveGroupVideoData();
        const groupEditModal = document.querySelector<HTMLDivElement>("#js-groupEditModal");
        const outputTextBox = groupEditModal.querySelector<HTMLInputElement>(".js-videoOutputInput");
        const groupOutputDropdown = groupEditModal.querySelector<HTMLDivElement>("#js-videoGroupOutputDropdown");

        LoadingHelper.showLoading();
        this._service.requestOutputDataAsync(videoData, true)
            .then(outData => {
                let converted = OutputResponseData.fromJson(outData);
                this._service.updateGroupVideoOutList(converted.videos);
                this._service.updateGroupVideoOutOrigin(converted.origin);
                this.triggerChangeOutputTextBox(outputTextBox, converted.videos[0].title);
                this.populateOutputDropdownList(groupOutputDropdown, converted.videos);
                LoadingHelper.hideLoading();
            }).catch(reason => {
            LoadingHelper.hideLoading();
        });
    }

    register(): void {
        // register event handlers on document
        const videoTypeRadios = document.querySelectorAll<HTMLInputElement>('.js-videoRow .js-videoTypeInput');
        videoTypeRadios.forEach((radio) => {
            radio.addEventListener('change', (event) => {
                this.handleVideoTypeChange(radio);
            });
        });

        const skipCacheCheckboxes = document.querySelectorAll<HTMLInputElement>('.js-videoRow .js-videoSkipCacheInput');
        skipCacheCheckboxes.forEach((checkBox) => {
            checkBox.addEventListener('change', (event) => {
                this._service.updateVideoSkipCache(this.findIndex(checkBox), checkBox.checked);
            });
        });

        const skipOnlineCheckboxes = document.querySelectorAll<HTMLInputElement>('.js-videoRow .js-videoSkipOnlineSearchInput');
        skipOnlineCheckboxes.forEach((checkBox) => {
            checkBox.addEventListener('change', (event) => {
                this._service.updateVideoSkipOnline(this.findIndex(checkBox), checkBox.checked);
            });
        });

        const outputTextBoxes = document.querySelectorAll<HTMLInputElement>('.js-videoRow .js-videoOutputInput');
        outputTextBoxes.forEach((textBox) => {
            textBox.addEventListener("keyup", (event) => {
                this.handleVideoOutputKeyup(this.findRow(textBox), this.findIndex(textBox), textBox);
            });
        });

        const onlineReSearchButtons = document.querySelectorAll<HTMLButtonElement>(".js-videoOutputOnlineReSearch");
        onlineReSearchButtons.forEach((btn) => {
            btn.addEventListener("click", (event) => {
                this.handleOnlineReSearchButtonClick(this.findRow(btn), this.findIndex(btn));
            });
        });

        const moveVideosButton = document.querySelector<HTMLButtonElement>("#js-moveVideosButton");
        moveVideosButton.addEventListener("click", (event) => {
            this.handleMoveVideosButtonClick();
        });

        // grouping listeners
        const groupEditButton = document.querySelector<HTMLButtonElement>("#js-groupEditButton");
        groupEditButton.addEventListener("click", (event) => {
            this.showGroupEditModal();
        });

        const groupEditCheckboxes = document.querySelectorAll<HTMLInputElement>(".js-videoRow .js-videoMultiEdit");
        groupEditCheckboxes.forEach((checkBox) => {
            checkBox.addEventListener("change", (event) => {
                this.handleGroupEditCheckBoxChange(this.findRow(checkBox), this.findIndex(checkBox), checkBox.checked);
            });
        });

        const groupEditNames = document.querySelectorAll<HTMLTableDataCellElement>(".js-videoRow .js-videoName");
        groupEditNames.forEach((name) => {
            name.addEventListener("click", (event) => {
                // TODO: if SHIFT is pressed
                // find current row
                // find first previous row already selected
                // change all rows between first and current
                // else: continue below

                const row = this.findRow(name);
                const multiEditCheckbox = row.querySelector<HTMLInputElement>(".js-videoMultiEdit");
                multiEditCheckbox.checked = !multiEditCheckbox.checked;
                this.handleGroupEditCheckBoxChange(row, this.findIndex(name), multiEditCheckbox.checked);
            });
        });

        const groupVideoTypeRadios = document.querySelectorAll<HTMLInputElement>('#js-groupEditModal .js-videoGroupTypeInput');
        groupVideoTypeRadios.forEach((radio) => {
            radio.addEventListener('change', (event) => {
                this.handleVideoGroupTypeChange(radio);
            });
        });

        const groupSkipCacheCheckbox = document.querySelector<HTMLInputElement>('#js-groupEditModal .js-videoGroupSkipCacheInput');
        groupSkipCacheCheckbox.addEventListener('change', (event) => {
            this._service.updateGroupVideoSkipCache(groupSkipCacheCheckbox.checked);
        });

        const groupSkipOnlineCheckbox = document.querySelector<HTMLInputElement>('#js-groupEditModal .js-videoGroupSkipOnlineSearchInput');
        groupSkipOnlineCheckbox.addEventListener('change', (event) => {
            this._service.updateGroupVideoSkipOnline(groupSkipOnlineCheckbox.checked);
        });

        const groupOutputTextBox = document.querySelector<HTMLInputElement>('#js-groupEditModal .js-videoOutputInput');
        groupOutputTextBox.addEventListener("keyup", (event) => {
            this.handleGroupVideoOutputKeyup(groupOutputTextBox);
        });

        const onlineGroupReSearchButton = document.querySelector<HTMLButtonElement>("#js-videoGroupOutputOnlineReSearch");
        onlineGroupReSearchButton.addEventListener("click", (event) => {
            this.handleGroupOnlineReSearchButtonClick();
        });

        // dynamic event handlers (elements that don't exist yet)
        const body = document.querySelector<HTMLBodyElement>("body");
        body.addEventListener("click", (event) => {
            const target = event.target as HTMLElement;
            const outputDropdownItem = target.closest<HTMLAnchorElement>(".js-output-dropdown-item");
            if (outputDropdownItem != null) {
                const textBox = outputDropdownItem.closest<HTMLDivElement>(".js-outputDropdownContainer").querySelector<HTMLInputElement>(".js-videoOutputInput");
                const title = outputDropdownItem.querySelector<HTMLDivElement>(".js-output-dropdown-item-title").innerText;
                this.triggerChangeOutputTextBox(textBox, title);
            }
        });

        this._modalHandler.register(this);
    }
}