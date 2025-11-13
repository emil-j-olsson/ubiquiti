package devicev1

import "errors"

func (r *UpdateDeviceRequest) Validate() error {
	if r == nil {
		return errors.New("empty request")
	}
	if r.GetDeviceStatus() == DeviceStatus_DEVICE_STATUS_UNSPECIFIED {
		return errors.New("device status is unspecified or unknown")
	}
	return nil
}
