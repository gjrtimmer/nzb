package nzb

import (
    "fmt"
    "regexp"
    "strconv"
    "strings"
)

var (
    subjectName             = regexp.MustCompile(`^([^-]*)`)
    subjectFilenameRegex    = regexp.MustCompile(`"([^"]*)"`)
    subjectPartsRegex       = regexp.MustCompile(`\[([0-9]+)\/([0-9]+)\]`)
    subjectYEncPartsRegex   = regexp.MustCompile(`\(([0-9]+)\/([0-9]+)\)`)
)

// Subject of NZB File
type Subject string

func (s Subject) extractName() (string, error) {
    match := subjectName.FindStringSubmatch(string(s))
    if match == nil {
        return "", fmt.Errorf("unable to extract download name from '%s'", string(s))
    }

    name := strings.TrimSpace(match[1])

    return name, nil
}

// Subject
func (s Subject) extractFilename() (string, error) {
    match := subjectFilenameRegex.FindStringSubmatch(string(s))
    if match == nil {
        return "", fmt.Errorf("unable to extract filename from '%s'", string(s))
    }

    // Filename extracted
    return match[1], nil
}

func (s Subject) extractPartNumber() (int, error) {
    match := subjectPartsRegex.FindStringSubmatch(string(s))
    if match == nil {
        return 0, fmt.Errorf("unable to extract part number from '%s'", string(s))
    }

    p, _ := strconv.Atoi(match[1])

    return p, nil
}

func (s Subject) extractTotalParts() (int, error) {
    match := subjectPartsRegex.FindStringSubmatch(string(s))
    if match == nil {
        return 0, fmt.Errorf("unable to extract total parts from '%s'", string(s))
    }

    p, _ := strconv.Atoi(match[2])

    return p, nil
}

func (s Subject) extractYEncPartNumber() (int, error) {
    match := subjectYEncPartsRegex.FindStringSubmatch(string(s))
    if match == nil {
        return 0, fmt.Errorf("unable to extract yEnc part from '%s'", string(s))
    }

    p, _ := strconv.Atoi(match[1])

    return p, nil
}

func (s Subject) extractYEncTotalParts() (int, error) {
    match := subjectYEncPartsRegex.FindStringSubmatch(string(s))
    if match == nil {
        return 0, fmt.Errorf("unable to extract yEnc total parts from '%s'", string(s))
    }

    p, _ := strconv.Atoi(match[2])

    return p, nil
}

// EOF
