import {VideoData} from "./VideoData";
import {VideoWebResult} from "./VideoWebResult";
import {VideoDataRepository} from "./VideoDataRepository";
import {OutputRequestData} from "./OutputRequestData";
import {MoveRequestData} from "./MoveRequestData";

export interface VideoDataService {
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

export class BasicVideoDataService implements VideoDataService {

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