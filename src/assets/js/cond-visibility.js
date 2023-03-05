"use strict";

/*
 * Implements behaviour for the custom data-hide-unless attribute:
 * add or remove the hidden attribute based on the target checkbox state.
 */

import * as htmx from "htmx.org";

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
