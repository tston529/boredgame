package engine

import (
	"bytes"
	"io"
)

type sz struct {
	x int8
	y int8
}

func lineCounter(r io.Reader) (sz, error) {
	buf := make([]byte, 32*1024)
	count := 0
	x := 0
	lineSep := []byte{'\n'}

xy_exit:
	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)
		x += len(bytes.Split(buf[:c], []byte("\n"))[0])

		switch {
		case err == io.EOF:
			break xy_exit
			//return count, nil

		case err != nil:
			break xy_exit
			//return count, err
		}
	}

	return sz{int8(x), int8(count)}, nil
}
