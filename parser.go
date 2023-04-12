package wingetsvc

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
)

var (
	searchHeaderRe, _ = regexp.Compile(`^.*(Name\s+)(Id\s+)(Version\s+)(Match\s+)(Source.*)$`)
	separatorRe, _    = regexp.Compile("^-+$")
)

func parseSearchOutput(output []byte) ([]ServiceInfo, error) {
	scanner := bufio.NewScanner(bytes.NewReader(output))
	var match []string
	for scanner.Scan() {
		if match = searchHeaderRe.FindStringSubmatch(scanner.Text()); match != nil {
			break
		}
	}

	if match == nil {
		// no items found
		return []ServiceInfo{}, nil
	}

	stops := []int{
		len(match[1]),
		len(match[1]) + len(match[2]),
		len(match[1]) + len(match[2]) + len(match[3]),
	}
	records := make([]ServiceInfo, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if separatorRe.MatchString(line) {
			// ignore separator
			continue
		}
		runes := []rune(line)
		info := ServiceInfo{
			Name:    strings.TrimSpace(string(runes[:stops[0]])),
			Id:      strings.TrimSpace(string(runes[stops[0]:stops[1]])),
			Version: strings.TrimSpace(string(runes[stops[1]:stops[2]])),
		}
		records = append(records, info)
	}
	return records, nil
}

func parseVersionsOutput(output []byte) ([]string, error) {
	scanner := bufio.NewScanner(bytes.NewReader(output))
	var records []string
	startRead := false
	for scanner.Scan() {
		line := scanner.Text()
		if separatorRe.MatchString(line) {
			startRead = true
			continue
		}

		if startRead {
			records = append(records, strings.TrimSpace(line))
		}
	}
	return records, nil
}
