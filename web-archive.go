package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

func main() {
	var subdomainListFile string
	var outputFileName string

	flag.StringVar(&subdomainListFile, "l", "", "Specify the subdomain list file")
	flag.StringVar(&outputFileName, "o", "web_archive_urls.txt", "Specify the output filename")
	flag.Parse()

	if subdomainListFile == "" {
		fmt.Println("Please provide a subdomain list file using the -l flag.")
		return
	}

	subdomains, err := readSubdomains(subdomainListFile)
	if err != nil {
		fmt.Println("Error reading subdomain list:", err)
		return
	}

	if len(subdomains) == 0 {
		fmt.Println("No subdomains found in the provided list.")
		return
	}

	// Create a progress bar
	bar := pb.StartNew(len(subdomains))
	defer bar.Finish()

	// Fetch URLs and write to the output file
	if err := fetchAndWriteURLs(subdomains, outputFileName, bar); err != nil {
		fmt.Println("Error fetching and writing URLs:", err)
		return
	}

	fmt.Printf("URLs written to %s\n", outputFileName)
}

func readSubdomains(filename string) ([]string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	subdomains := strings.Fields(string(data))
	return subdomains, nil
}

func fetchAndWriteURLs(subdomains []string, filename string, bar *pb.ProgressBar) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, sub := range subdomains {
		bar.Increment()

		url := fmt.Sprintf("https://web.archive.org/cdx/search/cdx?url=%s/*&output=text&fl=original&collapse=urlkey", sub)
		response, err := http.Get(url)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		lines := strings.Split(string(body), "\n")
		for _, line := range lines {
			if strings.Contains(line, sub) {
				fmt.Fprintf(file, "%s\n", line)
			}
		}
	}

	return nil
}
