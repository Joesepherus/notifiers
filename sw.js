// Install service worker
self.addEventListener('install', function(event) {
    console.log('Service Worker Installed');
    event.waitUntil(
        caches.open('static-cache').then(function(cache) {
            cache.addAll([
                '/manifest.json',
                '/static/js/main.js',
                '/static/img/logo.ico',
                '/static/css/styles.css',
                '/static/img/logo.png',
                '/',
                '/about',
                "/pricing",
                "/alerts",
                "/profile",
                "/reset-password-sent",
                "/reset-password-success",
                "/subscription-success",
                "/subscription-success-temp",
                "/subscription-cancel",
                "/subscription-cancel-temp",
                "/token-expired",
                "/docs",
                "/404",
                "/error",
            ]);
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
