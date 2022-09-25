export class OutputRequestData {

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