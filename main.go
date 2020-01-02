package close

import "io"

var closers []io.Closer

func Later(c io.Closer) {
	closers = append(closers, c)
}

func All() {
	errs := make([]error, 0, len(closers))
	for _, c := range closers {
		if err := c.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	for _, err := range errs {
		panic(err)
	}
}
