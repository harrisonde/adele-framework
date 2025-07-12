package logger

func (w *IoToLogWriter) Write(b []byte) (int, error) {
	n := len(b)
	if n > 0 && b[n-1] == '\n' {
		b = b[:n-1]
	}
	if w.Type == "Error" {
		w.Entry.Error(string(b))
	} else {
		w.Entry.Info(string(b))
	}
	return n, nil
}
