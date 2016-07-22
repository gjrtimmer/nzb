package nzb

// Segment of File in NZB
type Segment struct {
    ID      string          `xml:",innerxml" json:"ID"`
    Number  int             `xml:"number,attr" json:"Number"`
    Bytes   int             `xml:"bytes,attr" json:"Bytes"`
    Exists  bool            `json:"Exists"`
}

// EOF
