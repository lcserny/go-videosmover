$(document).ready(function () {
    $("input.js-videoTypeInput").change(function () {
        // TODO: ajax post to /ajax/output sending needed data
        var $data = $(this).parents("tr");
        var type = $(this).val();
        console.log($data.data());
        console.log(type);
    });
});
