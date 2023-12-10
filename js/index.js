import home from './home.js';
import about from './about.js';
import ContactPage from './contact.js';
import TestPage from './tys.js';


if ('serviceWorker' in navigator) {
    navigator.serviceWorker.register('/service-worker.js')
        .then(registration => {
            console.log('Service Worker registered with scope:', registration.scope);
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

function clearData() {
    clipboardData.setData('text','')
}



displayContent();

window.addEventListener('hashchange', displayContent);

