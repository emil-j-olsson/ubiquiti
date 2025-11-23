package monitorv1

import "errors"

func (r *DiagnosticsRequest) Validate() error {
	if r == nil {
		return errors.New("empty request")
	}
	if len(r.GetDeviceId()) == 0 {
		return errors.New("missing device_id in request")
	}
	return nil
}
