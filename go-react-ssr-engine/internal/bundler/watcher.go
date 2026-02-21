package bundler

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/thutasann/go-react-ssr-engine/internal/config"
)

// Watcher monitors pages/ and components/ for changes and triggers rebuilds.
// Uses debouncing to coalesce rapid file saves (e.g. editor auto-save)
// into a single rebuild. Without debounce, saving 5 files in 100ms
// would trigger 5 rebuilds - wasteful since each one takes <50ms anyway.
type Watcher struct {
	cfg      *config.Config
	bundler  *Bundler
	onChange func(*BuildResult) // callback when rebuild completes

	mu    sync.Mutex
	timer *time.Timer
}

// NewWatcher creates a file watcher. onChange is called after every successful rebuild
// with the new build result. Typically used to call engine.LoadBundle()
func NewWatcher(cfg *config.Config, bundler *Bundler, onChange func(*BuildResult)) *Watcher {
	return &Watcher{
		cfg:      cfg,
		bundler:  bundler,
		onChange: onChange,
	}
}

// Start begins watching in a background goroutine.
// Returns immediately. Errors are logged, not returned -
// a watcher failer shouldn't crash the server.
func (w *Watcher) Start() error {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	// Watch pages and components directories recursively
	dirs := []string{w.cfg.PagesDir, "components"}
	for _, dir := range dirs {
		if err := w.addRecursive(fsw, dir); err != nil {
			log.Printf("watcher: skipping %s: %v", dir, err)
		}
	}

	go w.loop(fsw)

	log.Printf("watcher: watching %v for changes", dirs)
	return nil
}

// loop is the main event loop running in its own goroutine.
// Reads fsnotify events and debounce rebuilds.
func (w *Watcher) loop(fsw *fsnotify.Watcher) {
	defer fsw.Close()

	for {
		select {
		case event, ok := <-fsw.Events:
			if !ok {
				return
			}
			// Only care about JS/TS file changes
			if !isJSFile(event.Name) {
				continue
			}
			w.debounceRebuild()

		case err, ok := <-fsw.Errors:
			if !ok {
				return
			}
			log.Printf("watcher: error: %v", err)
		}
	}
}

// debounceRebuild resets a 100ms timer on every file event.
// The rebuild only fires when no events arrive for 100ms.
// This means rapid saves result in exactly one rebuild.
func (w *Watcher) debounceRebuild() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.timer != nil {
		w.timer.Stop()
	}

	w.timer = time.AfterFunc(100*time.Millisecond, func() {
		log.Println("watcher: rebuilding...")
		start := time.Now()

		result, err := w.bundler.Build()
		if err != nil {
			log.Printf("watcher: rebuild failed: %v", err)
			return
		}

		log.Printf("watcher: rebuilt in %s", time.Since(start))
		w.onChange(result)
	})
}

// addRecursive walks a directory and adds all subdirs to the watcher.
func (w *Watcher) addRecursive(fsw *fsnotify.Watcher, root string) error {
	return filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return fsw.Add(path)
		}
		return nil
	})
}

func isJSFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".tsx" || ext == ".jsx" || ext == ".ts" || ext == ".js"
}
