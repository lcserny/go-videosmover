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

$(document).ready(function () {
    $("input.js-videoTypeInput").change(function () {
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
                populateOutputData(rowIndex, response);
            });
        }
    });

    $(document).on("change", "input.js-videoSkipCacheInput", function () {
        addToRowData($(this).data("index"), "skipcache", $(this).is(":checked"));
    }).on("change", "input.js-videoSkipOnlineSearchInput", function () {
        addToRowData($(this).data("index"), "skiponlinesearch", $(this).is(":checked"));
    }).on("keyup", "input.js-videoOutputInput", function () {
        populateOutputData($(this).data("index"), $(this).val());
    }).on("click", "a.js-videoOutputDropdownItem", function () {
        populateOutputData($(this).data("index"), $(this).text());
    });
});
