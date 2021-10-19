package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/likexian/whois"
	"os"
	"regexp"
	"strings"
)

const (
	UNREGISTERED_DOMAIN_REGEX = "^No match|^NOT FOUND|^Not fo|AVAILABLE|^No Data Fou|has not been regi|No entri|^Invalid query or domain name not known in"
	LEN_ALPHABET              = 26
)

func whoisDomainLookup(domain string) (string, error) {
	if domain == "" {
		return "", errors.New("No domain provided to whoisDomainLookup function")
	}

	result, err := whois.Whois(domain)
	if err != nil {
		return "", err
	}

	return result, nil
}

// check if the domain's whois page matches a "NOT REGISTERED" or equivalent
func isDomainAvailable(domain string) bool {
	whoisResult, err := whoisDomainLookup(domain)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}

	isAvailable, err := regexp.MatchString(UNREGISTERED_DOMAIN_REGEX, whoisResult)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}

	return isAvailable
}

func getAlphabet() []byte {
	p := make([]byte, LEN_ALPHABET)
	for i := range p {
		p[i] = 'a' + byte(i)
	}
	return p
}

func getPermutations(b byte, alphabet []byte) []string {
	r := make([]string, len(alphabet))

	for i, v := range alphabet {
		r[i] = strings.Join([]string{string(b), string(v)}, "")
	}

	return r
}

func outputer(resultChan chan string, completed chan bool) {
	for {
		domain, ok := <-resultChan

		if !ok {
			completed <- true
			return
		}

		fmt.Println(domain)
	}
}

func worker(workerChan, resultChan chan string, completed chan bool) {
	for {
		domain, ok := <-workerChan

		if !ok {
			completed <- true
			return
		}

		isAvailable := isDomainAvailable(domain)
		if isAvailable {
			resultChan <- domain
		}
	}
}

func dispatcher(workerChan chan string, exts []string, maxLen int) {
	alphabet := getAlphabet()

	var recurse func(domain string, depth int)
	recurse = func(domain string, depth int) {
		if depth == 0 {
			return
		}

		for _, l := range alphabet {
			for _, ext := range exts {
				newDomain := strings.Join([]string{domain, string(l)}, "")
				workerChan <- strings.Join([]string{newDomain, ext}, ".")
				recurse(newDomain, depth-1)
			}
		}
	}
	recurse("", maxLen)

	// close the workerChan so the worker can know it's final job
	close(workerChan)
}

func main() {
	extsFlag := flag.String("exts", "tk,ml,cf", "List of domain extensions (ie. .com, .io)") // ga, gq is not supported for now...
	lenFlag := flag.Int("len", 3, "Maximum length of domain name")
	sepFlag := flag.String("sep", ",", "Char used to separate the list of domain extensions")
	workersFlag := flag.Int("workers", 10, "Number of worker to query whois in parallel. Too many may overwhelm the service and get you blocked")
	flag.Parse()

	exts := strings.Split(*extsFlag, *sepFlag)
	maxLen := *lenFlag
	numOfWorkers := *workersFlag

	workerChan := make(chan string, 3*numOfWorkers)
	resultChan := make(chan string, 3*numOfWorkers)
	completed := make(chan bool, 10)

	go dispatcher(workerChan, exts, maxLen)
	go outputer(resultChan, completed)

	for i := 0; i < numOfWorkers; i++ {
		go worker(workerChan, resultChan, completed)
	}

	for i := 0; i < numOfWorkers; i++ {
		<-completed
	}

	close(resultChan)
	// wait until outputerService finishes
	<-completed
}
