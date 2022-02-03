//
// Listens to HTMX request results and shows an alert div when requests fail.
//

const errorTarget = document.getElementById("htmx-alert")
document.body.addEventListener('htmx:afterRequest', function (evt) {
    if (evt.detail.successful) {
        errorTarget.classList.add("d-none");
        errorTarget.innerText = "";
    } else {
        if (evt.detail.xhr.statusText) {
            errorTarget.innerText = "Unexpected server error: " + evt.detail.xhr.statusText;
        } else {
            errorTarget.innerText = "Unexpected network error, check your connection and try again.";
        }
        errorTarget.classList.remove("d-none");
    }
});
