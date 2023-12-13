// import home from './js/home.js';
// import about from './js/about.js';
// import ContactPage from './js/contact.js';
// import TestPage from './js/tys.js';



if ('serviceWorker' in navigator) {
    navigator.serviceWorker.register('./service-worker.js', {
        type: 'module'
    });

    navigator.serviceWorker.addEventListener('message', (event) => {
        if (event.data && event.data.type === 'CACHE_UPDATED') {
          // Reload the page when the cache is updated
          window.location.reload(true);
        }
      });
}


function displayContent() {
    const page = document.querySelector('main');
    const currentURL = (window.location.href).split('#')[1];
    console.log(currentURL);

    switch (currentURL) {
        case 'about':
            console.log('about');
            AboutPage(page);
            break;
        case 'contact':
            ContactPage(page);
            break;
        case 'test':
            TestPage(page);
            break;
        default:
            HomePage(page);
    }
}


displayContent();

window.addEventListener('hashchange', displayContent);

