import { LoadingHelper} from "./base.js";

class VideoWebResult {

    private _title: string;
    private _description: string;
    private _posterURL: string;
    private _cast: Array<string>;

    constructor(title: string, description: string, posterURL: string, cast: Array<string>) {
        this._title = title;
        this._description = description;
        this._posterURL = posterURL;
        this._cast = cast;
    }

    get title(): string {
        return this._title;
    }

    set title(value: string) {
        this._title = value;
    }

    get description(): string {
        return this._description;
    }

    set description(value: string) {
        this._description = value;
    }

    get posterURL(): string {
        return this._posterURL;
    }

    set posterURL(value: string) {
        this._posterURL = value;
    }

    get cast(): Array<string> {
        return this._cast;
    }

    set cast(value: Array<string>) {
        this._cast = value;
    }
}

class MoveRequestData {

    public video: string;
    public subs: Array<string>;
    public type: string;
    public diskPath: string;
    public outName: string;

    constructor(video: string, subs: Array<string>, type: string, outName: string, diskPath: string = null) {
        this.video = video;
        this.subs = subs;
        this.type = type;
        this.outName = outName;
        this.diskPath = diskPath;
    }
}

class MoveResponseData {

    public unmovedFolder: string;
    public reasons: Array<string>;

    constructor(unmovedFolder: string, reasons: Array<string>) {
        this.unmovedFolder = unmovedFolder;
        this.reasons = reasons;
    }

    public static fromJson(json: any): MoveResponseData {
        return new MoveResponseData(json.unmovedFolder, json.reasons);
    }
}

class OutputRequestData {

    public name: string;
    public type: string;
    public skipCache: boolean;
    public skipOnlineSearch: boolean;
    public diskPath: string;

    constructor(name: string, type: string, skipCache: boolean, skipOnlineSearch: boolean, diskPath: string = "") {
        this.name = name;
        this.type = type;
        this.skipCache = skipCache;
        this.skipOnlineSearch = skipOnlineSearch;
        this.diskPath = diskPath;
    }
}

class OutputResponseData {

    public origin: string;
    public videos: Array<VideoWebResult>;

    constructor(origin: string, videos: Array<VideoWebResult>) {
        this.origin = origin;
        this.videos = videos;
    }

    public static fromJson(json: any): OutputResponseData {
        let videos = new Array<VideoWebResult>();
        if (json.hasOwnProperty("videos")) {
            for (let video of json.videos) {
                let cast = new Array<string>();
                if (video.Cast != null) {
                    for (let actr of video.Cast) {
                        cast.push(actr);
                    }
                }
                let webResult = new VideoWebResult(video.Title, video.Description, video.PosterURL, cast);
                videos.push(webResult)
            }
        }
        return new OutputResponseData(json.origin, videos);
    }
}

class VideoData {

    private _index: number;
    private _name: string;
    private _fileName: string;
    private _path: string;
    private _subs: Array<string>;
    private _type: string;
    private _skipCache: boolean;
    private _skipOnline: boolean;
    private _output: string;
    private _outputList: Array<VideoWebResult>;
    private _outputOrigin: string;
    private _moving: boolean;
    private _grouping: boolean;

    constructor(index: number, name: string, fileName: string, path: string, subs: Array<string>) {
        this._index = index;
        this._name = name;
        this._fileName = fileName;
        this._path = path;
        this._subs = subs;
        this._type = "unknown";
        this._skipCache = false;
        this._skipOnline = false;
        this._output = "";
        this._outputList = new Array<VideoWebResult>();
        this._outputOrigin = "";
        this._moving = false;
        this._grouping = false;
    }

    get index(): number {
        return this._index;
    }

    set index(value: number) {
        this._index = value;
    }

    get name(): string {
        return this._name;
    }

    set name(value: string) {
        this._name = value;
    }

    get fileName(): string {
        return this._fileName;
    }

    set fileName(value: string) {
        this._fileName = value;
    }

    get path(): string {
        return this._path;
    }

    set path(value: string) {
        this._path = value;
    }

    get subs(): Array<string> {
        return this._subs;
    }

    set subs(value: Array<string>) {
        this._subs = value;
    }

    get type(): string {
        return this._type;
    }

    set type(value: string) {
        this._type = value;
    }

    get skipCache(): boolean {
        return this._skipCache;
    }

    set skipCache(value: boolean) {
        this._skipCache = value;
    }

    get skipOnline(): boolean {
        return this._skipOnline;
    }

    set skipOnline(value: boolean) {
        this._skipOnline = value;
    }

    get output(): string {
        return this._output;
    }

    set output(value: string) {
        this._output = value;
    }

    get outputList(): Array<VideoWebResult> {
        return this._outputList;
    }

    set outputList(value: Array<VideoWebResult>) {
        this._outputList = value;
    }

    get outputOrigin(): string {
        return this._outputOrigin;
    }

    set outputOrigin(value: string) {
        this._outputOrigin = value;
    }

