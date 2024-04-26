package anthem

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"strings"
)

// paths is a global var to hold all the file locations that match our criteria
var paths = make(map[string]struct{})

// ProcessIndex is the main entrypoint, and expects an index.json.gz file as input
func ProcessIndex(indexPath string, outputPath string) {
	file, err := os.Open(indexPath)
	if err != nil {
		log.Fatalf("Unable to open file: %v", err)
	}
	defer file.Close()

	gr, err := gzip.NewReader(file)
	if err != nil {
		log.Fatalf("Unable to create gzip reader: %v", err)
	}

	reader := bufio.NewReader(gr)

	for {
		line, err := reader.ReadBytes('\n')
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			log.Printf("Completed")
			break
		}
		if err != nil {
			log.Fatalf("unable to read bytes: %v", err)
		}
		// Trim trailing characters off the json object
		line = bytes.TrimRight(line, ",\n")

		var r ReportingStructure
		err = json.Unmarshal(line, &r)
		// We expect the first few lines of the file to be un-parseable, as they are not the json structs we are looking for
		if err != nil {
			log.Printf("Unable to parse line: %v\n", err)
			continue
		}

		processRS(r)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Unable to create output file: %v", err)
	}

	// Print out de-duped set of paths
	for path, _ := range paths {
		_, err := outFile.WriteString(path + "\n")
		if err != nil {
			log.Printf("Error writing path to file: %v", err)
		}
	}
}

// processRS takes in a ReportingStructure and adds all file locations that match our filter to paths
func processRS(r ReportingStructure) {
	for _, file := range r.NetworkFiles {
		// If we have seen this path before, just continue
		_, ok := paths[file.Location]
		if ok {
			continue
		}

		// Heuristic used to identify Anthem NY PPO plans
		if file.Description == "In-Network Negotiated Rates Files" && strings.Contains(file.Location, "NY_") {
			paths[file.Location] = struct{}{}
		}
	}
}
