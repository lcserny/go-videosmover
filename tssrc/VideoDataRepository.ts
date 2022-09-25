import {VideoData} from "./VideoData";

export interface VideoDataRepository {
    add(index: number, videoData: VideoData): void;
    addGroupVideoData(videoData: VideoData): void;
    clearGroupVideoData(): void;
    get(index: number): VideoData | null;
    getGroupVideoData(): VideoData | null;
    getAll(): Array<VideoData>;
}

export class InMemoryVideoDataRepository implements VideoDataRepository{

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