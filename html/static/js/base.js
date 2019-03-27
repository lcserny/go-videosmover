function addToRowData(index, key, val) {
    let data = $("#js-videoRow" + index).data();
    data[key] = val;
    return data;
}

function setOutputVal(index, val) {
    let $output = $("#videoOutput" + index);
    $output.val(val);
    $output.change();
}

$(document).ready(function () {
    $("input.js-videoTypeInput").change(function () {
        let rowIndex = $(this).data("index");
        let rowType = $(this).val();
        let rowData = addToRowData(rowIndex, "type", rowType);

        if (rowType === "unknown") {
            setOutputVal(rowIndex, "");
        } else {
            $.post("/ajax/output", {
                name: rowData["name"],
                type: rowData["type"],
                skipcache: rowData["skipcache"],
                skiponlinesearch: rowData["skiponlinesearch"],
            }, function (response) {
                if (typeof response === 'undefined' || response.length < 1) {
                    setOutputVal(rowIndex, "");
                    console.log("Output response invalid, check logs.");
                    return
                }

                setOutputVal(rowIndex, response["names"][0]); // TODO: this is an array, get first always?
                console.log(response); // TODO: remove me afterwards
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
        addToRowData($(this).data("index"), "output", $(this).val());
    });
});
