// TODO: refactor functions

function addToRowData(index, key, val) {
    let data = $("#js-videoRow" + index).data();
    data[key] = val;
    return data;
}

function populateOutputData(index, data, itemClass) {
    let videoRowData = $("#js-videoRow" + index).data();
    let outputVal = data;

    if (outputVal === videoRowData["output"]) {
        return;
    }

    if (typeof data !== "string") {
        outputVal = data["names"][0];
        addToRowData(index, "outputnames", data["names"]);
        addToRowData(index, "outputorigin", data["origin"]);
        populateOutputDropdown(index, data["names"], itemClass);
    }
    addToRowData(index, "output", outputVal);

    $("#videoOutput" + index).val(outputVal);
}

function populateOutputDropdown(index, outputNames, itemClass) {
    let $dropdownContainer = $("#js-videoOutputDropdown" + index);
    let dropdownContent = "";
    for (let i = 0; i < outputNames.length; i++) {
        dropdownContent += "<a class=\"" + itemClass + " dropdown-item\" data-index=\""
            + index + "\" href=\"#\">" + outputNames[i] + "</a>";
    }
    $dropdownContainer.html(dropdownContent);
}

function handleVideoTypeChange(index, type, itemClass, $loadingContainer) {
    let data = addToRowData(index, "type", type);
    let $row = $("tr#js-videoRow" + index);

    if (type === "unknown") {
        populateOutputData(index, "", itemClass);
        $row.removeClass("highlight-row")
    } else {
        $loadingContainer.show();
        $.post("/ajax/output", {
            name: data["name"],
            type: data["type"],
            skipcache: data["skipcache"],
            skiponlinesearch: data["skiponlinesearch"],
        }, function (response) {
            $loadingContainer.hide();
            if (typeof response === 'undefined' || response.length < 1) {
                response = "";
                console.log("Output response invalid, check logs.");
                $row.removeClass("highlight-row")
            }
            $("#moveVideos").show();
            populateOutputData(index, response, itemClass);
            $row.addClass("highlight-row")
        });
    }
}

function getListOfCheckmarkedVideoMultiEdit() {
    let list = [];
    $("input.js-videoMultiEdit").each(function () {
        if ($(this).prop("checked")) {
            list.push($(this));
        }
    });
    return list;
}

function capitalizeFirstLetter(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
}

