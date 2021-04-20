import {SearchViewHandler} from "./SearchViewHandler";
import {BasicVideoDataService} from "./VideoDataService";
import {InMemoryVideoDataRepository} from "./VideoDataRepository";
import {JQueryModalHandler} from "./ModalHandler";

$(function() {
    new SearchViewHandler(new BasicVideoDataService(new InMemoryVideoDataRepository()), new JQueryModalHandler()).register();
});