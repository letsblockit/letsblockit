/*!
  * Hot-reload of the frontend, added when the server runs with `--reload`.
  * The frontend keeps a SSE stream open to detect when the server restarts
  * (connection closed then successfully re-opened). This triggers a page reload.
  */

let reconnectTimeout = 100
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
        watchEventSource.close()
        setTimeout(watchForRestart, reconnectTimeout);
        if (reconnectTimeout < 2000) {
            reconnectTimeout += 100
        }
    };
}

watchForRestart();
