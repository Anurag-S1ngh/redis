package resp

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func Parse(reader *bufio.Reader) ([]string, error) {
	prefix, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	if prefix != '*' {
		return nil, fmt.Errorf("invalid RESP format, expected * got %q", prefix)
	}

	argCountStr, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	argCount, err := strconv.Atoi(strings.TrimSpace(argCountStr))
	if err != nil {
		return nil, err
	}

	fmt.Println("argCount", argCount)

	parts := make([]string, 0, argCount)

	for range argCount {
		typ, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}

		fmt.Println("typ", string(typ))

		switch typ {
		case '$':
			str, err := parseBulk(reader)
			if err != nil {
				return nil, err
			}
			parts = append(parts, str)
		}

	}

	return parts, nil
}

func parseBulk(reader *bufio.Reader) (string, error) {
	argNumberStr, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	argNumber, err := strconv.Atoi(strings.TrimSpace(argNumberStr))
	if err != nil {
		return "", err
	}
	fmt.Println("argNumber", argNumber)

	if argNumber == 0 {
		if _, err := reader.ReadString('\n'); err != nil {
			return "", err
		}
		return "", nil
	}

	buf := make([]byte, argNumber)
	if _, err := reader.Read(buf); err != nil {
		return "", err
	}
	line := string(buf)

	if _, err := reader.ReadString('\n'); err != nil {
		return "", err
	}

	return line, nil
}
