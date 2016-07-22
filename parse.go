package nzb

import (
    "encoding/xml"
    "fmt"
    "html"
    "io"
    "regexp"
    "strings"
    "strconv"

    "golang.org/x/net/html/charset"
)

// Credits for Regexp Goto: DanielMorsing (https://github.com/DanielMorsing)
var par2Filter = regexp.MustCompile(`(?i)(\.vol\d+\+(\d+))?\.par2$`)
// RegEx:
// Group 0: Prefix
// Group 1: Par2 Volume
// Group 2: Par2 Blocks

// Parse NZB file
func Parse(r io.Reader) (*NZB, error) {

    var err error

    // Setup Decoder
    decoder := xml.NewDecoder(r)
    decoder.CharsetReader = charset.NewReaderLabel

    // Create empty Xml NZB Struct
    x := new(xNZB)

    // Create new NZB
    n := new(NZB)

    // Decode NZB
    err = decoder.Decode(x)
    if err != nil {
        return nil, err
    }

    // Files marked processed into fileset
    processed := make(map[*File]bool)

    for _, f := range x.Files {

        // Extract filename
        f.Filename, err = f.Subject.ExtractFilename()
        if err != nil {
            return nil, err
        }

        // Calculate Filesize
        for _, s := range f.Segments {
            f.Size += Size(s.Bytes)
            s.ID = html.UnescapeString(s.ID)
        }

        // Calculate total size
        n.Size += Size(f.Size)

        // Search for par2 files
        pEx := par2Filter.FindStringSubmatch(f.Filename)
        if pEx == nil {
            // Current file is not a repair (par2) file
            continue
        }

        // Search for Par files
        // Check for RegEx Group 1 (.vol**+**)
        if pEx[1] != "" {
            blockstr := pEx[2]
            b, err := strconv.Atoi(blockstr)
            if err != nil {
                // Unable to extract block number from filename
                continue
            }

            prefix := f.Filename[:len(f.Filename) - len(pEx[0])]
            fileset := n.getFileSetByPrefix(prefix)
            parfile := f.convert2ParFile(b)
            fileset.ParSet.addParFile(parfile)

            processed[f] = true
        } else {
            // pEx[1] Group1 (.vol**+**) is empty
            // While there are matches
            // This means Group 0 (Prefix) has a match which
            // end whith the extension .par2
            // This file is expected to be the parent par2

            // Some statements are identical to above 'if'
            // reason == it's possible to have the first file
            // be processed within the NZB to be the parent par2

            // Get prefix
            prefix := f.Filename[:len(f.Filename) - len(pEx[0])]
            fileset := n.getFileSetByPrefix(prefix)
            fileset.ParSet.Parent = f

            // Mark file as processed
            processed[f] = true
        }
    }

    // Par2 files have been processed
    // result there are (n) FileSet(s) present with prefix (name) defined
    // this can now be used to process the regular files
    for _, f := range x.Files {

        // Search for par2 files
        // The previous used regex can now be used to get
        // all the files which are not par2 files
        // regardless of their extension or subformat like:
        // - .Part***.
        // - .nfo
        // - .sfv
        // - .r**
        pEx := par2Filter.FindStringSubmatch(f.Filename)
        if pEx == nil {
            // file is a regular file !par2

            // Foreach file (f) loop the previous
            // created FileSet(s)
            // check the prefix (name) of the FileSet
            // against the Prefix part of the current file.Filename
            // to decide if current file belongs to a fileset
            for _, fs := range n.FileSets {

                // Check fileset name (previous generated prefix)
                // against first part of filename
                if strings.HasPrefix(f.Filename, fs.Name) {
                    // Current file belongs to this FileSet
                    fs.addFile(f)

                    // mark file being processed
                    processed[f] = true
                }
            }

        }

    }

    // Remove processed files
    files := make([]*File, 0, 0)
	for _, f := range x.Files {
		if !processed[f] {
			files = append(files, f)
		}
	}
    x.Files = files

    // End Result should be that all files within the NZB
    // have been processed
    // xNZB (XML NZB) Files should be completely empty
    if len(x.Files) > 0 {
        return nil, fmt.Errorf("not all file(s) where processed into fileset(s)")
    }

    // Delete Xml NZB
    x = nil

    return n, nil
}

// EOF
