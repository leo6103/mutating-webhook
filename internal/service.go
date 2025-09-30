package internal

// MutateService logs AdmissionReview payloads for future mutation logic.
type MutateService struct {
	logPath string
}

// NewMutateService constructs a MutateService writing to the provided log path.
func NewMutateService(logPath string) *MutateService {
	return &MutateService{logPath: logPath}
}

// Mutate records the raw request body and returns success.
func (s *MutateService) Mutate(body []byte) error {
	if err := AppendBody(s.logPath, body); err != nil {
		return err
	}

	// TODO: Implement AdmissionReview mutation using the Kubernetes Go client.
	return nil
}
