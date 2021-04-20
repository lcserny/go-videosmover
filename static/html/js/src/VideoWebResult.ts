export class VideoWebResult {

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