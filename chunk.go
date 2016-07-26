package nzb

import (
    "fmt"
    "sync"
)

// Chunk of data to download
type Chunk struct {
    Groups []string
    Segment *Segment
}

// Chunks list of all chunks within NZB
type Chunks struct {
    c       []*Chunk
    mu      *sync.Mutex
    Total   int
    Marker  int
}

func (c *Chunks) addChunks(groups []string, segments []*Segment) {
    for _, s := range segments {
        cnk := &Chunk {
            Groups: groups,
            Segment: s,
        }

        c.c = append(c.c, cnk)
    }
}

// GetChunks limited by max
func (c *Chunks) GetChunks(max int) []*Chunk {
    c.mu.Lock()
    defer c.mu.Unlock()

    chunks := make([]*Chunk, max)

    if c.Marker < len(c.c) {
        for i := 0; i < max; i++ {
            chunks[i] = c.c[c.Marker]
            c.Marker++
        }
    }

    return chunks
}

// GetNext get next chunk from Chunks
func (c *Chunks) GetNext() *Chunk {
    c.mu.Lock()
    defer c.mu.Unlock()

    if c.Marker < len(c.c) {
        cnk := c.c[c.Marker]
        c.Marker++
        return cnk
    }

    return nil
}

// Remove chunk from chunks based on article id
// This allows to re-continue a download
// and remove already downloaded chunks from the chunk list
func (c *Chunks) Remove(id string) error {

    c.mu.Lock()
    defer c.mu.Unlock()

    rm := -1
    for i, cnk := range c.c {
        if cnk.Segment.ID == id {
            // Found chunk to remove
            rm = i
        }
    }

    // Remove element
    if rm > -1 {
        c.c = append(c.c[:rm], c.c[rm+1:]...)
        // Update Total
        c.Total = len(c.c)

        return nil
    }

    return fmt.Errorf("id does not exists")
}

// EOF
