package bundler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/thutasann/go-react-ssr-engine/internal/config"
)

type Bundler struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Bundler {
	return &Bundler{cfg: cfg}
}

type BuildResult struct {
	ServerBundle  string
	ClientEntries map[string]string
}

func (b *Bundler) Build() (*BuildResult, error) {
	entries, err := b.discoverPages()
	if err != nil {
		return nil, fmt.Errorf("bundler: page discovery failed: %w", err)
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("bundler: no pages found in %s", b.cfg.PagesDir)
	}

	serverJS, err := b.buildServer(entries)
	if err != nil {
		return nil, fmt.Errorf("bundler: server build failed: %w", err)
	}

	clientEntries, err := b.buildClient(entries)
	if err != nil {
		return nil, fmt.Errorf("bundler: client build failed: %w", err)
	}

	return &BuildResult{
		ServerBundle:  serverJS,
		ClientEntries: clientEntries,
	}, nil
}

func (b *Bundler) discoverPages() ([]string, error) {
	var entries []string

	err := filepath.Walk(b.cfg.PagesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if ext != ".tsx" && ext != ".jsx" {
			return nil
		}

		base := filepath.Base(path)
		if strings.HasPrefix(base, "_") {
			return nil
		}

		entries = append(entries, path)
		return nil
	})

	return entries, err
}

