package service

import (
	"errors"
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func AutoDetection(s string) (string, error) {
	emptyString := errors.New("input string is empty")
	if len(s) == 0 {
		return "", emptyString
	}
	if strings.ContainsAny(s, ".-") {
		return morse.ToText(s), nil
	}
	return morse.ToMorse(s), nil
}
