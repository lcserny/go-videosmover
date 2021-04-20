export class MoveResponseData {

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