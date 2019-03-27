function addToRowData(index, key, val) {
    let data = $("#js-videoRow" + index).data();
    data[key] = val;
    return data;
}

$(document).ready(function () {
    $("input.js-videoTypeInput").change(function () {
        let rowIndex = $(this).data("index");
        let rowData = addToRowData(rowIndex, "type", $(this).val());
        $.post("/ajax/output", {data: JSON.stringify(rowData)}, function (response) {
            let $output = $("#videoOutput" + rowIndex);
            $output.val(response); // TODO: this is an array, get first always?
            $output.change();
        });
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
