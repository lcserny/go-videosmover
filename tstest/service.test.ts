import {assert, expect} from "chai";
import {describe, it} from "mocha";
import {BasicVideoDataService} from "../tssrc/VideoDataService";
import {InMemoryVideoDataRepository} from "../tssrc/VideoDataRepository";

describe("service init", function () {
    it("new service is not null", function () {
        let repo = new InMemoryVideoDataRepository();
        let service = new BasicVideoDataService(repo);
        assert.notEqual(service, null);
    });
});