/* service.go is a local http server that runs as a daemon
   and accepts requests to create new blocks in Roam Reseach

   Roam Research is a note talking application. It markets itself as a tool for thought.
   I especially love Roam's queries, inline calculations, it's bi-directional linking,
   embedded code formatting and JS/CLJ execution, the LaTeX math formatting. I like roam.

   The private roam-api is a node package that needs your graph name, user name and password
   to connect. I set those in the Dockerfile.dapper and then pass the password in as an arg,
   using an environment variable to keep it out of the git repo.

   The reason I use Dapper to build roam-api is that I want to keep all the node.js stuff
   in a container.

   example usage:

     curl -d '{"note": "foo bar"}' \
	      -H 'Accept: text/json' \
		  -X POST http://localhost:8080
*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	cmdchain "github.com/rainu/go-command-chain"
)

const dapperDir = "/tmp/roamd/"
const dapperFilePath = "/tmp/roamd/Dockerfile.dapper"
const dapperFileTemplate = `FROM mhart/alpine-node:16.4.2
ENV ROAM_API_GRAPH=
ENV ROAM_API_EMAIL=
ENV ROAM_API_PASSWORD=
RUN npm i -g roam-research-private-api`

func main() {
	msgs := make(chan string, 10)
	handler := makeHandler(msgs)
	checkRoamApiDir()
	http.HandleFunc("/", handler)
	// Listen for notes in the background on port 8080
	go func() {
		http.ListenAndServe(":8080", nil)
	}()

	// Receive from the messages channel and print them out
	for {
		msg := <-msgs
		fmt.Printf("received message: %s\n", msg)
	}
}

func checkRoamApiDir() {
	statDir, err := os.Stat(dapperDir)
	if err != nil {
		os.Mkdir(dapperDir, os.ModeDir)
	}
	statDir, err = os.Stat(dapperDir)
	if err != nil {
		panic(err) // tried to recover from missing dir, failed
	}
	if statDir.IsDir() {
		// check if Dapper file is in directory
		_, err := os.Stat(dapperFilePath)
		if err != nil {
			// create the template file then error out to the human
			file, err := os.Create(dapperFilePath)
			if err == nil {
				defer file.Close()
				// write the template to the dapper file
				n, err := file.WriteString(dapperFileTemplate)
				if err != nil {
					panic(err)
				}
				if n != len(dapperFileTemplate) {
					panic(fmt.Sprintf("Should have written %d bytes to %s",
						n, dapperFilePath,
					))
				}
				// This is the point where we have successfully written /tmp/roamd/Dockerfile.dapper
				// Now we alert the user of the file contents, and to edit the file!
				fmt.Printf("PLEASE EDIT /tmp/roamd/Dockerfile.dapper to include your graph, username and password\n")
			} else {
				// file creation failed, time to panic
				panic(err)
			}
		}
	}
}

type Block struct {
	Note string `json:"note"`
}

func makeHandler(msgs chan<- string) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		var block Block
		json.NewDecoder(req.Body).Decode(&block)
		fmt.Printf("decoded request: %v\n", block)
		go roamApiCreateBlock(msgs, block.Note)
	}
}

// Each post to the Roam Private API is really slow, so we will put it
// into this background process called service. Then we use a goroutine
// and a channel to allow service to handle many requests quickly
func roamApiCreateBlock(msgs chan<- string, note string) error {
	cmdOut := &bytes.Buffer{}
	cmdErr := &bytes.Buffer{}

	create := func() error {
		err := cmdchain.Builder().
			Join(
				"dapper",
				"--directory", "/tmp/roamd/",
				"roam-api", "create", note,
			).
			Finalize().
			WithOutput(cmdOut).
			WithError(cmdErr).
			Run()
		msgs <- cmdOut.String()
		msgs <- cmdErr.String()
		return err
	}

	retries := 0
	for err := create(); err != nil; retries++ {
		msgs <- err.Error()
		if retries > 1 {
			msgs <- "too many retries"
			return nil
		}
	}
	return nil
}
