/*!
  * Hot-reload of the frontend, added when the server runs with `--reload`.
  * The frontend keeps a SSE stream open to detect when the server restarts
  * (connection closed then successfully re-opened). This triggers a page reload.
  */

let isConnected = false
let watchEventSource = null;

function watchForRestart() {
    watchEventSource = new EventSource('/should-reload');
    watchEventSource.onopen = function () {
        if (isConnected) {
            console.log("Restart detected, reloading page");
            location.reload()
        } else {
            console.log("Connection to server opened.");
            isConnected = true
        }
    }
    watchEventSource.onerror = function () {
        if (watchEventSource.readyState === 2) {
            setTimeout(watchForRestart, 100);
        }
    };
}
watchForRestart();
