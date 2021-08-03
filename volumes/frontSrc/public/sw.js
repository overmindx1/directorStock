self.addEventListener('install', (e) => {
    console.log('[Service Worker] - Install');
    self.skipWaiting();
});

const cacheName = 'directorStockCache';

self.addEventListener('activate', function(event) {
    console.log('Claiming control');
    return self.clients.claim();
});

// Cache Test
self.addEventListener('fetch', function (event) {    
    //console.log('on fetch sw')
    
    event.respondWith( (async () => {
        try {
            //console.log(event.request)
            let response = await fetch(event.request);
            let responseToCache = response.clone();
            //console.log(responseToCache)
            let cacheObj = await caches.open(cacheName);
            cacheObj.put(event.request, responseToCache);
            return response; 
        } catch (error) {
            console.error(error)
            const cachedResponse = await caches.match(event.request);
            if (cachedResponse) {
                return cachedResponse
            };
            return error
        }
    })());
});