// TODO: move each handler better in functions
function registerEventHandlers() {
    let $loadingContainer = $("#loading-container");

    $(document).on("change", "input.js-videoSkipCacheInput", function () {
        addToRowData($(this).data("index"), "skipcache", $(this).is(":checked"));
    }).on("change", "input.js-videoGroupSkipCacheInput", function () {
        let groupEditList = getListOfCheckmarkedVideoMultiEdit();
        for (let i = 0; i < groupEditList.length; i++) {
            let $ge = groupEditList[i];
            let rowIndex = $ge.data("index");
            addToRowData(rowIndex, "skipcache", $(this).is(":checked"));
            $("#videoSkipCache" + rowIndex).prop("checked", $(this).prop("checked")).trigger("change");
        }
        // TODO: populate also this popup input.js-videoGroupSkipCacheInput
    }).on("change", "input.js-videoSkipOnlineSearchInput", function () {
        addToRowData($(this).data("index"), "skiponlinesearch", $(this).is(":checked"));
    }).on("change", "input.js-videoGroupSkipOnlineSearchInput", function () {
        let groupEditList = getListOfCheckmarkedVideoMultiEdit();
        for (let i = 0; i < groupEditList.length; i++) {
            let $ge = groupEditList[i];
            let rowIndex = $ge.data("index");
            addToRowData(rowIndex, "skiponlinesearch", $(this).is(":checked"));
            $("#videoSkipOnlineSearch" + rowIndex).prop("checked", $(this).prop("checked")).trigger("change");
        }
        // TODO: populate also this popup input.js-videoGroupSkipOnlineSearchInput
    }).on("keyup", "input.js-videoOutputInput", function () {
        populateOutputData($(this).data("index"), $(this).val(), "js-videoOutputDropdownItem");
    }).on("keyup", "input.js-videoGroupOutputInput", function () {
        let groupEditList = getListOfCheckmarkedVideoMultiEdit();
        for (let i = 0; i < groupEditList.length; i++) {
            let $ge = groupEditList[i];
            let rowIndex = $ge.data("index");
            populateOutputData(rowIndex, $(this).val(), "js-videoGroupOutputDropdownItem");
        }
        // TODO: populate also this popup input.js-videoGroupOutputInput
    }).on("click", "a.js-videoOutputDropdownItem", function (event) {
        populateOutputData($(this).data("index"), $(this).text(), "js-videoOutputDropdownItem");
        event.preventDefault();
    }).on("click", "a.js-videoGroupOutputDropdownItem", function (event) {
        let groupEditList = getListOfCheckmarkedVideoMultiEdit();
        for (let i = 0; i < groupEditList.length; i++) {
            let $ge = groupEditList[i];
            let rowIndex = $ge.data("index");
            populateOutputData(rowIndex, $(this).text(), "js-videoGroupOutputDropdownItem");
        }
        // TODO: populate also this popup input.js-videoGroupOutputDropdownItem
        event.preventDefault();
    }).on("click", "button#groupEdit", function () {
        $("#groupEditModal").modal("show");
    }).on("click", "input.js-videoMultiEdit", function () {
        let groupEditListSize = getListOfCheckmarkedVideoMultiEdit().length;
        if (groupEditListSize > 0) {
            $("#groupEditCount").text("(" + groupEditListSize + ")");
        }
        let $groupEditButton = $("#groupEdit");
        if (groupEditListSize > 0) {
            $groupEditButton.show();
        } else {
            $groupEditButton.hide();
        }

        let rowIndex = $(this).data("index");
        let $row = $("tr#js-videoRow" + rowIndex);
        if ($(this).prop("checked")) {
            $row.addClass("highlight-border")
        } else {
            $row.removeClass("highlight-border")
        }
    }).on("change", "input.js-videoTypeInput", function () {
        handleVideoTypeChange($(this).data("index"), $(this).val(), "js-videoOutputDropdownItem", $loadingContainer);
    }).on("change", "input.js-videoGroupTypeInput", function () {
        let groupEditList = getListOfCheckmarkedVideoMultiEdit();
        for (let i = 0; i < groupEditList.length; i++) {
            let $ge = groupEditList[i];
            let rowIndex = $ge.data("index");
            handleVideoTypeChange(rowIndex, $(this).val(), "js-videoGroupOutputDropdownItem", $loadingContainer);
            $("#videoType" + capitalizeFirstLetter($(this).val()) + rowIndex).prop("checked", $(this).prop("checked")).trigger("change");
        }
    }).on("hidden.bs.modal", "#moveIssuesModal", function () {
        $("#searchVideos").submit();
    }).on("hidden.bs.modal", "#groupEditModal", function () {
        let groupEditList = getListOfCheckmarkedVideoMultiEdit();
        for (let i = 0; i < groupEditList.length; i++) {
            let $ge = groupEditList[i];
            $ge.trigger("click");
        }

        $(".js-videoGroupTypeInput").prop("checked", false).trigger("change");
        $(".js-videoGroupSkipCacheInput").prop("checked", false).trigger("change");
        $(".js-videoGroupSkipOnlineSearchInput").prop("checked", false).trigger("change");
        $(".js-videoGroupOutputInput").val("").trigger("change");
        $(".js-videoGroupOutputDropdown").html("").trigger("change");
    }).on("click", "#moveVideos", function () {
        $loadingContainer.show();

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
            $loadingContainer.hide();

            if (response.length === 0) {
                $("#searchVideos").submit();
                return;
            }
            $("#moveIssuesModal .modal-body pre").html(JSON.stringify(response, undefined, 2));
            $("#moveIssuesModal").modal("show");
        });
    });
}

$(document).ready(function () {
    registerEventHandlers();
});
