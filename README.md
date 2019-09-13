# Golang Short Domain Finder

This Golang program will attempt to find available short domains based on a few input parameters such as domain extension and maximum length of domain name

This is usefull is you're trying to find short and unique domains.
Particularly usefull for url shortener websites or finding a free subdomain with few letters.

This was originally created to find a short `.me` domain as a student since namecheap offers those.


## Requirements

Uses the [whois.go](https://github.com/likexian/whois-go) module.

To install:

`go get -u github.com/likexian/whois-go`


## Usage

```
Usage of ./short-domain-finder:
  -exts string
        List of domain extensions (ie. .com, .io) (default "tk,ml,ga,cf")
  -len int
        Maximum length of domain name (default 3)
  -sep string
        Char used to separate the list of domain extensions (default ",")
  -workers int
        Number of worker to query whois in parallel. Too many may overwhelm the service and get you blocked (default 10)
```

The error messages are printed to stderr so you can simply redirect valid names to files as such:

```
./short-domain-finder > output.txt
```


## TODO

- Some domains don't seem to be supported such as gq with current whois server
    - Potential fix: use a list of whois servers such as `https://github.com/whois-server-list/whois-server-list`
