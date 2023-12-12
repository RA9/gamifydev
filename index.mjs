import home from './js/home.js';
import about from './js/about.js';
import ContactPage from './js/contact.js';
import TestPage from './js/tys.js';



if ('serviceWorker' in navigator) {
    navigator.serviceWorker.register('./service-worker.js', {
        type: 'module'
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

