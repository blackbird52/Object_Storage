package rs

import (
	"fmt"
	"io"
	"lib/objectstream"
)

type RSPutStream struct {
	*encoder
}

func NewRSPutStream(dataServers []string, hash string, size int64) (*RSPutStream, error) {
	if len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}
	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	writers := make([]io.Writer, ALL_SHARDS)
	var err error
	for i := range writers {
		writers[i], err = objectstream.NewTempPutStream(dataServers[i], fmt.Sprintf("%s.%d", hash, i), perShard)
		if err != nil {
			return nil, err
		}
	}
	enc := NewEncoder(writers)
	return &RSPutStream{enc}, nil
}

func (s *RSPutStream) Commit(success bool) {
	s.Flush()
	for i := range s.writers {
		s.writers[i].(*objectstream.TempPutStream).Commit(success)
	}
}
