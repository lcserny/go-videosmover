import {VideoWebResult} from "./VideoWebResult";

export class OutputResponseData {

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