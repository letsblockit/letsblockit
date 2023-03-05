"use strict";

/*
 * Hooks into HTMX afterRequest events and shows an alert div when requests fail.
 */

document.body.addEventListener('htmx:afterRequest', function (evt) {
    const errorTarget = document.getElementById("htmx-alert")
    if (evt.detail.successful) {
        errorTarget.setAttribute("hidden", "true")
        errorTarget.innerText = "";
    } else {
        if (evt.detail.xhr.statusText) {
            errorTarget.innerText = "Unexpected server error: " + evt.detail.xhr.statusText;
        } else {
            errorTarget.innerText = "Unexpected network error, check your connection and try again.";
        }
        errorTarget.removeAttribute("hidden");
    }
});