    get moving(): boolean {
        return this._moving;
    }

    set moving(value: boolean) {
        this._moving = value;
    }

    get grouping(): boolean {
        return this._grouping;
    }

    set grouping(value: boolean) {
        this._grouping = value;
    }
}

interface VideoDataRepository {
    add(index: number, videoData: VideoData): void;
    addGroupVideoData(videoData: VideoData): void;
    clearGroupVideoData(): void;
    get(index: number): VideoData | null;
    getGroupVideoData(): VideoData | null;
    getAll(): Array<VideoData>;
}

// data layer
class InMemoryVideoDataRepository implements VideoDataRepository{

    private readonly _list: Array<VideoData>;
    private _groupVideoData: VideoData;

    constructor() {
        this._list = new Array<VideoData>();
        this._groupVideoData = null;
    }

    add(index: number, videoData: VideoData): void {
        this._list[index] = videoData;
    }

    addGroupVideoData(videoData: VideoData): void {
        this._groupVideoData = videoData;
    }

    clearGroupVideoData(): void {
        this._groupVideoData = null;
    }

    get(index: number): VideoData | null {
        return this._list[index];
    }

    getGroupVideoData(): VideoData | null {
        return this._groupVideoData;
    }

    getAll(): Array<VideoData> {
        return this._list;
    }
}

interface VideoDataService {
    updateVideoType(index: number, value: string): VideoData;
    updateGroupVideoType(value: string): VideoData;
    updateVideoSkipCache(index: number, value: boolean): VideoData
    updateGroupVideoSkipCache(value: boolean): VideoData
    updateVideoSkipOnline(index: number, value: boolean): VideoData;
    updateGroupVideoSkipOnline(value: boolean): VideoData;
    updateVideoOutput(index: number, value: string): VideoData;
    updateGroupVideoOutput(value: string): VideoData;
    updateVideoOutList(index: number, values: Array<VideoWebResult>): VideoData;
    updateGroupVideoOutList(values: Array<VideoWebResult>): VideoData;
    updateVideoOutOrigin(index: number, value: string): VideoData;
    updateGroupVideoOutOrigin(value: string): VideoData;
    updateVideoGrouping(index: number, value: boolean): VideoData;
    save(videoData: VideoData): void;
    saveGroupVideoData(videoData: VideoData): void;
    retrieve(index: number): VideoData | null;
    retrieveGroupVideoData(): VideoData | null;
    requestOutputDataAsync(videoData: VideoData, useOutputInsteadOfName: boolean): Promise<any>;
    requestMoveVideosAsync(): Promise<any>;
    shouldShowMoveButton(): boolean;
    shouldShowGroupEditButton(): boolean;
    getGroupedVideosCount(): number;
    getMovingVideosCount(): number;
    getAllGroupedVideos(): Array<VideoData>;
    saveVideoDataGroupingLeader(): VideoData | null;
    resetGroupVideoLeader(): void;
    applyLeaderToGroupedVideos(): Array<VideoData>;
}

// service layer
class BasicVideoDataService implements VideoDataService {

    private _repo: VideoDataRepository;

    constructor(repo: VideoDataRepository) {
        this._repo = repo;
    }

    updateVideoType(index: number, value: string): VideoData {
        let videoData = this._repo.get(index);
        videoData.type = value;
        videoData.moving = value !== "unknown";
        this.save(videoData);
        return videoData;
    }

    updateGroupVideoType(value: string): VideoData {
        let videoData = this._repo.getGroupVideoData();
        videoData.type = value;
        videoData.moving = value !== "unknown";
        this.saveGroupVideoData(videoData);
        return videoData;
    }

    updateVideoSkipCache(index: number, value: boolean): VideoData {
        let videoData = this._repo.get(index);
        videoData.skipCache = value;
        this.save(videoData);
        return videoData;
    }

    updateGroupVideoSkipCache(value: boolean): VideoData {
        let videoData = this._repo.getGroupVideoData();
        videoData.skipCache = value;
        this.saveGroupVideoData(videoData);
        return videoData;
    }

    updateVideoSkipOnline(index: number, value: boolean): VideoData {
        let videoData = this._repo.get(index);
        videoData.skipOnline = value;
        this.save(videoData);
        return videoData;
    }

    updateGroupVideoSkipOnline(value: boolean): VideoData {
        let videoData = this._repo.getGroupVideoData();
        videoData.skipOnline = value;
        this.saveGroupVideoData(videoData);
        return videoData;
    }

    updateVideoOutput(index: number, value: string): VideoData {
        let videoData = this._repo.get(index);
        videoData.output = value;
        this.save(videoData);
        return videoData;
    }

    updateGroupVideoOutput(value: string): VideoData {
        let videoData = this._repo.getGroupVideoData();
        videoData.output = value;
        this.saveGroupVideoData(videoData);
        return videoData;
    }

    updateVideoOutList(index: number, values: Array<VideoWebResult>): VideoData {
        let videoData = this._repo.get(index);
        videoData.outputList = values;
        this.save(videoData);
        return videoData;
    }

