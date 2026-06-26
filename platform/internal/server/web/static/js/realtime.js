// realtime.js — opens a WebSocket to /ws and, on any server push, asks htmx to
// refresh the live regions of the page (elements with hx-trigger="refresh
// from:body"). The server is the source of truth; the socket only signals.
//
// A small [data-live-status] dot, if present, reflects the connection state.
(function () {
  const statusEl = document.querySelector("[data-live-status]");
  const setStatus = (ok) => {
    if (statusEl) statusEl.classList.toggle("is-live", ok);
  };

  let ws = null;
  let retry = 1000;

  function connect() {
    const proto = location.protocol === "https:" ? "wss" : "ws";
    ws = new WebSocket(`${proto}://${location.host}/ws`);

    ws.onopen = () => {
      retry = 1000;
      setStatus(true);
    };
    ws.onmessage = () => {
      // Any topic message triggers a refresh of the live regions.
      document.body.dispatchEvent(new Event("refresh"));
    };
    ws.onclose = () => {
      setStatus(false);
      // Reconnect with capped backoff (handles deploys / sleep / network drops).
      setTimeout(connect, Math.min(retry, 15000));
      retry = Math.min(retry * 2, 15000);
    };
    ws.onerror = () => {
      try { ws.close(); } catch (_) {}
    };
  }

  // Pause when the tab is hidden; reconnect when it returns.
  document.addEventListener("visibilitychange", () => {
    if (document.visibilityState === "visible" && (!ws || ws.readyState === WebSocket.CLOSED)) {
      connect();
    }
  });

  connect();
})();
