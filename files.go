package nzb

import (
    "strings"
)

// FileSet .
type FileSet struct {
    Name        string      `json:"Name"`
    ParSet      *ParSet     `json:"ParSet"`
    Files       []*File     `json:"Files"`
    Size        Size        `json:"Size"`
}

// File .
type File struct {
    Filename    string      `json:"Filename"`
    Size        Size        `json:"Size"`

    // XML
    Poster      string      `xml:"poster,attr" json:"Poster"`
    Date        int         `xml:"date,attr" json:"Date"`
    Subject     Subject     `xml:"subject,attr"`

    Groups      []string    `xml:"groups>group" json:"Groups"`
    Segments    []*Segment  `xml:"segments>segment" json:"Subject"`
}

// ParSet .
type ParSet struct {
    Parent      *File       `json:"Parent"`
    TotalBlocks int         `json:"TotalBlocks"`
    Files       []*ParFile  `json:"Files"`
    Size        Size        `json:"Size"`
}

// ParFile .
type ParFile struct {
    *File
    Blocks      int         `json:"Blocks"`
}

func (fs *FileSet) addFile(file *File) {

    // Check to make sure file is not already present for unknown reason
    // mmore failsafe then necessary
    for _, f := range fs.Files {
        if strings.Compare(f.Filename, file.Filename) == 0 {
            // File already present in fileset
            return
        }
    }

    // file not present append
    fs.Files = append(fs.Files, file)

    // update fileset size
    fs.Size += Size(file.Size)
}

// Convert2ParFile convert *File -> *ParFile
func (f *File) convert2ParFile(blocks int) *ParFile {

    p := &ParFile {
        File: f,
        Blocks: blocks,
    }

    return p
}

func (p *ParSet) addParFile(file *ParFile) {

    // Check to make sure file is not already present for unknown reason
    // mmore failsafe then necessary
    for _, f := range p.Files {
        if strings.Compare(f.Filename, file.Filename) == 0 {
            // File already present in repair set
            return
        }
    }

    // file not present append
    p.Files = append(p.Files, file)

    // Update parset size
    p.Size += Size(file.Size)

    // Update Total present block
    p.TotalBlocks += file.Blocks
}

// EOF
