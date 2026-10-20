package hw10programoptimization

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type DomainStat map[string]int

type user struct {
	Email string `json:"email"`
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	domain = strings.ToLower(strings.TrimPrefix(domain, "."))
	if domain == "" {
		return result, nil
	}

	suffix := "." + domain
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}

		var u user
		if err := jsoniter.ConfigFastest.Unmarshal(line, &u); err != nil {
			return nil, fmt.Errorf("decode user: %w", err)
		}

		email := strings.ToLower(u.Email)
		at := strings.LastIndexByte(email, '@')
		if at < 0 || at == len(email)-1 {
			continue
		}

		host := email[at+1:]
		if strings.HasSuffix(host, suffix) {
			result[host]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read users: %w", err)
	}

	return result, nil
}
