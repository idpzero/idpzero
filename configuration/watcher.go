package configuration

import (
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	w    *fsnotify.Watcher
	done chan struct{}
}

func NewWatcher(ci *ConfigInformation, changed func(x *IDPConfiguration)) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return nil, err
	}

	w := &Watcher{
		w:    watcher,
		done: make(chan struct{}),
	}

	// Start listening for events.
	go func() {
		defer func() {
			w.done <- struct{}{}
		}()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					if event.Name == ci.ConfigFilePath() {
						t, err := ci.Load()
						if err != nil {
							color.Red("Error loading config file from watch")
						}
						if t != nil {
							changed(t)
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				color.Red("Error occured during watch: %v", err)
			}
		}
	}()

	if err := watcher.Add(ci.dirPath); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Watcher) Close() {
	w.w.Close()
	<-w.done // wait for the go routine to finish
}
