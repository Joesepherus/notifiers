// Helper function to wait for a specific time (1 second)
function delay(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

// Install service worker
self.addEventListener('install', function(event) {
    console.log('Service Worker Installed');
    event.waitUntil(
        caches.open('static-cache').then(async function(cache) {
            const staticAssets = [
                '/manifest.json',
                '/static/js/main.js',
                '/static/img/favicon.ico',
                '/static/css/styles.css',
                '/static/img/logo.png'
            ];

            const dynamicPages = [
                '/',
                '/about',
                '/pricing',
            ];

            // Cache static assets one by one with a delay
            for (let asset of staticAssets) {
                await cache.add(asset);
            }

            // Cache dynamic pages one by one with a delay
            for (let page of dynamicPages) {
                await cache.add(page);
                await delay(1000); // 1 second delay
            }
        })
    );
});

// Fetch event (for offline support) with redirect handling
self.addEventListener('fetch', function(event) {
    event.respondWith(
        fetch(event.request).then(function(response) {
            // Check if the response was redirected and handle accordingly
            if (response.redirected) {
                return fetch(response.url);
            }
            return response;
        }).catch(function() {
            // If the request fails, try to serve from cache
            return caches.match(event.request);
        })
    );
});
