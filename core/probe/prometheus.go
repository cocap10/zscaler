package probe

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Prometheus probe
//
// Metrics are retrived using HTTP/Text protocol
// TODO protobuf support ?
type Prometheus struct {
	URL string // Sample: http://localhost:9100/metrics
	Key string // Sample: node_cpu{cpu="cpu6",mode="idle"}
}

// Name of the probe
func (p Prometheus) Name() string {
	return "Prometheus probe for " + p.URL + " [" + p.Key + "]"
}

// Value make the request and parse content
func (p Prometheus) Value() float64 {
	resp, err := http.Get(p.URL)
	if err != nil {
		fmt.Printf("%s", err)
		return -1.0
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("%s", err)
		return -1.0
	}
	fvalue, err := p.findValue(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		return -1.0
	}
	return fvalue
}

// find the matching token and parse probe value
func (p Prometheus) findValue(body io.Reader) (float64, error) {
	scanLines := bufio.NewScanner(body)
	for scanLines.Scan() { // for each line
		// now we're splitting by spaces
		scanWords := bufio.NewScanner(strings.NewReader(scanLines.Text()))
		scanWords.Split(bufio.ScanWords)

		if scanWords.Scan() { // line is not empty
			switch scanWords.Text() {
			case p.Key:
				if scanWords.Scan() {
					fvalue, err := strconv.ParseFloat(scanWords.Text(), 64)
					if err != nil {
						fmt.Printf("%s", err)
						return -1.0, err
					}
					return fvalue, nil
				}
			default:
				// if it doesn't start with the key, discard the line
				break
			}

		}
	}
	return 0, errors.New("Token " + p.Key + " not found")
}