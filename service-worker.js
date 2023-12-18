const cacheName = "gamifydev-v0.72";

self.addEventListener("install", (event) => {
  // if there is a new service worker, skip waiting
  self.skipWaiting();

  event.waitUntil(
    caches.open(cacheName).then((cache) => {
      return cache.addAll([
        "/",
        "/index.html",
        "/css/main.css",
        "/index.js",
        "/js/tys.js",
        "/js/utils.js",
        "/js/about.js",
        "/js/contact.js",
        "/js/quiz.js",
        "/js/external/dexie.js",
        "/js/tailwind.min.js",
        "/js/storage.js",
        "/js/modules.js",
        "/data/questions.json",
        "/data/app.json",
        "/data/notes/history_of_the_web.md",
        "/data/notes/exercises.json",
        // // '/images/icon.png'
      ]);
    })
  );
});

self.addEventListener("activate", (event) => {
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames
          .filter((name) => name !== cacheName)
          .map((name) => {
            return caches.delete(name);
          })
      );
    })
  );
});

// self.addEventListener("fetch", (event) => {
//   event.respondWith(
//     caches.open(cacheName).then(async (cache) => {
//       return cache.match(event.request).then((response) => {
//         return response || fetch(event.request);
//       });
//     })
//   );
// });

function isCacheable(request) {
  const url = new URL(request.url);
  return !url.pathname.endsWith(".json");
}

async function cacheFirstWithRefresh(request) {
  const fetchResponsePromise = fetch(request).then(async (networkResponse) => {
    if (networkResponse.ok) {
      const cache = await caches.open(cacheName);
      if (
        request.url.startsWith("chrome-extension") ||
        request.url.includes("extension") ||
        !(request.url.indexOf("http") === 0)
      ) {
        return;
      }
      cache.put(request, networkResponse.clone());
    }
    return networkResponse;
  });

  // if (
  //   request.url.startsWith("chrome-extension") ||
  //   request.url.includes("extension") ||
  //   !(request.url.indexOf("http") === 0)
  // ) {
  //   return;
  // }

  return (await caches.match(request)) || (await fetchResponsePromise);
}

self.addEventListener("fetch", (event) => {
  if (isCacheable(event.request)) {
    event.respondWith(cacheFirstWithRefresh(event.request));
  }
});
