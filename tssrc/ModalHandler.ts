import $ = require('jquery')
import {SearchViewHandler} from "./SearchViewHandler";

export interface ModalHandler {
    showGroupEditModal(): void;
    showMoveIssuesModal(): void;
    register(viewHandler: SearchViewHandler): void;
}

// TODO: JQuery modals, try to remove later...
export class JQueryModalHandler implements ModalHandler{

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