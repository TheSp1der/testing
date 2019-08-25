package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

// check to see if an ip is in a network range
func isIPinNet(ip string, network string) bool {
	_, ipnet, _ := net.ParseCIDR(network)

	if ipnet.Contains(net.ParseIP(ip)) {
		return true
	}

	return false
}

// getEnvString returns string from environment variable.
func getEnvString(env string, def string) string {
	val := os.Getenv(env)
	if len(val) == 0 {
		return def
	}
	return val
}

// getEnvBool returns boolean from environment variable.
func getEnvBool(env string, def bool) bool {
	var (
		err error
		val = os.Getenv(env)
		ret bool
	)

	if len(val) == 0 {
		return def
	}

	if ret, err = strconv.ParseBool(val); err != nil {
		log.Fatal(val + " environment variable is not boolean")
	}

	return ret
}

// getEnvInt returns int from environment variable.
func getEnvInt(env string, def int) int {
	var (
		err error
		val = os.Getenv(env)
		ret int
	)

	if len(val) == 0 {
		return def
	}

	if ret, err = strconv.Atoi(val); err != nil {
		log.Fatal(env + " environment variable is not numeric")
	}

	return ret
}

var (
	filename string
	follow   bool
)

func init() {
	// read command line options
	cmdLnFile := flag.String("file", getEnvString("FILE", ""), "(FILE)\nFile to tail.")
	cmdLnFollow := flag.Bool("follow", getEnvBool("FOLLOW", true), "(FOLLOW)\nFollow file if rotated.")
	flag.Parse()

	// set global variables
	filename = *cmdLnFile
	follow = *cmdLnFollow

	if len(filename) == 0 || filename == "" {
		flag.PrintDefaults()
		log.Fatal("No file to follow defined.")
	}
}

func main() {
	readFile()
}

func readFile() {
	// initialize output parser
	logParser := make(chan string)
	go parseLine(logParser)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	for {
		// read the next line
		line, err := reader.ReadString('\n')
		if err != nil {
			// have we reached the end of the file
			if err == io.EOF {
				// wait for new data
				time.Sleep(250 * time.Millisecond)

				// get the current file position
				curReadPos, err := file.Seek(0, 1)
				if err != nil {
					break
				}

				// get current file size
				fileInfo, err := file.Stat()
				if err != nil {
					break
				}

				// check to see if the file was truncated
				if fileInfo.Size() < curReadPos {
					// file was truncated
					break
				}

				// see if we can determine if the file moved
				f, err := os.Stat(filename)
				if err != nil {
					break
				}
				if !os.SameFile(f, fileInfo) {
					// file location has changed
					break
				}
			} else {
				break
			}
		}
		logParser <- line
	}
}

func parseLine(l <-chan string) {
	for {
		rx := regexp.MustCompile(`^(?i)(([0-9]|[a-z]|\-)+)\s(([0-9]|\:|\.)+)\s([a-z]+\:)\s([a-z]+\:)\s([a-z]+)\s(@0x[0-9a-f]+)\s(([0-9]{1,3}\.){3}[0-9]{1,3})\#[0-9]+\s\((([a-z0-9]|_|\-|\.)+)\)`)
		fields := rx.FindStringSubmatch(<-l)
		if len(fields) < 12 {
			continue
		}

		date := fields[1]
		time := fields[3]
		reqIP := fields[9]
		query := fields[11]

		rx = regexp.MustCompile(`test-chamber-13\.lan`)

		if !isIPinNet(reqIP, "127.0.0.1/32") && !rx.MatchString(query) {
			reqHost, err := net.LookupAddr(reqIP)
			if err != nil {
				fmt.Println(color.MagentaString(date+" "+time) + " " + color.YellowString(reqIP) + " " + color.GreenString(query))
			} else {
				fmt.Println(color.MagentaString(date+" "+time) + " " + color.YellowString(strings.Join(reqHost, "")) + " " + color.GreenString(query))
			}
		}
	}
}
