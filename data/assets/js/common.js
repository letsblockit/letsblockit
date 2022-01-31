//
// Listens to HTMX request results and shows an alert div when requests fail.
//

const errorTarget = document.getElementById("htmx-alert")
document.body.addEventListener('htmx:afterRequest', function (evt) {
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

//
// Enables triggering item visibility on a checkbox, for optional parameters
//

function toggleConditionalVisibility(checkbox, target) {
    if (checkbox.checked){
        target.removeAttribute("hidden")
    } else {
        target.setAttribute("hidden", "true")
    }
}

function wireConditionalVisibility(content) {
    content.querySelectorAll('[data-hide-unless]').forEach(function (target) {
        const sourceId = target.getAttribute("data-hide-unless")
        const source = document.getElementById(sourceId)
        if (source) {
            source.addEventListener('change', function () {
                toggleConditionalVisibility(source, target)
            })
            toggleConditionalVisibility(source, target)
        } else {
            console.error("cannot find hide-unless source " + sourceId)
        }
    })
}

wireConditionalVisibility(document.body)
htmx.onLoad(wireConditionalVisibility)
