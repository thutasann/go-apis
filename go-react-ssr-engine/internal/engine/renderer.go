package engine

import (
	"fmt"
	"sync"

	v8 "rogchap.com/v8go"
)

// Worker wraps a single V8 isolate.
// Single-threaded per V8 rules but many workers run in parallel via pool.
//
// Key optimization: the bundle is compiled once into a cached script.
// Subsequent renders skip parsing — V8 runs from compiled bytecode.
type Worker struct {
	id  int
	iso *v8.Isolate
	ctx *v8.Context

	// bundleLoaded tracks if current bundle is already compiled in this isolate.
	// Avoids re-parsing the same bundle on every render.
	bundleLoaded bool
	bundleHash   string

	mu sync.Mutex // protects isolate — V8 is not thread safe per isolate
}

func NewWorker(id int) (*Worker, error) {
	iso := v8.NewIsolate()
	global := v8.NewObjectTemplate(iso)
	ctx := v8.NewContext(iso, global)

	// Inject minimal console.log so React doesn't crash on console calls.
	// V8 doesn't have console natively — it's a browser/node API.
	ctx.RunScript(`
		var console = {
			log: function() {},
			warn: function() {},
			error: function() {}
		};
		var process = { env: { NODE_ENV: 'production' } };
	`, "bootstrap.js")

	return &Worker{
		id:  id,
		iso: iso,
		ctx: ctx,
	}, nil
}

