function addToRowData(index, key, val) {
    let data = $("#js-videoRow" + index).data();
    data[key] = val;
    return data;
}

function populateOutputData(index, data) {
    let videoRowData = $("#js-videoRow" + index).data();
    let $output = $("#videoOutput" + index);
    let outputVal = data;

    if (outputVal === videoRowData["output"]) {
        return;
    }

    if (typeof data !== "string") {
        outputVal = data["names"][0];
        addToRowData(index, "outputnames", data["names"]);
        addToRowData(index, "outputorigin", data["origin"]);
        populateOutputDropdown(index, data["names"]);
    }
    addToRowData(index, "output", outputVal);

    $output.val(outputVal);
}

function populateOutputDropdown(index, outputNames) {
    let $dropdownContainer = $("#js-videoOutputDropdown" + index);
    let dropdownContent = "";
    for (let i = 0; i < outputNames.length; i++) {
        dropdownContent += "<a class=\"js-videoOutputDropdownItem dropdown-item\" data-index=\""
            + index + "\" href=\"#\">" + outputNames[i] + "</a>";
    }
    $dropdownContainer.html(dropdownContent);
}

function registerEventHandlers() {
    $(document).on("change", "input.js-videoSkipCacheInput", function () {
        addToRowData($(this).data("index"), "skipcache", $(this).is(":checked"));
    }).on("change", "input.js-videoSkipOnlineSearchInput", function () {
        addToRowData($(this).data("index"), "skiponlinesearch", $(this).is(":checked"));
    }).on("keyup", "input.js-videoOutputInput", function () {
        populateOutputData($(this).data("index"), $(this).val());
    }).on("click", "a.js-videoOutputDropdownItem", function () {
        populateOutputData($(this).data("index"), $(this).text());
    }).on("change", "input.js-videoTypeInput", function () {
        let rowIndex = $(this).data("index");
        let rowType = $(this).val();
        let rowData = addToRowData(rowIndex, "type", rowType);

        if (rowType === "unknown") {
            populateOutputData(rowIndex, "");
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
                $("#moveVideos").show();
                populateOutputData(rowIndex, response);
            });
        }
    }).on("hidden.bs.modal", "#moveIssuesModal", function () {
        $("#searchVideos").submit();
    }).on("click", "#moveVideos", function () {
        let dataList = [];
        $(".js-videoRow").each(function (i, row) {
            let rowData = $(row).data();
            let type = rowData["type"];
            if (type === "unknown") {
                return true;
            }

            let moveData = {
                video: rowData["path"],
                subs: rowData["subs"],
                type: type,
                outName: rowData["output"]
            };
            dataList.push(moveData);
        });

        $.post("/ajax/move", {movedata: JSON.stringify(dataList)}, function (response) {
            if (response.length === 0) {
                $("#searchVideos").submit();
                return;
            }
            $("#moveIssuesModal .modal-body").html(JSON.stringify(response));
            $("#moveIssuesModal").modal("show");
        });
    });
}

$(document).ready(function () {
    registerEventHandlers();
});
