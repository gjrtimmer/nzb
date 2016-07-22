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

// ExtractFilename from subject
func (s Subject) ExtractFilename() (string, error) {
    match := subjectFilenameRegex.FindStringSubmatch(string(s))
    if match == nil {
        return "", fmt.Errorf("unable to extract filename from '%s'", string(s))
    }

    // Filename extracted
    return match[1], nil
}

// ExtractPartNumber from subject
// Example: ...[12/77]..., this will extract 12
// only if it is available within the subject
func (s Subject) ExtractPartNumber() (int, error) {
    match := subjectPartsRegex.FindStringSubmatch(string(s))
    if match == nil {
        return 0, fmt.Errorf("unable to extract part number from '%s'", string(s))
    }

    p, _ := strconv.Atoi(match[1])

    return p, nil
}

// ExtractTotalParts from subject
// Example: ...[12/77]..., this will extract 77
// only if it is available within the subject
func (s Subject) ExtractTotalParts() (int, error) {
    match := subjectPartsRegex.FindStringSubmatch(string(s))
    if match == nil {
        return 0, fmt.Errorf("unable to extract total parts from '%s'", string(s))
    }

    p, _ := strconv.Atoi(match[2])

    return p, nil
}

// ExtractYEncPartNumber from subject
// Example: ...(001/420)..., this will extract 001
// yEnc encoding parts can almost
// alwwys be found at the end of a subject string
func (s Subject) ExtractYEncPartNumber() (int, error) {
    match := subjectYEncPartsRegex.FindStringSubmatch(string(s))
    if match == nil {
        return 0, fmt.Errorf("unable to extract yEnc part from '%s'", string(s))
    }

    p, _ := strconv.Atoi(match[1])

    return p, nil
}

// ExtractYEncTotalParts from subject
// Example: ...(001/420)..., this will extract 420
// yEnc encoding parts can almost
// alwwys be found at the end of a subject string
func (s Subject) ExtractYEncTotalParts() (int, error) {
    match := subjectYEncPartsRegex.FindStringSubmatch(string(s))
    if match == nil {
        return 0, fmt.Errorf("unable to extract yEnc total parts from '%s'", string(s))
    }

    p, _ := strconv.Atoi(match[2])

    return p, nil
}

// EOF
