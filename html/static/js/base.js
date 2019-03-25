$(document).ready(function () {
    $("input.js-videoTypeInput").change(function () {
        // TODO: ajax post to /ajax/output sending needed data
        var index = $(this).data("index");
        var $data = $("#js-videoRow" + index).data();
        $data["type"] = $(this).val();
        console.log($data);
    });
});
