// Install service worker
self.addEventListener('install', function(event) {
    console.log('Service Worker Installed');
    event.waitUntil(
        caches.open('static-cache').then(function(cache) {
            cache.addAll([
                '/',
                '/manifest.json',
                '/static/js/main.js',
                '/static/img/logo.png',
                "/templates/index.html",
                "/templates/pricing.html",
                "/templates/about.html",
                "/templates/alerts.html",
                "/templates/profile.html",
                "/templates/reset-password-sent.html",
                "/templates/reset-password-success.html",
                "/templates/subscription-success.html",
                "/templates/subscription-success-temp.html",
                "/templates/subscription-cancel.html",
                "/templates/subscription-cancel-temp.html",
                "/templates/token-expired.html",
                "/templates/docs.html",
                "/templates/404.html",
                "/templates/error.html",
            ]);
        })
    );
});

// Fetch event (for offline support)
self.addEventListener('fetch', function(event) {
    event.respondWith(
        caches.match(event.request).then(function(response) {
            return response || fetch(event.request);
        })
    );
});

// Listen for push notifications
self.addEventListener('push', function(event) {
    const data = event.data.json();
    self.registration.showNotification(data.title, {
        body: data.message,
        icon: '/icon.png'
    });
});
