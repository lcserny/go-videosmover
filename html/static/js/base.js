function addToRowData(index, key, val) {
    let data = $("#js-videoRow" + index).data();
    data[key] = val;
    return data;
}

function populateOutputData(index, data) {
    let $output = $("#videoOutput" + index);
    let outputVal = data;
    let outputNames = [data];
    let outputOrigin = "MANUAL";

    if (typeof data !== "string") {
        outputVal = data["names"][0];
        outputNames = data["names"];
        outputOrigin = data["origin"];
    }

    $output.val(outputVal);
    addToRowData(index, "output", outputVal);
    addToRowData(index, "outputnames", outputNames);
    addToRowData(index, "outputorigin", outputOrigin);
}

$(document).ready(function () {
    $("input.js-videoTypeInput").change(function () {
        let rowIndex = $(this).data("index");
        let rowType = $(this).val();
        let rowData = addToRowData(rowIndex, "type", rowType);

        if (rowType === "unknown") {
            populateOutputData(rowIndex, "");
        } else {
            $.post("/ajax/output", {
                name: rowData["name"],
                type: rowData["type"],
                skipcache: rowData["skipcache"],
                skiponlinesearch: rowData["skiponlinesearch"],
            }, function (response) {
                if (typeof response === 'undefined' || response.length < 1) {
                    response = "";
                    console.log("Output response invalid, check logs.");
                }
                populateOutputData(rowIndex, response);
            });
        }
    });

    $("input.js-videoSkipCacheInput").change(function () {
        addToRowData($(this).data("index"), "skipcache", $(this).is(":checked"));
    });

    $("input.js-videoSkipOnlineSearchInput").change(function () {
        addToRowData($(this).data("index"), "skiponlinesearch", $(this).is(":checked"));
    });

    $("input.js-videoOutputInput").keyup(function () {
        populateOutputData($(this).data("index"), $(this).val());
    });
});
