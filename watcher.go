package main

import (
	"regexp"
	"time"

	"github.com/radovskyb/watcher"
	log "github.com/sirupsen/logrus"
)

type Watcher struct {
	uploader    *Uploader
	pollWatcher *watcher.Watcher
}

func NewWatcher(uploader *Uploader) *Watcher {
	w := new(Watcher)
	w.uploader = uploader

	w.pollWatcher = watcher.New()
	w.pollWatcher.FilterOps(watcher.Create, watcher.Rename, watcher.Move)

	return w
}

func (w Watcher) Watch(watchPath string, interval time.Duration, filterRegex string) error {
	// first upload everything in this directory
	w.uploader.uploadAll(watchPath, true)

	// watch directory recursively for changes
	if err := w.pollWatcher.AddRecursive(watchPath); err != nil {
		return err
	}

	if len(filterRegex) > 0 {
		r := regexp.MustCompile(filterRegex)
		w.pollWatcher.AddFilterHook(watcher.RegexFilterHook(r, false))
	}

	go w.processOps()

	// start the watching process
	log.Infof("Watching: %v\n", watchPath)
	if err := w.pollWatcher.Start(interval); err != nil {
		return err
	}
	return nil
}

func (w Watcher) processOps() {
	for {
		select {
		case event := <-w.pollWatcher.Event:
			log.Infof("action: found, path: '%s', op: '%s'", event.Path, event.Op)
			if !event.IsDir() {
				w.uploader.uploadAll(event.Path, true)
			}
		case err := <-w.pollWatcher.Error:
			log.Fatalln(err)
		case <-w.pollWatcher.Closed:
			return
		}
	}
}
