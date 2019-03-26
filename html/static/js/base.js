function addToRowData(index, key, val) {
    var data = $("#js-videoRow" + index).data();
    data[key] = val;
}

$(document).ready(function () {
    $("input.js-videoTypeInput").change(function () {
        addToRowData($(this).data("index"), "type", $(this).val());
        // TODO: ajax post to /ajax/output sending needed data then populate output field
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
