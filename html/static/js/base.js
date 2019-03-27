function addToRowData(index, key, val) {
    let data = $("#js-videoRow" + index).data();
    data[key] = val;
    return data;
}

function handleOutputValChange(index, data, triggerChange) {
    let outputVal = data;
    let outputNames = [data];
    let outputOrigin = "";

    if (typeof data !== "string") {
        outputVal = data["names"][0];
        outputNames = data["names"];
        outputOrigin = data["origin"];
    }

    addToRowData(index, "output", outputVal);
    addToRowData(index, "outputnames", outputNames);
    addToRowData(index, "outputorigin", outputOrigin);

    if (triggerChange) {
        let $output = $("#videoOutput" + index);
        $output.val(outputVal);
        $output.change();
    }
}

$(document).ready(function () {
    $("input.js-videoTypeInput").change(function () {
        let rowIndex = $(this).data("index");
        let rowType = $(this).val();
        let rowData = addToRowData(rowIndex, "type", rowType);

        if (rowType === "unknown") {
            handleOutputValChange(rowIndex, "", true);
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
                handleOutputValChange(rowIndex, response, true);
            });
        }
    });

    $("input.js-videoSkipCacheInput").change(function () {
        addToRowData($(this).data("index"), "skipcache", $(this).is(":checked"));
    });

    $("input.js-videoSkipOnlineSearchInput").change(function () {
        addToRowData($(this).data("index"), "skiponlinesearch", $(this).is(":checked"));
    });

    $("input.js-videoOutputInput").change(function () {
        handleOutputValChange($(this).data("index"), $(this).val(), false);
    });
});
