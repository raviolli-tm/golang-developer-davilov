package hw10programoptimization

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
)

type DomainStat map[string]int

var domainRegex = regexp.MustCompile(`@[a-zA-Z0-9\.]*`)

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	byteDomain := []byte(domain)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Bytes()
		domainFind := domainRegex.Find(line)
		if bytes.HasSuffix(domainFind, byteDomain) {
			res := bytes.ToLower(domainFind[1:])
			result[string(res)]++
		}
	}

	return result, nil
}
