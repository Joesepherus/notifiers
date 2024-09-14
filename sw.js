// Install service worker
self.addEventListener('install', function(event) {
    console.log('Service Worker Installed');
    event.waitUntil(
        caches.open('static-cache').then(function(cache) {
            cache.addAll([
                '/',
                '/index.html',
                '/manifest.json',
                '/main.js',
                '/icon.png',
                "/index.html",
                "/pricing.html",
                "/about.html",
                "/alerts.html",
                "/profile.html",
                "/reset-password-sent.html",
                "/reset-password-success.html",
                "/subscription-success.html",
                "/subscription-success-temp.html",
                "/subscription-cancel.html",
                "/subscription-cancel-temp.html",
                "/token-expired.html",
                "/docs.html",
                "/404.html",
                "/error.html",
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
