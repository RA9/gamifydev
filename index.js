// if ('serviceWorker' in navigator) {
//     navigator.serviceWorker.register('./service-worker.js');

//     navigator.serviceWorker.addEventListener('message', (event) => {
//         if (event.data && event.data.type === 'CACHE_UPDATED') {
//           // Reload the page when the cache is updated
//           window.location.reload(true);
//         }
//       });
// }


function displayContent() {
    const page = document.querySelector('main');
    const currentURL = (window.location.href).split('#')[1];
    // console.log(currentURL);
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

