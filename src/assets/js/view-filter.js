"use strict";

/*
 * Helper functions for the view-filter page form.
 */

import * as htmx from "htmx.org";

function resetFilterParamRow(row) {
    row.getElementsByTagName('input')[0].value = ""
}

function focusFilterParamRow(row) {
    row.getElementsByTagName('input')[0].focus()
}

window.copyFilterOutput = function (clicked) {
    navigator.clipboard.writeText(document.getElementById("output-code").innerText);
    document.getElementById("output-card").classList.add("border-success")
    document.getElementById("output-copy-text").innerText = "Copied to clipboard"
    clicked.classList.add("btn-success")
}

window.deleteFilterParamRow = function (clicked) {
    const thisRow = clicked.closest(".input-group")
    const rowCount = thisRow.parentElement.getElementsByClassName("input-group").length

    if (rowCount > 1) {
        thisRow.remove()
    } else { // Don't allow removing the last row
        resetFilterParamRow(thisRow)
        focusFilterParamRow(thisRow)
    }
    htmx.trigger(htmx.find("#filter_input"), "input", {});
}

window.addFilterParamRow = function (clicked) {
    const thisRow = clicked.closest(".input-group")
    const clonedRow = thisRow.cloneNode(true)
    resetFilterParamRow(clonedRow)
    thisRow.parentNode.insertBefore(clonedRow, thisRow.nextSibling);
    focusFilterParamRow(clonedRow)
}