func (b *Bundler) buildServer(entries []string) (string, error) {
	virtualEntry := b.generateServerEntry(entries)

	serverDir := filepath.Join(b.cfg.BuildDir, "server")
	os.MkdirAll(serverDir, 0755)

	entryPath := filepath.Join(serverDir, "_entry.jsx")
	if err := os.WriteFile(entryPath, []byte(virtualEntry), 0644); err != nil {
		return "", err
	}

	absNodeModules, _ := filepath.Abs("node_modules")

	// Shim out Node.js builtins that React references but never actually
	// calls in renderToString. These modules don't exist in V8.
	// Marking them external tells esbuild to emit require('util') as-is,
	// but since we use IIFE format that would fail. Instead we inject
	// empty shims so the references resolve to no-ops.
	shimDir := filepath.Join(serverDir, "_shims")
	os.MkdirAll(shimDir, 0755)

	shims := map[string]string{
		"util": `
var TextEncoder = function() {};
TextEncoder.prototype.encode = function(s) {
  var arr = [];
  for (var i = 0; i < s.length; i++) arr.push(s.charCodeAt(i));
  return new Uint8Array(arr);
};
var TextDecoder = function() {};
TextDecoder.prototype.decode = function(arr) {
  if (typeof arr === 'string') return arr;
  var s = '';
  for (var i = 0; i < arr.length; i++) s += String.fromCharCode(arr[i]);
  return s;
};
module.exports = {
  TextEncoder: TextEncoder,
  TextDecoder: TextDecoder,
  inspect: function(obj) { return JSON.stringify(obj); },
  format: function() { return Array.prototype.slice.call(arguments).join(' '); },
  deprecate: function(fn) { return fn; },
  inherits: function(ctor, superCtor) {
    ctor.prototype = Object.create(superCtor.prototype);
    ctor.prototype.constructor = ctor;
  },
  isArray: Array.isArray,
  isBuffer: function() { return false; },
  isNull: function(v) { return v === null; },
  isNullOrUndefined: function(v) { return v == null; },
  isUndefined: function(v) { return v === undefined; },
  isString: function(v) { return typeof v === 'string'; },
  isNumber: function(v) { return typeof v === 'number'; },
  isObject: function(v) { return typeof v === 'object' && v !== null; },
  isFunction: function(v) { return typeof v === 'function'; },
  isRegExp: function(v) { return v instanceof RegExp; },
  isDate: function(v) { return v instanceof Date; },
  isError: function(v) { return v instanceof Error; },
  isPrimitive: function(v) { return v === null || (typeof v !== 'object' && typeof v !== 'function'); },
  types: {
    isUint8Array: function(v) { return v instanceof Uint8Array; },
    isArrayBuffer: function(v) { return v instanceof ArrayBuffer; },
  },
  promisify: function(fn) { return fn; }
};`,
		"stream": `
function Stream() {}
Stream.prototype.on = function() { return this; };
Stream.prototype.once = function() { return this; };
Stream.prototype.emit = function() { return this; };
Stream.prototype.pipe = function() { return this; };
Stream.prototype.removeListener = function() { return this; };
Stream.prototype.addListener = function() { return this; };
Stream.Readable = Stream;
Stream.Writable = Stream;
Stream.Duplex = Stream;
Stream.Transform = Stream;
Stream.PassThrough = Stream;
Stream.Stream = Stream;
module.exports = Stream;`,
		"crypto": `
module.exports = {
  createHash: function() {
    return {
      update: function() { return this; },
      digest: function() { return 'mock-hash'; }
    };
  },
  randomBytes: function(n) { return new Uint8Array(n); },
  randomUUID: function() { return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'; }
};`,
		"buffer": `
function Buffer(arg) {
  if (typeof arg === 'number') return new Uint8Array(arg);
  if (typeof arg === 'string') {
    var arr = [];
    for (var i = 0; i < arg.length; i++) arr.push(arg.charCodeAt(i));
    return new Uint8Array(arr);
  }
  return new Uint8Array(arg || 0);
}
Buffer.from = function(data, encoding) {
  if (data instanceof Uint8Array) return data;
  if (typeof data === 'string') return new Buffer(data);
  return new Uint8Array(data || 0);
};
Buffer.alloc = function(n) { return new Uint8Array(n); };
Buffer.allocUnsafe = function(n) { return new Uint8Array(n); };
Buffer.isBuffer = function(obj) { return obj instanceof Uint8Array; };
Buffer.concat = function(list) {
  var total = 0;
  for (var i = 0; i < list.length; i++) total += list[i].length;
  var result = new Uint8Array(total);
  var offset = 0;
  for (var i = 0; i < list.length; i++) {
    result.set(list[i], offset);
    offset += list[i].length;
  }
  return result;
};
Buffer.byteLength = function(s) { return typeof s === 'string' ? s.length : s.byteLength || 0; };
module.exports = { Buffer: Buffer };`,
		"process": `
module.exports = {
  env: { NODE_ENV: 'production' },
  nextTick: function(cb) { cb(); },
  version: 'v18.0.0',
  versions: { node: '18.0.0' },
  platform: 'linux',
  stderr: { write: function() {} },
  stdout: { write: function() {} },
  argv: [],
  pid: 1,
  cwd: function() { return '/'; },
  exit: function() {},
  on: function() { return this; },
  once: function() { return this; },
  off: function() { return this; },
  removeListener: function() { return this; },
  emit: function() { return this; },
  binding: function() { return {}; },
  hrtime: function() { return [0, 0]; }
};`,
		"async_hooks": `
function AsyncLocalStorage() {
  this._store = undefined;
}
AsyncLocalStorage.prototype.getStore = function() { return this._store; };
AsyncLocalStorage.prototype.run = function(store, fn) {
  this._store = store;
  try { return fn(); } finally { this._store = undefined; }
};
AsyncLocalStorage.prototype.enterWith = function(store) { this._store = store; };
AsyncLocalStorage.prototype.disable = function() { this._store = undefined; };
module.exports = {
  createHook: function() { return { enable: function(){}, disable: function(){} }; },
  executionAsyncId: function() { return 0; },
  triggerAsyncId: function() { return 0; },
  executionAsyncResource: function() { return {}; },
  AsyncLocalStorage: AsyncLocalStorage,
  AsyncResource: function(type) { this.type = type; }
};`,
		"events": `
function EventEmitter() { this._events = {}; }
EventEmitter.prototype.on = function(e, fn) {
  if (!this._events[e]) this._events[e] = [];
  this._events[e].push(fn);
  return this;
};
EventEmitter.prototype.once = function(e, fn) { return this.on(e, fn); };
EventEmitter.prototype.off = function(e, fn) {
  if (this._events[e]) this._events[e] = this._events[e].filter(function(f) { return f !== fn; });
  return this;
};
EventEmitter.prototype.removeListener = EventEmitter.prototype.off;
EventEmitter.prototype.addListener = EventEmitter.prototype.on;
EventEmitter.prototype.emit = function(e) {
  var args = Array.prototype.slice.call(arguments, 1);
  if (this._events[e]) this._events[e].forEach(function(fn) { fn.apply(null, args); });
  return this;
};
EventEmitter.prototype.removeAllListeners = function(e) {
  if (e) delete this._events[e]; else this._events = {};
  return this;
};
EventEmitter.prototype.listenerCount = function(e) {
  return this._events[e] ? this._events[e].length : 0;
};
EventEmitter.prototype.setMaxListeners = function() { return this; };
module.exports = EventEmitter;
module.exports.EventEmitter = EventEmitter;`,
		"path":                "module.exports = { join: function() { return Array.prototype.slice.call(arguments).join('/'); }, resolve: function() { return Array.prototype.slice.call(arguments).join('/'); }, dirname: function(p) { return p.split('/').slice(0,-1).join('/'); }, basename: function(p) { var parts = p.split('/'); return parts[parts.length-1]; }, extname: function(p) { var m = p.match(/\\.[^.]+$/); return m ? m[0] : ''; }, sep: '/', delimiter: ':' };",
		"fs":                  "module.exports = { readFileSync: function() { return ''; }, existsSync: function() { return false; }, writeFileSync: function() {}, mkdirSync: function() {}, readdirSync: function() { return []; }, statSync: function() { return { isFile: function(){return false}, isDirectory: function(){return false} }; } };",
		"url":                 "module.exports = { parse: function(u) { return { href: u, pathname: u }; }, resolve: function(a, b) { return b; }, format: function(o) { return o.href || ''; }, URL: typeof URL !== 'undefined' ? URL : function(u) { this.href = u; this.pathname = u; } };",
		"string_decoder":      "function StringDecoder() {} StringDecoder.prototype.write = function(buf) { return typeof buf === 'string' ? buf : String(buf); }; StringDecoder.prototype.end = function() { return ''; }; module.exports = { StringDecoder: StringDecoder };",
		"net":                 "module.exports = { createServer: function() { return { listen: function(){}, on: function(){ return this; }, close: function(){} }; }, connect: function() { return { on: function(){ return this; }, write: function(){}, end: function(){} }; }, isIP: function() { return 0; } };",
		"tls":                 "module.exports = {};",
		"os":                  "module.exports = { platform: function(){ return 'linux'; }, cpus: function(){ return []; }, totalmem: function(){ return 0; }, freemem: function(){ return 0; }, homedir: function(){ return '/'; }, tmpdir: function(){ return '/tmp'; }, EOL: '\\n', hostname: function(){ return 'localhost'; } };",
		"zlib":                "module.exports = { createGzip: function(){return {}}, createGunzip: function(){return {}}, createDeflate: function(){return {}}, createInflate: function(){return {}} };",
		"http":                "module.exports = { createServer: function() { return { listen: function(){}, on: function(){ return this; } }; }, request: function() { return { on: function(){ return this; }, write: function(){}, end: function(){} }; }, get: function() { return { on: function(){ return this; } }; }, Agent: function(){}, METHODS: [], STATUS_CODES: {} };",
		"https":               "module.exports = { createServer: function() { return { listen: function(){}, on: function(){ return this; } }; }, request: function() { return { on: function(){ return this; }, write: function(){}, end: function(){} }; }, get: function() { return { on: function(){ return this; } }; }, Agent: function(){} };",
		"child_process":       "module.exports = { exec: function(){}, execSync: function(){return '';}, spawn: function(){ return { on: function(){ return this; }, stdout: { on: function(){ return this; } }, stderr: { on: function(){ return this; } } }; } };",
		"assert":              "module.exports = function(val) { if (!val) throw new Error('assertion failed'); }; module.exports.ok = module.exports; module.exports.equal = function(){}; module.exports.deepEqual = function(){}; module.exports.strictEqual = function(){}; module.exports.notEqual = function(){};",
		"querystring":         "module.exports = { parse: function(s) { var o = {}; if (!s) return o; s.split('&').forEach(function(p) { var kv = p.split('='); o[decodeURIComponent(kv[0])] = decodeURIComponent(kv[1] || ''); }); return o; }, stringify: function(o) { return Object.keys(o).map(function(k) { return encodeURIComponent(k) + '=' + encodeURIComponent(o[k]); }).join('&'); }, encode: function(o) { return this.stringify(o); }, decode: function(s) { return this.parse(s); } };",
		"punycode":            "module.exports = { encode: function(s){return s;}, decode: function(s){return s;}, toASCII: function(s){return s;}, toUnicode: function(s){return s;} };",
		"vm":                  "module.exports = { createContext: function(o){return o||{};}, runInContext: function(code){return eval(code);}, runInNewContext: function(code){return eval(code);}, Script: function(code){ this.runInContext = function(){return eval(code);}; } };",
		"worker_threads":      "module.exports = { isMainThread: true, parentPort: null, Worker: function(){}, workerData: null };",
		"perf_hooks":          "module.exports = { performance: { now: function(){ return Date.now(); }, mark: function(){}, measure: function(){} }, PerformanceObserver: function(){ this.observe = function(){}; } };",
		"diagnostics_channel": "module.exports = { channel: function(){ return { subscribe: function(){}, unsubscribe: function(){}, hasSubscribers: false }; }, hasSubscribers: function(){ return false; }, subscribe: function(){}, unsubscribe: function(){} };",
	}

	for mod, content := range shims {
		shimPath := filepath.Join(shimDir, mod+".js")
		os.WriteFile(shimPath, []byte(content), 0644)
	}

	// Build alias map: "util" -> "/path/to/_shims/util.js"
	aliases := make(map[string]string, len(shims))
	for mod := range shims {
		absShim, _ := filepath.Abs(filepath.Join(shimDir, mod+".js"))
		aliases[mod] = absShim
	}

	result := api.Build(api.BuildOptions{
		EntryPoints:      []string{entryPath},
		Bundle:           true,
		Write:            false,
		Platform:         api.PlatformNeutral,
		Format:           api.FormatIIFE,
		Target:           api.ES2020,
		JSX:              api.JSXAutomatic,
		Sourcemap:        api.SourceMapNone,
		MinifySyntax:     !b.cfg.Dev,
		MinifyWhitespace: !b.cfg.Dev,
		NodePaths:        []string{absNodeModules},
		Alias:            aliases,
		Define: map[string]string{
			"process.env.NODE_ENV": fmt.Sprintf(`"%s"`, b.envMode()),
		},
	})

	if len(result.Errors) > 0 {
		return "", fmt.Errorf("esbuild server: %s", result.Errors[0].Text)
	}

	return string(result.OutputFiles[0].Contents), nil
}

