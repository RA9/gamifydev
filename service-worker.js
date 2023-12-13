const cacheName = "gamifydev-v25";
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
        "/data/questions.json",
        "/data/app.json",
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

self.addEventListener("fetch", (event) => {
  event.respondWith(
    caches.open(cacheName).then(async (cache) => {
      return cache.match(event.request).then((response) => {
        return response || fetch(event.request);
      });
    })
  );
});
