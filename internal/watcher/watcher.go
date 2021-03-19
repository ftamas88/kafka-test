package watcher

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/ftamas88/kafka-test/internal/domain"
)

type FileWatcher struct {
	incoming chan<- domain.Payload
	notify   <-chan domain.Payload
}

func NewFileWatcher(ch chan<- domain.Payload, notify <-chan domain.Payload) *FileWatcher {
	return &FileWatcher{
		incoming: ch,
		notify:   notify,
	}
}

func (f FileWatcher) Run(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("unable to create watcher: %s", err.Error())
		return fmt.Errorf("unable to create watcher: %s", err.Error())
	}

	defer func() {
		_ = watcher.Close()
	}()

	log.Println("File watcher running..")

	lastFile := ""

	done := make(chan bool)

	f.watchFileCreation(watcher, lastFile)
	f.watchFileDeletion(watcher)

	err = watcher.Add("tmp/ingest")
	if err != nil {
		return fmt.Errorf("unable to watchFileCreation the folder: %s", err.Error())
	}
	<-done

	log.Println("File watcher exiting..")

	return nil
}

func (f FileWatcher) watchFileDeletion(watcher *fsnotify.Watcher) {
	go func() {
		log.Println("waiting for file deletion event")
		for {
			select {
			case event, ok := <-f.notify:
				if !ok {
					return
				}

				fmt.Printf("deleting file: %s", event.Filename)
				if err := os.Remove(event.Filename); err != nil {
					log.Printf("unable to delete file: [%s]", event.Filename)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
}

func (f FileWatcher) watchFileCreation(watcher *fsnotify.Watcher, lastFile string) {
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					// Due to the behaviour of the mock generator which creates the file first with empty content
					// I cannot subscribe to the file created event here.
					if lastFile == event.Name {
						continue
					}
					lastFile = event.Name

					r, err := os.ReadFile(lastFile)
					if err != nil {
						log.Printf("unable to read file: %s", err.Error())
						continue
					}
					log.Printf("File change detected: [%s]. Publishing [%s]", lastFile, event.Op)

					go func(fName string) {
						f.incoming <- domain.Payload{
							Data:     r,
							Filename: fName,
						}
					}(lastFile)

				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
}
