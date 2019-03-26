function addToRowData(that, key, val) {
    var index = $(that).data("index");
    var data = $("#js-videoRow" + index).data();
    data[key] = val;
}

$(document).ready(function () {
    $("input.js-videoTypeInput").change(function () {
        addToRowData(this, "type", $(this).val());
        // TODO: ajax post to /ajax/output sending needed data
    });

    $("input.js-videoSkipCacheInput").change(function () {
        addToRowData(this, "skipcache", $(this).is(":checked"));
    });

    $("input.js-videoSkipOnlineSearchInput").change(function () {
        addToRowData(this, "skiponlinesearch", $(this).is(":checked"));
    });
});
