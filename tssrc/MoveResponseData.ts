export class MoveResponseData {

    public unmovedFolder: string;
    public reasons: Array<string>;

    constructor(unmovedFolder: string, reasons: Array<string>) {
        this.unmovedFolder = unmovedFolder;
        this.reasons = reasons;
    }
}