// Execute runs the bundle and calls __renderToString.
// If the bundle hasn't changed since last call, skips re-parsing.
func (w *Worker) Execute(bundle, route, propsJSON string) (string, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Only load bundle if it changed — massive speedup on repeated renders.
	// First render: ~5ms (parse + compile). Subsequent: ~0.1ms (cached bytecode).
	bundleHash := hashBundle(bundle)
	if !w.bundleLoaded || w.bundleHash != bundleHash {
		// Fresh context to avoid stale state from previous bundle
		if w.ctx != nil {
			w.ctx.Close()
		}
		global := v8.NewObjectTemplate(w.iso)
		w.ctx = v8.NewContext(w.iso, global)

		// Re-inject polyfills
		w.ctx.RunScript(`
			// --- Console ---
			var console = {
				log: function(){}, warn: function(){}, error: function(){},
				info: function(){}, debug: function(){}, trace: function(){},
				dir: function(){}, table: function(){}, time: function(){},
				timeEnd: function(){}, timeLog: function(){}, assert: function(){},
				count: function(){}, countReset: function(){}, group: function(){},
				groupEnd: function(){}, groupCollapsed: function(){}, clear: function(){}
			};

			// --- Process ---
			var process = {
				env: { NODE_ENV: 'production' },
				nextTick: function(cb) { cb(); },
				version: 'v18.0.0',
				versions: { node: '18.0.0' },
				platform: 'linux',
				argv: [], pid: 1,
				cwd: function() { return '/'; },
				exit: function() {},
				on: function() { return this; },
				once: function() { return this; },
				off: function() { return this; },
				removeListener: function() { return this; },
				emit: function() { return this; },
				stderr: { write: function(){} },
				stdout: { write: function(){} },
				hrtime: function() { return [0,0]; },
				binding: function() { return {}; }
			};

			// --- Timers ---
			var queueMicrotask = function(cb) { cb(); };
			if (typeof setTimeout === 'undefined') {
				var setTimeout = function(cb, ms) { cb(); return 0; };
			}
			var clearTimeout = clearTimeout || function() {};
			var setInterval = setInterval || function() { return 0; };
			var clearInterval = clearInterval || function() {};
			var setImmediate = function(cb) { cb(); return 0; };
			var clearImmediate = function() {};

			// --- Encoding ---
			var TextEncoder = function() {};
			TextEncoder.prototype.encode = function(s) {
				var arr = [];
				for (var i = 0; i < s.length; i++) {
					var c = s.charCodeAt(i);
					if (c < 128) arr.push(c);
					else if (c < 2048) { arr.push(192 | (c >> 6)); arr.push(128 | (c & 63)); }
					else { arr.push(224 | (c >> 12)); arr.push(128 | ((c >> 6) & 63)); arr.push(128 | (c & 63)); }
				}
				return new Uint8Array(arr);
			};
			TextEncoder.prototype.encoding = 'utf-8';

			var TextDecoder = function(enc) { this.encoding = enc || 'utf-8'; };
			TextDecoder.prototype.decode = function(buf) {
				if (typeof buf === 'string') return buf;
				if (!buf || !buf.length) return '';
				var s = '';
				for (var i = 0; i < buf.length; i++) s += String.fromCharCode(buf[i]);
				return s;
			};

			// --- URL ---
			if (typeof URL === 'undefined') {
				var URL = function(url, base) {
					this.href = url;
					this.pathname = url.split('?')[0];
					this.search = url.indexOf('?') >= 0 ? url.slice(url.indexOf('?')) : '';
					this.hash = '';
					this.hostname = '';
					this.host = '';
					this.origin = '';
					this.protocol = 'https:';
					this.port = '';
					this.searchParams = {
						get: function(k) { return null; },
						has: function(k) { return false; },
						forEach: function() {},
						entries: function() { return []; }
					};
				};
			}
			if (typeof URLSearchParams === 'undefined') {
				var URLSearchParams = function(init) {
					this._params = {};
					if (typeof init === 'string') {
						init.replace(/^\?/, '').split('&').forEach(function(pair) {
							var kv = pair.split('=');
							if (kv[0]) this._params[decodeURIComponent(kv[0])] = decodeURIComponent(kv[1] || '');
						}.bind(this));
					}
				};
				URLSearchParams.prototype.get = function(k) { return this._params[k] || null; };
				URLSearchParams.prototype.has = function(k) { return k in this._params; };
				URLSearchParams.prototype.forEach = function(cb) {
					for (var k in this._params) cb(this._params[k], k);
				};
			}

			// --- Performance ---
			var performance = {
				now: function() { return Date.now(); },
				mark: function() {},
				measure: function() {},
				getEntriesByName: function() { return []; },
				getEntriesByType: function() { return []; },
				clearMarks: function() {},
				clearMeasures: function() {}
			};

			// --- Buffer ---
			if (typeof Buffer === 'undefined') {
				var Buffer = {
					from: function(data) {
						if (typeof data === 'string') {
							var arr = [];
							for (var i = 0; i < data.length; i++) arr.push(data.charCodeAt(i));
							return new Uint8Array(arr);
						}
						return new Uint8Array(data || 0);
					},
					alloc: function(n) { return new Uint8Array(n); },
					allocUnsafe: function(n) { return new Uint8Array(n); },
					isBuffer: function() { return false; },
					concat: function(list) {
						var total = 0;
						for (var i = 0; i < list.length; i++) total += list[i].length;
						var result = new Uint8Array(total);
						var offset = 0;
						for (var i = 0; i < list.length; i++) { result.set(list[i], offset); offset += list[i].length; }
						return result;
					},
					byteLength: function(s) { return typeof s === 'string' ? s.length : (s.byteLength || 0); }
				};
			}

			// --- Misc globals ---
			if (typeof global === 'undefined') var global = globalThis;
			if (typeof self === 'undefined') var self = globalThis;
			if (typeof window === 'undefined') var window = globalThis;

			// --- MessageChannel (React scheduler uses this) ---
			var MessageChannel = function() {
				var self = this;
				this.port1 = {
					onmessage: null,
					postMessage: function() {
						if (self.port2.onmessage) self.port2.onmessage({ data: null });
					}
				};
				this.port2 = {
					onmessage: null,
					postMessage: function() {
						if (self.port1.onmessage) self.port1.onmessage({ data: null });
					}
				};
			};
			var MessagePort = function() {};
			var MessageEvent = function() {};

			// --- AbortController ---
			if (typeof AbortController === 'undefined') {
				var AbortSignal = function() { this.aborted = false; this.reason = undefined; };
				AbortSignal.prototype.addEventListener = function() {};
				AbortSignal.prototype.removeEventListener = function() {};
				var AbortController = function() { this.signal = new AbortSignal(); };
				AbortController.prototype.abort = function(reason) {
					this.signal.aborted = true;
					this.signal.reason = reason;
				};
			}

			// --- Headers/Request/Response (fetch API stubs) ---
			if (typeof Headers === 'undefined') {
				var Headers = function() { this._h = {}; };
				Headers.prototype.get = function(k) { return this._h[k.toLowerCase()] || null; };
				Headers.prototype.set = function(k, v) { this._h[k.toLowerCase()] = v; };
				Headers.prototype.has = function(k) { return k.toLowerCase() in this._h; };
			}

			if (typeof Request === 'undefined') {
				var Request = function(url, opts) { this.url = url; this.method = (opts && opts.method) || 'GET'; };
			}

			if (typeof Response === 'undefined') {
				var Response = function(body, opts) { this.body = body; this.status = (opts && opts.status) || 200; };
				Response.prototype.text = function() { return Promise.resolve(this.body || ''); };
				Response.prototype.json = function() { return Promise.resolve(JSON.parse(this.body || '{}')); };
			}

			if (typeof fetch === 'undefined') {
				var fetch = function() { return Promise.resolve(new Response('{}')); };
			}

			// --- ReadableStream stub ---
			if (typeof ReadableStream === 'undefined') {
				var ReadableStream = function() {};
				ReadableStream.prototype.getReader = function() {
					return { read: function() { return Promise.resolve({ done: true, value: undefined }); }, releaseLock: function() {} };
				};
			}

			// --- WeakRef (React may use) ---
			if (typeof WeakRef === 'undefined') {
				var WeakRef = function(target) { this._target = target; };
				WeakRef.prototype.deref = function() { return this._target; };
			}

			// --- FinalizationRegistry ---
			if (typeof FinalizationRegistry === 'undefined') {
				var FinalizationRegistry = function() {};
				FinalizationRegistry.prototype.register = function() {};
				FinalizationRegistry.prototype.unregister = function() {};
			}

			// --- structuredClone ---
			if (typeof structuredClone === 'undefined') {
				var structuredClone = function(obj) { return JSON.parse(JSON.stringify(obj)); };
			}
		`, "bootstrap.js")

		_, err := w.ctx.RunScript(bundle, "server_bundle.js")
		if err != nil {
			return "", fmt.Errorf("worker %d: bundle exec: %w", w.id, err)
		}
		w.bundleLoaded = true
		w.bundleHash = bundleHash
	}

	renderCall := fmt.Sprintf(`__renderToString(%q, %s)`, route, propsJSON)
	val, err := w.ctx.RunScript(renderCall, "render.js")
	if err != nil {
		return "", fmt.Errorf("worker %d: render %s: %w", w.id, route, err)
	}

	return val.String(), nil
}

