const cache = "neringa";

self.addEventListener("install", e => {
  self.skipWaiting();
  e.waitUntil(caches.open(cache).then(c => c.addAll(["./", "main.wasm"])));
});

self.addEventListener("activate", e => e.waitUntil(self.clients.claim()));

self.addEventListener("fetch", e => {
  const u = new URL(e.request.url);
  if (e.request.method !== "GET" || u.origin !== location.origin || u.pathname.startsWith("/tvarkarastis/")) return;
  e.respondWith(caches.open(cache).then(async c => {
    const hit = await c.match(e.request);
    const net = fetch(e.request).then(r => {
      if (r.ok) c.put(e.request, r.clone());
      return r;
    });
    if (hit) {
      e.waitUntil(net.catch(() => {}));
      return hit;
    }
    return net;
  }));
});
