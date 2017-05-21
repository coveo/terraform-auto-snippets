package utils

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"sync"
)

// GetDocument returns a goquery.Document that could be used to parse
// the content of an HTML file
func getDocument(uri url.URL) (result *goquery.Document, err error) {
	response, err := http.Get(uri.String())
	switch {
	case response == nil:
		return
	case response.StatusCode >= 400:
		err = fmt.Errorf("%s -▶ %s", uri.String(), response.Status)
	}
	if err == nil {
		result, err = goquery.NewDocumentFromResponse(response)
	}
	return
}

// StartDocumentWorkers starts the specified number of threads to handle tasks
func StartDocumentWorkers(nbWorkers int) (wf *WorkForce) {
	Assert(nbWorkers > 0, "The number of workers must be greater than 0")

	wf = &WorkForce{
		terminate: make(chan interface{}),
		task:      make(chan workerTask, nbWorkers),
	}

	for i := 0; i < nbWorkers; i++ {
		go func(workerId int) {
			defer wf.waitGroup.Done()

			for {
				select {
				case task := <-wf.task:
					// Received a task to process
					if doc, err := getDocument(task.uri); err == nil {
						// Defer the processing to the handler
						go task.handler(doc, task.result)
					} else {
						task.result <- err
					}
				case _ = <-wf.terminate:
					// Received a signal to terminate the task
					return
				}
			}
		}(i)
	}
	wf.waitGroup.Add(nbWorkers)
	return
}

// DocumentHandler is a function that could
type DocumentHandler func(*goquery.Document, chan error)

type workerTask struct {
	uri     url.URL
	handler DocumentHandler
	result  chan error
}

// IWorkForce represents the minimal interface for task force objects
type IWorkForce interface {
	ProcessDocument(url.URL, DocumentHandler) error
	TerminateAll()
}

// WorkForce handles a dedicated work force to process jobs
type WorkForce struct {
	task      chan workerTask
	terminate chan interface{}
	waitGroup sync.WaitGroup
}

// ProcessDocument add a task to be processed by the work force
func (wf *WorkForce) ProcessDocument(uri url.URL, handler DocumentHandler) error {
	result := make(chan error)
	wf.task <- workerTask{uri, handler, result}
	return <-result
}

// TerminateAll send the order to worker task to end their processing
func (wf *WorkForce) TerminateAll() {
	for i := 0; i < cap(wf.task); i++ {
		wf.terminate <- nil
	}
	wf.waitGroup.Wait()
}
