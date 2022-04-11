/* service.go is a local http server that runs as a daemon
   and accepts requests to create new blocks in Roam Reseach

   example usage:

     curl -d '{"note": "foo bar"}' \
	      -H 'Accept: text/json' \
		  -X POST http://localhost:8080
 */
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Queue struct {
	Notes []string
	Mutex sync.Mutex
}

func (q *Queue) Enqueue(note string) {
	q.Mutex.Lock()
	q.Notes = append(q.Notes, note)
	q.Mutex.Unlock()
}

func (q *Queue) Dequeue() string {
	var note string
	q.Mutex.Lock()
	note = q.Notes[0]
	q.Notes = q.Notes[1:]
	q.Mutex.Unlock()
	return note
}

var q Queue

type Block struct {
	Note string `json:"note"`
}

func handler(res http.ResponseWriter, req *http.Request) {
	var block Block
	json.NewDecoder(req.Body).Decode(&block)
	q.Enqueue(block.Note)
	fmt.Printf("queue: '%v'\n", q.Notes)
}

// Each post to the Roam Private API is really slow
func dumpNotesIntoCloud() {
	for {
		if len(q.Notes) > 0 {
			note := q.Dequeue()
			fmt.Printf("posting note: '%s' to roam\n", note)
		}
	}
}

func main() {
	go dumpNotesIntoCloud()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