func (b *Bundler) buildClient(entries []string) (map[string]string, error) {
	clientDir := filepath.Join(b.cfg.BuildDir, "client")
	os.MkdirAll(clientDir, 0755)

	// Generate per-page hydration entry files.
	// Each page gets a tiny JS file that imports the component
	// and calls hydrateRoot. esbuild bundles each one separately.
	hydrateEntries, err := b.generateClientEntries(entries, clientDir)
	if err != nil {
		return nil, err
	}

	absNodeModules, _ := filepath.Abs("node_modules")

	result := api.Build(api.BuildOptions{
		EntryPoints:      hydrateEntries,
		Bundle:           true,
		Write:            true,
		Outdir:           clientDir,
		Platform:         api.PlatformBrowser,
		Format:           api.FormatESModule,
		Target:           api.ES2020,
		JSX:              api.JSXAutomatic,
		Splitting:        true,
		ChunkNames:       "chunks/[name]-[hash]",
		Sourcemap:        api.SourceMapLinked,
		MinifySyntax:     !b.cfg.Dev,
		MinifyWhitespace: !b.cfg.Dev,
		NodePaths:        []string{absNodeModules},
		Define: map[string]string{
			"process.env.NODE_ENV": fmt.Sprintf(`"%s"`, b.envMode()),
		},
	})

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("esbuild client: %s", result.Errors[0].Text)
	}

	// Map route -> output JS URL path
	clientMap := make(map[string]string)
	for _, entry := range entries {
		route := b.filePathToRoute(entry)
		// Hydrate entry mirrors page structure: pages/index.tsx -> _hydrate_index.js
		name := b.hydrateEntryName(entry)
		outFile := filepath.Join(clientDir, name+".js")
		if _, err := os.Stat(outFile); err == nil {
			clientMap[route] = outFile
		}
	}

	return clientMap, nil
}

