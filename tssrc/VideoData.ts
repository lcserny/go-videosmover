import {VideoWebResult} from "./VideoWebResult";

export class VideoData {

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