package ranger

import (
	"errors"
	"io"
)

var defaultBuffSize = 1024 * 512
var defaultFetchSize = 1024 * 128
var defaultRoutinesize = 2

type RingBuffReader struct {
	// the range fetcher with which to download blocks
	fetcher RangeFetcher

	buff       []byte
	readPoint  int64
	writePoint int64

	Length int64

	//suggested parameter (buffSzie = N * fetchSize * routineSize, N is integer)
	buffSize    int
	fetchSize   int
	routineSize int
}

// ReadAt reads len(p) bytes from the ranged-over source.
// It returns the number of bytes read and the error, if any.
// ReadAt always returns a non-nil error when n < len(b). At end of file, that error is io.EOF.
func (r *RingBuffReader) ReadAt(p []byte, off int64) (int, error) {
	if off < 0 {
		return 0, errors.New("read before beginning of file")
	}
	if off >= r.Length {
		return 0, errors.New("read beyond end of file")
	}
	r.readPoint = off
	return r.Read(p)
}

// Read reads len(p) bytes from ranged-over source.
// It returns the number of bytes read and the error, if any.
// EOF is signaled by a zero count with err set to io.EOF.
func (r *RingBuffReader) Read(p []byte) (int, error) {
	distance := r.writePoint - r.readPoint
	// all the zone is [-length,0](0,buffSize](buffSize,length]
	// (0,buffSize]
	if distance > 0 && distance <= int64(r.buffSize) {
		readIndex := r.readPoint % int64(r.buffSize)
		writeIndex := r.readPoint % int64(r.buffSize)
		if len(p) < int(distance) {
			length := ringRead(p, r.buff, int(readIndex), int(writeIndex))
			r.readPoint = r.readPoint + int64(length)
			if r.readPoint >= r.Length {
				return length, io.EOF
			} else {
				return length, nil
			}
		} else {
			length := ringRead(p, r.buff, int(readIndex), int(writeIndex))
			r.readPoint = r.readPoint + int64(length)
			if r.readPoint >= r.Length {
				return length, io.EOF
			}
			err := r.fillBuff()
			if err != nil {
				return length, err
			}
			return r.Read(p[length:])
		}
	} else {
		// (buffSize,length] U [-length,0)
		r.writePoint = r.readPoint
		r.fillBuff()
		return r.Read(p)
	}
}

// Seek sets the offset for the next Read to offset, interpreted
// according to whence: 0 means relative to the origin of the file, 1 means relative
// to the current offset, and 2 means relative to the end. It returns the new offset
// and an error, if any.
func (r *RingBuffReader) Seek(off int64, whence int) (int64, error) {

	switch whence {
	case 0: // set
		r.readPoint = off
	case 1: // cur
		off = r.readPoint + off
	case 2: // end
		off = r.readPoint + off
	}

	if off > r.Length {
		return 0, errors.New("seek beyond end of file")
	}

	if off < 0 {
		return 0, errors.New("seek before beginning of file")
	}

	r.readPoint = off
	return r.readPoint, nil
}

func ringRead(p []byte, ringBuff []byte, readIndex int, writeIndex int) int {
	if writeIndex > readIndex {
		return copy(p, ringBuff[readIndex:writeIndex])
	} else {
		length := copy(p, ringBuff[readIndex:])
		length = length + copy(p[length:], ringBuff[:writeIndex])
		return length
	}
}

func ringWrite(p []byte, ringBuff []byte, readIndex int, writeIndex int) int {
	if writeIndex >= readIndex {
		length := copy(ringBuff[writeIndex:], p)
		length = length + copy(ringBuff[:readIndex], p[length:])
		return length
	} else {
		return copy(ringBuff[writeIndex:readIndex], p)
	}
}

func (r *RingBuffReader) fillBuff() error {
	if r.writePoint >= r.Length {
		//reach the end
		return nil
	}
	distance := r.writePoint - r.readPoint
	// distance [0,buffSize)
	if distance >= 0 && distance < int64(r.buffSize) {
		fillSize := r.buffSize - int(distance)
		httpRangeStart := r.writePoint
		httpRangeEnd := r.writePoint + int64(fillSize) - 1
		ranges := make([]ByteRange, 0, 5)
		rangeIndex := httpRangeStart
		i := 0
		for {
			if httpRangeEnd-rangeIndex+1 > int64(r.fetchSize) {
				byteRange := ByteRange{
					Start: rangeIndex,
					End:   rangeIndex + int64(r.fetchSize) - 1,
				}
				if byteRange.End == httpRangeEnd {
					break
				}
				rangeIndex = byteRange.End + 1
				ranges = append(ranges, byteRange)
			} else {
				byteRange := ByteRange{
					Start: rangeIndex,
					End:   httpRangeEnd,
				}
				ranges = append(ranges, byteRange)
				break
			}
			i++
		}
		
		blocks, err := r.fetcher.FetchRanges(ranges)
		if err != nil {
			return err
		}
		for _, value := range blocks {
			writeIndex := r.writePoint % int64(r.buffSize)
			readIndex := r.readPoint % int64(r.buffSize)
			length := ringWrite(value.Data[:value.Length], r.buff, int(readIndex), int(writeIndex))
			r.writePoint = r.writePoint + int64(length)
		}
		return nil
	}
	// distance [buffSize,buffSize]
	if distance == int64(r.buffSize) {
		//no need to fill
		return nil
	}
	// distance (-∞,0) U (buffSize,+∞)
	if distance > int64(r.buffSize) {
		r.writePoint = r.readPoint
		return r.fillBuff()
	}
	return nil
}

func NewRingBuffReader(fetcher RangeFetcher, size ...int) (*RingBuffReader, error) {
	r := &RingBuffReader{
		fetcher:     fetcher,
		buffSize:    defaultBuffSize,
		fetchSize:   defaultBuffSize,
		routineSize: defaultRoutinesize,
	}
	for index, value := range size {
		switch index {
		case 0:
			r.buffSize = value
		case 1:
			r.fetchSize = value
		case 2:
			r.routineSize = value
		}
	}
	if r.buffSize <= 0 {
		return r, errors.New("buff size must be greater than 0")
	}
	if r.fetchSize <= 0 {
		return r, errors.New("fetch size must be greater than 0")
	}
	if r.routineSize <= 0 {
		return r, errors.New("routine size must be greater than 0")
	}
	if r.buffSize < r.fetchSize {
		return r, errors.New("buff size must be greater than or equal to fetch size")
	}
	r.fetcher.ExpectedLength()

	length, err := r.fetcher.ExpectedLength()
	if err != nil {
		return r, err
	}
	r.Length = length
	r.buff = make([]byte, r.buffSize, r.buffSize)
	return r, nil
}