// generateClientEntries creates per-page hydration scripts with SPA router.
// After initial SSR hydration, all subsequent navigation is client-side:
// 1. Intercept <a> clicks
// 2. Fetch page data from /__data?route=/path endpoint
// 3. Swap component without full reload
func (b *Bundler) generateClientEntries(entries []string, clientDir string) ([]string, error) {
	var hydrateFiles []string

	// Route map for client-side navigation
	var routeImports strings.Builder
	for _, entry := range entries {
		absPage, _ := filepath.Abs(entry)
		route := b.filePathToRoute(entry)
		routeImports.WriteString(fmt.Sprintf("  '%s': () => import('%s'),\n", route, absPage))
	}

	for _, entry := range entries {
		absPage, _ := filepath.Abs(entry)
		name := b.hydrateEntryName(entry)

		// SPA hydration script.
		// Key insight: initial hydration renders ONLY the page component
		// (matching server HTML exactly). SPA navigation activates after
		// hydration completes via useEffect (which only runs on client).
		script := fmt.Sprintf(`import { hydrateRoot } from 'react-dom/client';
import { createElement, useState, useEffect, useCallback, useRef } from 'react';
import InitialPage from '%s';

const routes = {
%s};

function routeToRegex(pattern) {
  const paramNames = [];
  const regexStr = pattern
    .split('/')
    .map(seg => {
      if (seg.startsWith(':')) {
        paramNames.push(seg.slice(1));
        return '([^/]+)';
      }
      return seg;
    })
    .join('/');
  return { regex: new RegExp('^' + regexStr + '$'), paramNames };
}

const compiledRoutes = Object.keys(routes).map(pattern => ({
  pattern,
  ...routeToRegex(pattern),
  load: routes[pattern]
}));

function matchRoute(path) {
  for (const route of compiledRoutes) {
    const match = path.match(route.regex);
    if (match) {
      const params = {};
      route.paramNames.forEach((name, i) => { params[name] = match[i + 1]; });
      return { pattern: route.pattern, params, load: route.load };
    }
  }
  return null;
}

async function fetchProps(path) {
  try {
    const res = await fetch('/__data?path=' + encodeURIComponent(path));
    if (!res.ok) return {};
    return await res.json();
  } catch (e) {
    console.error('props fetch failed:', e);
    return {};
  }
}

function SPAShell({ initialProps }) {
  const [currentPath, setPath] = useState(window.location.pathname);
  const [PageComponent, setPage] = useState(() => InitialPage);
  const [pageProps, setPageProps] = useState(initialProps);
  const [loading, setLoading] = useState(false);
  const hydrated = useRef(false);

  const navigate = useCallback(async (path, pushState) => {
    if (path === currentPath) return;

    const matched = matchRoute(path);
    if (!matched) {
      window.location.href = path;
      return;
    }

    setLoading(true);
    try {
      const [mod, props] = await Promise.all([
        matched.load(),
        fetchProps(path)
      ]);
      const Component = mod.default || mod;
      setPage(() => Component);
      setPageProps({ ...props, ...matched.params });
      setPath(path);
      if (pushState) window.history.pushState({ path }, '', path);
      window.scrollTo(0, 0);
    } catch (e) {
      console.error('navigation failed:', e);
      window.location.href = path;
    } finally {
      setLoading(false);
    }
  }, [currentPath]);

  // Activate SPA link interception AFTER hydration.
  // useEffect never runs on server, so this doesn't affect SSR match.
  useEffect(() => {
    hydrated.current = true;

    function handleClick(e) {
      const anchor = e.target.closest('a');
      if (!anchor) return;
      const href = anchor.getAttribute('href');
      if (!href || href.startsWith('http') || href.startsWith('//')) return;
      if (anchor.target === '_blank') return;
      if (href.startsWith('#')) return;
      if (e.metaKey || e.ctrlKey || e.shiftKey || e.altKey) return;
      e.preventDefault();
      navigate(href, true);
    }
    document.addEventListener('click', handleClick);
    return () => document.removeEventListener('click', handleClick);
  }, [navigate]);

  useEffect(() => {
    function handlePop() {
      navigate(window.location.pathname, false);
    }
    window.addEventListener('popstate', handlePop);
    return () => window.removeEventListener('popstate', handlePop);
  }, [navigate]);

  // First render: just the page with its props — matches server HTML exactly.
  // Loading bar only shows after hydration (client-only via useEffect state).
  return createElement('div', null,
    loading && hydrated.current
      ? createElement('div', {
          style: {
            position: 'fixed', top: 0, left: 0, width: '100%%',
            height: '3px', background: '#0070f3', zIndex: 9999,
            animation: 'reactgo_loading 1s ease-in-out infinite'
          }
        })
      : null,
    createElement(PageComponent, pageProps)
  );
}

// Inject loading animation
const style = document.createElement('style');
style.textContent = '@keyframes reactgo_loading { 0%% { transform: translateX(-100%%); } 50%% { transform: translateX(0%%); } 100%% { transform: translateX(100%%); } }';
document.head.appendChild(style);

// --- Hydration ---
// Initial render of SPAShell produces:
//   <div>{null}<InitialPage {...props} /></div>
// Which React renders as:
//   <div><InitialPage {...props} /></div>
//
// Server produces (via __renderToString):
//   <div>...InitialPage HTML...</div>
//
// These DON'T match because server doesn't have the outer <div>.
// Fix: server wraps in a plain <div> too.

const container = document.getElementById('__reactgo');
const initialProps = window.__REACTGO_DATA__ || {};
hydrateRoot(container, createElement(SPAShell, { initialProps }));
`, absPage, routeImports.String())

		hydratePath := filepath.Join(clientDir, name+".jsx")
		if err := os.WriteFile(hydratePath, []byte(script), 0644); err != nil {
			return nil, err
		}
		hydrateFiles = append(hydrateFiles, hydratePath)
	}

	return hydrateFiles, nil
}

