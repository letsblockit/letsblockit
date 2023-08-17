"use strict";

/*
 * Hooks into HTMX afterRequest events and shows an alert div when requests fail.
 */

document.body.addEventListener('htmx:afterRequest', function (evt) {
    const errorTarget = document.getElementById("htmx-alert")
    if (evt.detail.error || evt.detail.failed) {
        if (evt.detail.xhr.statusText) {
            errorTarget.innerText = "Unexpected server error: " + evt.detail.xhr.statusText;
        } else {
            errorTarget.innerText = "Unexpected network error, check your connection and try again.";
        }
        errorTarget.removeAttribute("hidden");
    } else {
        const errorTarget = document.getElementById("htmx-alert")
        errorTarget.setAttribute("hidden", "true")
        errorTarget.innerText = "";
    }
});
