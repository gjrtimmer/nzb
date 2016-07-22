package nzb

import (
    "encoding/json"
    "io/ioutil"
    "strings"
    "sync"
)


/*
 * NOTE:
 * The code to filter the files as fileset(s) was
 * based upon the code of DanielMorsing (https://github.com/DanielMorsing)
 * Code Origin for filtering par2 files: https://github.com/DanielMorsing/gonzbee/blob/master/par2.go
 */

// NZB File
type NZB struct {
    Size        Size        `json:"Size"`
    FileSets    []*FileSet  `json:"FileSets"`
}

type xNZB struct {
    Files       []*File     `xml:"file"`
}

/*
 * Get FileSet by prefix
 * if it does not exists create a new one
 */
func (n *NZB) getFileSetByPrefix(p string) *FileSet {
    for _, fs := range n.FileSets {
        if strings.Compare(fs.Name, p) == 0 {
            return fs
        }
    }

    // Fileset does not exists
    fs := &FileSet {
        Name: p,
        ParSet: &ParSet {},
    }

    n.FileSets = append(n.FileSets, fs)

    return fs
}

// Save NZB File
func (n *NZB) Save(file string) error {
    data, _ := json.MarshalIndent(n, "", "  ")
    return ioutil.WriteFile(file, data, 0644)
}

// Load NZB from JSON
func Load(file string) (*NZB, error) {
    b, err := ioutil.ReadFile(file)
    if err != nil {
        return nil, err
    }

    // Create new NZB
    n := new(NZB)
    err = json.Unmarshal(b, n)
    if err != nil {
        return nil, err
    }

    return n, nil
}

// GenerateChunkList generate new chunk list from NZB
func (n *NZB) GenerateChunkList() *Chunks {

    c := &Chunks {
        c: make([]*Chunk, 0),
        mu: new(sync.Mutex),
        Total: 0,
        Marker: 0,
    }

    for _, fs := range n.FileSets {

        // First Append all chuncks for parent par file
        c.addChunks(fs.ParSet.Parent.Groups, fs.ParSet.Parent.Segments)

        // Add Files
        for _, f := range fs.Files {
            c.addChunks(f.Groups, f.Segments)
        }

        // Add Par Files
        for _, p := range fs.ParSet.Files {
            c.addChunks(p.Groups, p.Segments)
        }
    }

    c.Total = len(c.c)

    return c
}


// EOF
