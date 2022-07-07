"use strict";

/*
 * Helper functions for the view-filter page form.
 */

function resetFilterParamRow(row) {
    row.getElementsByTagName('input')[0].value = ""
}

window.copyFilterOutput = function (clicked) {
    navigator.clipboard.writeText(document.getElementById("output-code").innerText);
    document.getElementById("output-card").classList.add("border-success")
    document.getElementById("output-copy-text").innerText = "Copied to clipboard"
    clicked.classList.add("btn-success")
}

window.deleteFilterParamRow = function (clicked) {
    const thisRow = clicked.closest(".input-group")
    if (thisRow.parentNode.childElementCount === 1) {
        resetFilterParamRow(thisRow) // Don't allow removing the last row
    } else {
        thisRow.remove()
    }
    htmx.trigger(htmx.find("#filter_input"), "input", {});
}

window.addFilterParamRow = function (clicked) {
    const thisRow = clicked.closest(".input-group")
    const clonedRow = thisRow.cloneNode(true)
    resetFilterParamRow(clonedRow)
    thisRow.parentNode.insertBefore(clonedRow, thisRow.nextSibling);
}
