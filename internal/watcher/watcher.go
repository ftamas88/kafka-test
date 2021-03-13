package watcher

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/ftamas88/kafka-test/internal/domain"
	"log"
	"os"
	"time"
)

type FileWatcher struct {
	incoming chan<- domain.Payload
	cloudIncoming chan<- domain.Payload
	notify <-chan domain.Payload
}

func NewFileWatcher(ch chan<- domain.Payload, cloudCh chan<- domain.Payload, notify <-chan domain.Payload) *FileWatcher {
	return &FileWatcher{
		incoming: ch,
		notify: notify,
		cloudIncoming: cloudCh,
	}
}


func (f FileWatcher) Run(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("unable to create watcher: %s", err.Error())
	}

	defer func() {
		_ = watcher.Close()
	}()

	done := make(chan bool)
	go func() {
		log.Println("watching folder..")
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write{
					// Due to the behavior of the mock generator which creates the file first with empty content
					// I cannot subscribe to the file created event here.

					r, err := os.ReadFile(event.Name)
					if err != nil {
						log.Printf("unable to read file: %s", err.Error())
						continue
					}
					log.Printf("File change detected: [%s]. Publishing [%s]", event.Name, event.Op)

					go func() {
						f.cloudIncoming <- domain.Payload{
							Data:     r,
							Status:   false,
							Filename: event.Name,
						}

						f.incoming <- domain.Payload{
							Data:     r,
							Status:   false,
							Filename: event.Name,
						}
					}()

				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	go func() {
		log.Println("waiting for file deletion event")
		for {
			select {
			case event, ok := <-f.notify:
				if !ok {
					return
				}

				if event.Status {
					<- time.After(1 * time.Second)
					fmt.Printf("deleting file: %s", event.Filename)
					if err := os.Remove(event.Filename); err != nil {
						log.Printf("unable to delete file: [%s]", event.Filename)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("tmp/ingest")
	if err != nil {
		return fmt.Errorf("unable to watch the folder: %s", err.Error())
	}
	<-done

	return nil
}