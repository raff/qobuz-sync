package client

import (
	"io"
	"log"
	"os"
	"runtime"
	"sync"

	//"github.com/rs/zerolog/log"
	"github.com/gobs/httpclient"
)

const (
	maxConns = 4
)

func copyPart(r io.ReaderAt, w io.WriterAt, start, end int64) (err error) {
	var b [65536]byte

	for p := start; p < end; {
		n := int(end - p)
		if n > len(b) {
			n = len(b)
		}
		n, err = r.ReadAt(b[:n], p)
		if err != nil {
			log.Println("error reading at", p, err)
			return err
		}

		if _, err := w.WriteAt(b[:n], p); err != nil {
			log.Println("error writing at", p, err)
			return err
		}

		p += int64(n)
	}

	return nil
}

func DownloadFile(fileURL string, w *os.File) error {
	nConns := runtime.GOMAXPROCS(0)
	if nConns > maxConns {
		nConns = maxConns
	}

	r, err := httpclient.OpenHttpFile(fileURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer r.Close()
	fsize := r.Size()

	if err := w.Truncate(fsize); err != nil {
		return err
	}

	psize := fsize / int64(nConns)

	var wg sync.WaitGroup

	for i := 0; i < nConns; i++ {
		wg.Add(1)

		go func(i int) {
			pstart := psize * int64(i)
			pend := pstart + psize
			if pend > fsize {
				pend = fsize
			}

			err := copyPart(r, w, pstart, pend)
			if err != nil {
				log.Fatal(err)
			}

			wg.Done()
		}(i)
	}

	wg.Wait()
	return nil
}
