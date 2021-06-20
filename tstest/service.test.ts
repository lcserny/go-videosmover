import {assert, expect} from "chai";
import {describe, it} from "mocha";
import {BasicVideoDataService} from "../tssrc/VideoDataService";

describe("service init", function () {
    it("new service is not null", function () {
        let service = new BasicVideoDataService(null);
        assert.notEqual(service, null);
    });
});