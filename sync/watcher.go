package sync

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"go-minio-sync/config"
)

func StartFileWatcher(cfg *config.Config, callback func(event fsnotify.Event)) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		defer func() {
			_ = watcher.Close()
		}()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Printf("Events 出现错误，退出监听")
					return
				}
				time.Sleep(time.Duration(cfg.Watch.Delay) * time.Second)
				callback(event)
			case err, ok := <-watcher.Errors:
				if !ok {
					log.Printf("Errors 出现错误，退出监听")
					return
				}
				log.Printf("监听错误: %v", err)
			}
		}
	}()

	return watcher.Add(cfg.Watch.Dir)
}
