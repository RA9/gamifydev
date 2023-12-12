self.addEventListener("install", (event) => {
  event.waitUntil(
    caches.open("gamifydev-v3").then((cache) => {
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
        cacheNames.filter((name) => name !== "gamifydev-v3").map((name) => {
          return caches.delete(name);
        })
      );
    })
  );
});

self.addEventListener("fetch", async (event) => {
  const cache = await caches.open("gamifydev-v3");
  const response = await cache.match(event.request);
  return response || fetch(event.request).then((networkResponse) => {
    cache.put(event.request, networkResponse.clone());
    return networkResponse;
  });
});