package internal

import (
	"encoding/json"
	"fmt"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

const (
	targetLabelKey      = "app.kubernetes.io/name"
	targetLabelValue    = "coder-workspace"
	initContainerName   = "clone-repo"
	initContainerImage  = "alpine/git:2.45.2"
	initContainerScript = `set -eux
cd /home/coder
rm -rf ray-example
git clone --depth=1 https://github.com/leo6103/ray-example ray-example
wget -O /home/coder/kubox-ray-sync https://github.com/leo6103/kubox-ray-sync/releases/download/v0.0.1/kubox-ray-sync
chmod +x /home/coder/kubox-ray-sync
chown -R 1000:1000 /home/coder/ray-example || true
chown 1000:1000 /home/coder/kubox-ray-sync || true
`
	sharedVolumeName    = "home"
	sharedMountPath     = "/home/coder"
)

// MutateService logs AdmissionReview payloads for future mutation logic.
type MutateService struct {
	logPath string
}

// NewMutateService constructs a MutateService writing to the provided log path.
func NewMutateService(logPath string) *MutateService {
	return &MutateService{logPath: logPath}
}

func (s *MutateService) logLine(msg string) {
	if msg == "" {
		return
	}
	_ = AppendBody(s.logPath, []byte(msg+"\n"))
}

// Mutate processes an AdmissionReview and returns the AdmissionReview response payload.
func (s *MutateService) Mutate(body []byte) ([]byte, error) {
	if err := AppendBody(s.logPath, body); err != nil {
		s.logLine("error log request")
		return nil, fmt.Errorf("log request: %w", err)
	}

	var review admissionv1.AdmissionReview
	if err := json.Unmarshal(body, &review); err != nil {
		s.logLine("error decode admission review")
		return nil, fmt.Errorf("decode admission review: %w", err)
	}

	response := admissionv1.AdmissionReview{
		TypeMeta: review.TypeMeta,
	}
	if response.Kind == "" {
		response.Kind = "AdmissionReview"
	}
	if response.APIVersion == "" {
		response.APIVersion = admissionv1.SchemeGroupVersion.String()
	}

	resp := &admissionv1.AdmissionResponse{Allowed: true}
	response.Response = resp

	if review.Request == nil {
		encoded, err := json.Marshal(response)
		if err != nil {
			s.logLine("error encode admission response")
			return nil, fmt.Errorf("encode admission response: %w", err)
		}
		return encoded, nil
	}

	resp.UID = review.Request.UID

	patch, err := s.buildInitPatch(review.Request)
	if err != nil {
		s.logLine("error build init patch")
		return nil, err
	}

	if len(patch) > 0 {
		resp.Patch = patch
		patchType := admissionv1.PatchTypeJSONPatch
		resp.PatchType = &patchType
	}

	encoded, err := json.Marshal(response)
	if err != nil {
		s.logLine("error encode admission response")
		return nil, fmt.Errorf("encode admission response: %w", err)
	}

	return encoded, nil
}

func (s *MutateService) buildInitPatch(req *admissionv1.AdmissionRequest) ([]byte, error) {
	if req == nil {
		s.logLine("buildInitPatch: empty request")
		return nil, nil
	}

	if req.Operation != admissionv1.Create {
		s.logLine("buildInitPatch: non-create operation")
		return nil, nil
	}

	if req.Kind.Kind != "Pod" || req.Kind.Group != "" {
		s.logLine("buildInitPatch: unsupported kind")
		return nil, nil
	}

	var pod corev1.Pod
	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		s.logLine("error decode pod")
		return nil, fmt.Errorf("decode pod: %w", err)
	}

	if pod.Labels[targetLabelKey] != targetLabelValue {
		s.logLine("buildInitPatch: label mismatch")
		return nil, nil
	}

	if hasInitContainer(pod.Spec.InitContainers) {
		s.logLine("buildInitPatch: init already present")
		return nil, nil
	}

	initContainer := corev1.Container{
		Name:    initContainerName,
		Image:   initContainerImage,
		Command: []string{"sh", "-c", initContainerScript},
		VolumeMounts: []corev1.VolumeMount{{
			Name:      sharedVolumeName,
			MountPath: sharedMountPath,
		}},
	}

	patchOps := make([]patchOperation, 0, 1)

	if len(pod.Spec.InitContainers) == 0 {
		patchOps = append(patchOps, patchOperation{
			Op:   "add",
			Path: "/spec/initContainers",
			Value: []corev1.Container{
				initContainer,
			},
		})
	} else {
		patchOps = append(patchOps, patchOperation{
			Op:    "add",
			Path:  "/spec/initContainers/-",
			Value: initContainer,
		})
	}

	patch, err := json.Marshal(patchOps)
	if err != nil {
		s.logLine("error encode patch")
		return nil, fmt.Errorf("encode patch: %w", err)
	}

	s.logLine("buildInitPatch: patch created")
	return patch, nil
}

func hasInitContainer(containers []corev1.Container) bool {
	for _, c := range containers {
		if c.Name == initContainerName {
			return true
		}
	}
	return false
}

// patchOperation is a helper for JSON Patch serialization.
type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}
