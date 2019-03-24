$(document).ready(function () {
    $("input.js-videoTypeInput").change(function () {
        // TODO: ajax post to /ajax/output sending needed data

        console.log($(this).data("type"));
        console.log($(this).data("name"));
    });
});
