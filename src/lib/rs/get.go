package rs

import (
	"fmt"
	"io"
	"lib/objectstream"
)

type RSGetStream struct {
	*decoder
}

func NewRSGetStream(locateInfo map[int]string, dataServers []string, hash string, size int64) (*RSGetStream, error) {
	if len(locateInfo)+len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}

	readers := make([]io.Reader, ALL_SHARDS)
	for i := 0; i < ALL_SHARDS; i++ {
		server := locateInfo[i]
		if server == "" {
			locateInfo[i] = dataServers[0]
			dataServers = dataServers[1:]
			continue
		}
		reader, err := objectstream.NewGetStream(server, fmt.Sprintf("%s.%d", hash, i))
		if err == nil {
			readers[i] = reader
		}
	}

	writers := make([]io.Writer, ALL_SHARDS)
	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	var err error
	for i := range readers {
		if readers[i] == nil {
			writers[i], err = objectstream.NewTempPutStream(locateInfo[i], fmt.Sprintf("%s.%d", hash, i), perShard)
			if err != nil {
				return nil, err
			}
		}
	}

	dec := NewDecoder(readers, writers, size)
	return &RSGetStream{dec}, nil
}

func (s *RSGetStream) Close() {
	for _, w := range s.writers {
		if w != nil {
			w.(*objectstream.TempPutStream).Commit(true)
		}
	}
}
