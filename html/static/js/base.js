function addToRowData(index, key, val) {
    let data = $("#js-videoRow" + index).data();
    data[key] = val;
    return data;
}

function handleOutputValChange(index, val) {
    let $output = $("#videoOutput" + index);
    if (val === "") {
        $output.val(val);
        $output.change();
        addToRowData(index, "outputnames", []);
        addToRowData(index, "outputorigin", "");
        return
    }

    $output.val(val["names"][0]);
    $output.change();
    addToRowData(index, "outputnames", val["names"]);
    addToRowData(index, "outputorigin", val["origin"]);
}

$(document).ready(function () {
    $("input.js-videoTypeInput").change(function () {
        let rowIndex = $(this).data("index");
        let rowType = $(this).val();
        let rowData = addToRowData(rowIndex, "type", rowType);

        if (rowType === "unknown") {
            handleOutputValChange(rowIndex, "");
        } else {
            $.post("/ajax/output", {
                name: rowData["name"],
                type: rowData["type"],
                skipcache: rowData["skipcache"],
                skiponlinesearch: rowData["skiponlinesearch"],
            }, function (response) {
                if (typeof response === 'undefined' || response.length < 1) {
                    handleOutputValChange(rowIndex, "");
                    console.log("Output response invalid, check logs.");
                    return
                }
                handleOutputValChange(rowIndex, response);
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
