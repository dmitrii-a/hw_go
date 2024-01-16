package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/mailru/easyjson"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	stat, err := getStat(r, domain)
	if err != nil {
		return nil, fmt.Errorf("get error: %w", err)
	}
	return stat, nil
}

func getStat(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	result := make(DomainStat)
	for scanner.Scan() {
		var user User
		if err := easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}
		if len(strings.Split(user.Email, ".")) < 2 {
			return nil, fmt.Errorf("invalid email(without domain): %s", user.Email)
		}
		if strings.HasSuffix(user.Email, domain) {
			emailParts := strings.SplitN(user.Email, "@", 2)
			if len(emailParts) < 2 {
				return nil, fmt.Errorf("invalid email(without @): %s", user.Email)
			}
			result[strings.ToLower(emailParts[1])]++
		}
	}
	return result, nil
}
