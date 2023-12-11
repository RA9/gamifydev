import home from './js/home.js';
import about from './js/about.js';
import ContactPage from './js/contact.js';
import TestPage from './js/tys.js';

if ('serviceWorker' in navigator) {
    navigator.serviceWorker.register('/service-worker.js')
        .then(registration => {
            console.log('Service Worker registered with scope:', registration.scope);
            // update service worker if there is a new version
            if (registration.waiting) {
                registration.waiting.postMessage({type: 'SKIP_WAITING'});
            }
            
            registration.addEventListener('updatefound', () => {
                // An updated service worker has appeared in registration.installing!
                const newWorker = registration.installing;
                newWorker.addEventListener('statechange', () => {
                    // Has service worker state changed?
                    switch (newWorker.state) {
                        case 'installed':
                            // There is a new service worker available, show the notification
                            if (navigator.serviceWorker.controller) {
                                console.log('New content is available; please refresh.');
                            } else {
                                console.log('Content is cached for offline use.');
                            }
                            break;
                    }
                });
            });
        })
        .catch(error => {
            console.error('Service Worker registration failed:', error);
        });
}

function displayContent() {
    const page = document.querySelector('main');
    const currentURL = window.location.href;
    console.log(currentURL);

    switch (currentURL.split('#')[1]) {
        case 'about':
            about(page);
            break;
        case 'contact':
            ContactPage(page);
            break;
        case 'test':
            TestPage(page);
            break;
        default:
            home(page);
    }
}


displayContent();

window.addEventListener('hashchange', displayContent);