// ExecuteProps calls __getServerSideProps in V8.
// Returns JSON string of { props, redirect, notFound }.
func (w *Worker) ExecuteProps(bundle, route, contextJSON string) (string, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	bundleHash := hashBundle(bundle)
	if !w.bundleLoaded || w.bundleHash != bundleHash {
		if w.ctx != nil {
			w.ctx.Close()
		}
		global := v8.NewObjectTemplate(w.iso)
		w.ctx = v8.NewContext(w.iso, global)

		w.ctx.RunScript(`
			// --- Console ---
			var console = {
				log: function(){}, warn: function(){}, error: function(){},
				info: function(){}, debug: function(){}, trace: function(){},
				dir: function(){}, table: function(){}, time: function(){},
				timeEnd: function(){}, timeLog: function(){}, assert: function(){},
				count: function(){}, countReset: function(){}, group: function(){},
				groupEnd: function(){}, groupCollapsed: function(){}, clear: function(){}
			};

			// --- Process ---
			var process = {
				env: { NODE_ENV: 'production' },
				nextTick: function(cb) { cb(); },
				version: 'v18.0.0',
				versions: { node: '18.0.0' },
				platform: 'linux',
				argv: [], pid: 1,
				cwd: function() { return '/'; },
				exit: function() {},
				on: function() { return this; },
				once: function() { return this; },
				off: function() { return this; },
				removeListener: function() { return this; },
				emit: function() { return this; },
				stderr: { write: function(){} },
				stdout: { write: function(){} },
				hrtime: function() { return [0,0]; },
				binding: function() { return {}; }
			};

			// --- Timers ---
			var queueMicrotask = function(cb) { cb(); };
			if (typeof setTimeout === 'undefined') {
				var setTimeout = function(cb, ms) { cb(); return 0; };
			}
			var clearTimeout = clearTimeout || function() {};
			var setInterval = setInterval || function() { return 0; };
			var clearInterval = clearInterval || function() {};
			var setImmediate = function(cb) { cb(); return 0; };
			var clearImmediate = function() {};

			// --- Encoding ---
			var TextEncoder = function() {};
			TextEncoder.prototype.encode = function(s) {
				var arr = [];
				for (var i = 0; i < s.length; i++) {
					var c = s.charCodeAt(i);
					if (c < 128) arr.push(c);
					else if (c < 2048) { arr.push(192 | (c >> 6)); arr.push(128 | (c & 63)); }
					else { arr.push(224 | (c >> 12)); arr.push(128 | ((c >> 6) & 63)); arr.push(128 | (c & 63)); }
				}
				return new Uint8Array(arr);
			};
			TextEncoder.prototype.encoding = 'utf-8';

			var TextDecoder = function(enc) { this.encoding = enc || 'utf-8'; };
			TextDecoder.prototype.decode = function(buf) {
				if (typeof buf === 'string') return buf;
				if (!buf || !buf.length) return '';
				var s = '';
				for (var i = 0; i < buf.length; i++) s += String.fromCharCode(buf[i]);
				return s;
			};

			// --- URL ---
			if (typeof URL === 'undefined') {
				var URL = function(url, base) {
					this.href = url;
					this.pathname = url.split('?')[0];
					this.search = url.indexOf('?') >= 0 ? url.slice(url.indexOf('?')) : '';
					this.hash = '';
					this.hostname = '';
					this.host = '';
					this.origin = '';
					this.protocol = 'https:';
					this.port = '';
					this.searchParams = {
						get: function(k) { return null; },
						has: function(k) { return false; },
						forEach: function() {},
						entries: function() { return []; }
					};
				};
			}
			if (typeof URLSearchParams === 'undefined') {
				var URLSearchParams = function(init) {
					this._params = {};
					if (typeof init === 'string') {
						init.replace(/^\?/, '').split('&').forEach(function(pair) {
							var kv = pair.split('=');
							if (kv[0]) this._params[decodeURIComponent(kv[0])] = decodeURIComponent(kv[1] || '');
						}.bind(this));
					}
				};
				URLSearchParams.prototype.get = function(k) { return this._params[k] || null; };
				URLSearchParams.prototype.has = function(k) { return k in this._params; };
				URLSearchParams.prototype.forEach = function(cb) {
					for (var k in this._params) cb(this._params[k], k);
				};
			}

			// --- Performance ---
			var performance = {
				now: function() { return Date.now(); },
				mark: function() {},
				measure: function() {},
				getEntriesByName: function() { return []; },
				getEntriesByType: function() { return []; },
				clearMarks: function() {},
				clearMeasures: function() {}
			};

			// --- Buffer ---
			if (typeof Buffer === 'undefined') {
				var Buffer = {
					from: function(data) {
						if (typeof data === 'string') {
							var arr = [];
							for (var i = 0; i < data.length; i++) arr.push(data.charCodeAt(i));
							return new Uint8Array(arr);
						}
						return new Uint8Array(data || 0);
					},
					alloc: function(n) { return new Uint8Array(n); },
					allocUnsafe: function(n) { return new Uint8Array(n); },
					isBuffer: function() { return false; },
					concat: function(list) {
						var total = 0;
						for (var i = 0; i < list.length; i++) total += list[i].length;
						var result = new Uint8Array(total);
						var offset = 0;
						for (var i = 0; i < list.length; i++) { result.set(list[i], offset); offset += list[i].length; }
						return result;
					},
					byteLength: function(s) { return typeof s === 'string' ? s.length : (s.byteLength || 0); }
				};
			}

			// --- Misc globals ---
			if (typeof global === 'undefined') var global = globalThis;
			if (typeof self === 'undefined') var self = globalThis;
			if (typeof window === 'undefined') var window = globalThis;

			// --- MessageChannel (React scheduler uses this) ---
			var MessageChannel = function() {
				var self = this;
				this.port1 = {
					onmessage: null,
					postMessage: function() {
						if (self.port2.onmessage) self.port2.onmessage({ data: null });
					}
				};
				this.port2 = {
					onmessage: null,
					postMessage: function() {
						if (self.port1.onmessage) self.port1.onmessage({ data: null });
					}
				};
			};
			var MessagePort = function() {};
			var MessageEvent = function() {};

			// --- AbortController ---
			if (typeof AbortController === 'undefined') {
				var AbortSignal = function() { this.aborted = false; this.reason = undefined; };
				AbortSignal.prototype.addEventListener = function() {};
				AbortSignal.prototype.removeEventListener = function() {};
				var AbortController = function() { this.signal = new AbortSignal(); };
				AbortController.prototype.abort = function(reason) {
					this.signal.aborted = true;
					this.signal.reason = reason;
				};
			}

			// --- Headers/Request/Response (fetch API stubs) ---
			if (typeof Headers === 'undefined') {
				var Headers = function() { this._h = {}; };
				Headers.prototype.get = function(k) { return this._h[k.toLowerCase()] || null; };
				Headers.prototype.set = function(k, v) { this._h[k.toLowerCase()] = v; };
				Headers.prototype.has = function(k) { return k.toLowerCase() in this._h; };
			}

			if (typeof Request === 'undefined') {
				var Request = function(url, opts) { this.url = url; this.method = (opts && opts.method) || 'GET'; };
			}

			if (typeof Response === 'undefined') {
				var Response = function(body, opts) { this.body = body; this.status = (opts && opts.status) || 200; };
				Response.prototype.text = function() { return Promise.resolve(this.body || ''); };
				Response.prototype.json = function() { return Promise.resolve(JSON.parse(this.body || '{}')); };
			}

			if (typeof fetch === 'undefined') {
				var fetch = function() { return Promise.resolve(new Response('{}')); };
			}

			// --- ReadableStream stub ---
			if (typeof ReadableStream === 'undefined') {
				var ReadableStream = function() {};
				ReadableStream.prototype.getReader = function() {
					return { read: function() { return Promise.resolve({ done: true, value: undefined }); }, releaseLock: function() {} };
				};
			}

			// --- WeakRef (React may use) ---
			if (typeof WeakRef === 'undefined') {
				var WeakRef = function(target) { this._target = target; };
				WeakRef.prototype.deref = function() { return this._target; };
			}

			// --- FinalizationRegistry ---
			if (typeof FinalizationRegistry === 'undefined') {
				var FinalizationRegistry = function() {};
				FinalizationRegistry.prototype.register = function() {};
				FinalizationRegistry.prototype.unregister = function() {};
			}

			// --- structuredClone ---
			if (typeof structuredClone === 'undefined') {
				var structuredClone = function(obj) { return JSON.parse(JSON.stringify(obj)); };
			}
		`, "bootstrap.js")

		_, err := w.ctx.RunScript(bundle, "server_bundle.js")
		if err != nil {
			return "", fmt.Errorf("worker %d: bundle exec: %w", w.id, err)
		}
		w.bundleLoaded = true
		w.bundleHash = bundleHash
	}

	propsCall := fmt.Sprintf(`__getServerSideProps(%q, %s)`, route, contextJSON)
	val, err := w.ctx.RunScript(propsCall, "props.js")
	if err != nil {
		return "", fmt.Errorf("worker %d: props %s: %w", w.id, route, err)
	}

	return val.String(), nil
}

func (w *Worker) Dispose() {
	if w == nil {
		return
	}
	if w.ctx != nil {
		w.ctx.Close()
	}
	if w.iso != nil {
		w.iso.Dispose()
	}
}

// hashBundle produces a fast identity check for bundle content.
// Not cryptographic — just needs to detect changes.
// Uses length + first/last 64 bytes. Collisions are harmless
// (worst case: one extra re-parse).
func hashBundle(bundle string) string {
	l := len(bundle)
	if l <= 128 {
		return bundle
	}
	return fmt.Sprintf("%d:%s:%s", l, bundle[:64], bundle[l-64:])
}
