const cache = "neringa";

self.addEventListener("install", e => {
  self.skipWaiting();
  e.waitUntil(caches.open(cache).then(c => c.addAll(["./", "main.wasm"].map(u => new Request(u, {cache: "reload"})))));
});

self.addEventListener("activate", e => e.waitUntil(self.clients.claim()));

self.addEventListener("fetch", e => {
  const u = new URL(e.request.url);
  if (e.request.method !== "GET" || u.origin !== location.origin || /^\/(tvarkarastis|gyventojui|data|v1)\//.test(u.pathname)) return;
  e.respondWith(caches.open(cache).then(async c => {
    const hit = await c.match(e.request);
    const net = fetch(e.request.url, {cache: "no-cache"}).then(r => {
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
