package logger

// These two structs are dummy structs that implement Write() to satisfy
// io.Writer.
type logOut struct {
	log *Log
}

type logErr struct {
	log *Log
}

func newDummyLogs(l *Log) (*logOut, *logErr) {
	return &logOut{l}, &logErr{l}
}

func (l *logOut) Write(p []byte) (int, error) {
	return l.log.stdout(p)
}

func (l *logErr) Write(p []byte) (int, error) {
	return l.log.stderr(p)
}
