"use strict";

/*
 * Web component to implement the `onlyif` parameter property:
 * add or remove the hidden attribute based on the target checkbox state.
 */

class condVisible extends HTMLElement {
    setVisibility(source) {
        this.hidden = !source.checked
    }

    connectedCallback() {
        const sourceId = this.getAttribute("only-if")
        const source = document.getElementById(sourceId)
        if (source) {
            source.addEventListener('change', () => this.setVisibility(source))
            this.setVisibility(source)
        } else {
            console.error("cannot find only-if source " + sourceId)
        }
    }
}

customElements.define("cond-visible", condVisible);
