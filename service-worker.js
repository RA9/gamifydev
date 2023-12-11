self.addEventListener('install', event => {
    event.waitUntil(
        caches.open('gamifydev-v1').then(cache => {
            return cache.addAll([
                '/',
                '/index.html',
                '/css/main.css',
                '/index.js',
                '/js/tys.js',
                '/js/utils.js',
                '/js/about.js',
                '/js/contact.js',
                '/js/quiz.js',
                '/js/external/dexie.js',
                '/js/tailwind.min.js',
                '/js/storage.js',
                '/data/questions.json',
                '/data/app.json'
                // // '/images/icon.png'
            ]);
        })
    );
});



self.addEventListener('fetch', event => {
    event.respondWith(
        caches.match(event.request).then(response => {
            return response || fetch(event.request).catch(error => {
                console.error('Error fetching', event.request.url, error);
            });
        })
    );
});

// update caches with the latest
self.addEventListener('activate', event => {
    const cacheWhitelist = ['gamifydev-v1'];
    event.waitUntil(
        caches.keys().then(keyList => {
            return Promise.all(keyList.map(key => {
                if (!cacheWhitelist.includes(key)) {
                    return caches.delete(key);
                }
            }));
        })
    );
});


self.addEventListener('fetch', event => {
    event.respondWith(
        caches.match(event.request).then(response => {
            return response || fetch(event.request);
        })
    );
});