    updateGroupVideoOutList(values: Array<VideoWebResult>): VideoData {
        let videoData = this._repo.getGroupVideoData();
        videoData.outputList = values;
        this.saveGroupVideoData(videoData);
        return videoData;
    }

    updateVideoOutOrigin(index: number, value: string): VideoData {
        let videoData = this._repo.get(index);
        videoData.outputOrigin = value;
        this.save(videoData);
        return videoData;
    }

    updateGroupVideoOutOrigin(value: string): VideoData {
        let videoData = this._repo.getGroupVideoData();
        videoData.outputOrigin = value;
        this.saveGroupVideoData(videoData);
        return videoData;
    }

    updateVideoGrouping(index: number, value: boolean): VideoData {
        let videoData = this._repo.get(index);
        videoData.grouping = value;
        this.save(videoData);
        return videoData;
    }

    save(videoData: VideoData): void {
        this._repo.add(videoData.index, videoData);
    }

    saveGroupVideoData(videoData: VideoData): void {
        this._repo.addGroupVideoData(videoData);
    }

    retrieve(index: number): VideoData | null {
        return this._repo.get(index);
    }

    retrieveGroupVideoData(): VideoData | null {
        return this._repo.getGroupVideoData();
    }

    async requestOutputDataAsync(videoData: VideoData, useOutputInsteadOfName: boolean): Promise<any> {
        const response = await fetch("/ajax/output", {
            method: 'POST',
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(new OutputRequestData(
                useOutputInsteadOfName ? videoData.output : videoData.name,
                videoData.type,
                videoData.skipCache,
                videoData.skipOnline
            ))
        });
        return response.json();
    }

    async requestMoveVideosAsync(): Promise<any> {
        let moveDataList = new Array<MoveRequestData>();
        for (let videoData of this._repo.getAll()) {
            if (!videoData.moving) {
                continue;
            }
            moveDataList.push(new MoveRequestData(videoData.path, videoData.subs, videoData.type, videoData.output));
        }
        const response = await fetch("/ajax/move", {
            method: 'POST',
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(moveDataList)
        });
        return response.json();
    }

    shouldShowMoveButton(): boolean {
        return this.getMovingVideosCount() > 0;
    }

    shouldShowGroupEditButton(): boolean {
        return this.getGroupedVideosCount() > 0;
    }

    getGroupedVideosCount(): number {
        return this.getAllGroupedVideos().length;
    }

    getMovingVideosCount(): number {
        let count = 0;
        for (let videoData of this._repo.getAll()) {
            if (videoData.moving) {
                count++;
            }
        }
        return count;
    }

    getAllGroupedVideos(): Array<VideoData> {
        const groupedVideoDataList = new Array<VideoData>();
        for (let videoData of this._repo.getAll()) {
            if (videoData.grouping) {
                groupedVideoDataList.push(videoData);
            }
        }
        return groupedVideoDataList;
    }

    saveVideoDataGroupingLeader(): VideoData | null {
        for (let videoData of this._repo.getAll()) {
            if (videoData.grouping) {
                this._repo.addGroupVideoData(videoData);
                return videoData;
            }
        }
        return null;
    }

    resetGroupVideoLeader(): void {
        this._repo.clearGroupVideoData();
    }

    applyLeaderToGroupedVideos(): Array<VideoData> {
        const leader = this._repo.getGroupVideoData();
        const groupedVideos = this.getAllGroupedVideos();

        for (let video of groupedVideos) {
            video.type = leader.type;
            video.skipCache = leader.skipCache;
            video.skipOnline = leader.skipOnline;
            video.output = leader.output;
            video.outputList = leader.outputList;
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
        for (let radio of groupTypeRadios) {
            radio.checked = false;
        }

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

interface ModalHandler {
    showGroupEditModal(): void;
    showMoveIssuesModal(): void;
    register(viewHandler: SearchViewHandler): void;
}

// TODO: JQuery modals, try to remove later...
class JQueryModalHandler implements ModalHandler{

    private _groupEditModal: JQuery<HTMLDivElement>;
    private _moveIssuesModal: JQuery<HTMLDivElement>;

    constructor() {
        this._groupEditModal = $("#js-groupEditModal");
        this._moveIssuesModal = $("#js-moveIssuesModal");
    }

    showGroupEditModal(): void {
        this._groupEditModal.modal("show");
    }

    showMoveIssuesModal(): void {
        this._moveIssuesModal.modal("show");
    }

    register(viewHandler: SearchViewHandler): void {
        this._groupEditModal.on('hidden.bs.modal', function () {
            viewHandler.handleGroupEditModalClose();
        });

        this._moveIssuesModal.on('hidden.bs.modal', function () {
            viewHandler.triggerSearchVideosButton();
        });
    }
}

// init
$(function() {
    new SearchViewHandler(new BasicVideoDataService(new InMemoryVideoDataRepository()), new JQueryModalHandler()).register();
});