// generateServerEntry creates the V8 entry that registers all pages
// and exposes __renderToString and __getServerSideProps globals.
func (b *Bundler) generateServerEntry(entries []string) string {
	var sb strings.Builder

	sb.WriteString(`var React = require('react');
var ReactDOMServer = require('react-dom/server');

var routes = {};
var propsLoaders = {};

`)

	for i, entry := range entries {
		absPath, _ := filepath.Abs(entry)
		route := b.filePathToRoute(entry)

		sb.WriteString(fmt.Sprintf("var Page%d = require('%s');\n", i, absPath))
		// Support both default and named exports
		sb.WriteString(fmt.Sprintf("var Comp%d = Page%d.default || Page%d;\n", i, i, i))
		sb.WriteString(fmt.Sprintf("routes['%s'] = Comp%d;\n", route, i))

		// Register getServerSideProps if exported
		sb.WriteString(fmt.Sprintf("if (Page%d.getServerSideProps) { propsLoaders['%s'] = Page%d.getServerSideProps; }\n\n", i, route, i))
	}

	// Global render bridge — called by Worker.Execute()
	sb.WriteString(`
globalThis.__renderToString = function(route, props) {
  var Component = routes[route];
  if (!Component) {
    return '<div>404 - Page not found</div>';
  }
  try {
    return ReactDOMServer.renderToString(React.createElement(Component, props));
  } catch(e) {
    return '<div>Render Error: ' + e.message + '</div>';
  }
};

globalThis.__getServerSideProps = function(route, context) {
  var loader = propsLoaders[route];
  if (!loader) {
    return JSON.stringify({ props: {} });
  }
  try {
    var result = loader(context);
    return JSON.stringify(result);
  } catch(e) {
    return JSON.stringify({ props: {}, error: e.message });
  }
};

globalThis.__hasServerProps = function(route) {
  return !!propsLoaders[route];
};
`)

	return sb.String()
}

func (b *Bundler) filePathToRoute(filePath string) string {
	route := strings.TrimPrefix(filePath, b.cfg.PagesDir)
	route = strings.TrimSuffix(route, filepath.Ext(route))
	route = filepath.ToSlash(route)
	route = strings.ReplaceAll(route, "[", ":")
	route = strings.ReplaceAll(route, "]", "")

	if strings.HasSuffix(route, "/index") {
		route = strings.TrimSuffix(route, "/index")
	}
	if route == "" {
		route = "/"
	}
	return route
}

// hydrateEntryName creates a unique flat filename for hydration entries.
// pages/posts/[id].tsx -> _hydrate_posts_id
func (b *Bundler) hydrateEntryName(filePath string) string {
	name := strings.TrimPrefix(filePath, b.cfg.PagesDir)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	name = filepath.ToSlash(name)
	name = strings.Trim(name, "/")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "[", "")
	name = strings.ReplaceAll(name, "]", "")
	return "_hydrate_" + name
}

func (b *Bundler) envMode() string {
	if b.cfg.Dev {
		return "development"
	}
	return "production"
}
