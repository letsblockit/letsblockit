function copyOutput(clicked) {
    navigator.clipboard.writeText(document.getElementById("output-code").innerText);
    document.getElementById("output-card").classList.add("border-success")
    document.getElementById("output-copy-text").innerText = "Copied to clipboard"
    const icon = document.getElementById("output-copy-icon")
    icon.classList.replace("far", "fas")
    icon.classList.replace("fa-clipboard", "fa-clipboard-check")
    clicked.classList.add("btn-success")
}

function resetRow(row) {
    row.getElementsByTagName('input')[0].value = ""
}

function deleteRow(clicked) {
    const thisRow = clicked.closest(".input-group")
    if (thisRow.parentNode.childElementCount === 1) {
        resetRow(thisRow) // Don't allow removing the last row
    } else {
        thisRow.remove()
    }
}

function addRow(clicked) {
    const thisRow = clicked.closest(".input-group")
    const clonedRow = thisRow.cloneNode(true)
    resetRow(clonedRow)
    thisRow.parentNode.insertBefore(clonedRow, thisRow.nextSibling);
}
