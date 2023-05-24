package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var input = bufio.NewReader(os.Stdin)

func GetString(prompt string) (string, error) {
	fmt.Print(prompt)
	str, err := input.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(str), nil
}

func GetInt(prompt string) (int, error) {
	inp, err := GetString(prompt)

	if err != nil {
		return 0, err
	}

	ret, err := strconv.Atoi(inp)

	if err != nil {
		return 0, err
	}

	return ret, nil
}
