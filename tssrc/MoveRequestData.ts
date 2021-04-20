export class MoveRequestData {

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