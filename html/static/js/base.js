function addToRowData(index, key, val) {
    let data = $("#js-videoRow" + index).data();
    data[key] = val;
    return data;
}

$(document).ready(function () {
    $("input.js-videoTypeInput").change(function () {
        let index = $(this).data("index");
        let data = addToRowData(index, "type", $(this).val());

        $.ajax({
            url: 'ajax/output',
            type: 'post',
            dataType: 'html',
            data : { data: data},
            success : function(response) {
                let $output = $("#videoOutput" + index);
                $output.val(response);
                $output.change();

                console.log("Response is: " + response);
            },
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
