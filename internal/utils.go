package internal

import "os"

// AppendBody writes the request body followed by ",\n" to the log file.
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
	if _, err := f.Write([]byte(",\n")); err != nil {
		return err
	}

	return nil
}
