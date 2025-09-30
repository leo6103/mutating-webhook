package internal

import "os"

// AppendBody writes the provided bytes to the log file. If the payload is
// missing a trailing newline we append ",\n" to keep JSON entries comma
// delimited. Callers that want a plain log line should include a newline.
func AppendBody(path string, body []byte) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	if len(body) > 0 {
		if _, err := f.Write(body); err != nil {
			return err
		}
	}
	if len(body) == 0 || body[len(body)-1] != '\n' {
		if _, err := f.Write([]byte(",\n")); err != nil {
			return err
		}
	}

	return nil
}
