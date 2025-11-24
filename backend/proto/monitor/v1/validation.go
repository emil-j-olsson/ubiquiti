package monitorv1

import "errors"

func (r *RegisterDeviceRequest) Validate() error {
	if r == nil {
		return errors.New("empty request")
	}
	if len(r.GetHost()) == 0 {
		return errors.New("missing host in request")
	}
	if r.GetPort() <= 0 || r.GetPort() > 65535 {
		return errors.New("invalid port in request")
	}
	if r.GetPortGateway() <= 0 || r.GetPortGateway() > 65535 {
		return errors.New("invalid gateway port in request")
	}
	if r.GetProtocol() == Protocol_PROTOCOL_UNSPECIFIED {
		return errors.New("missing protocol in request")
	}
	return nil
}

func (r *DiagnosticsRequest) Validate() error {
	if r == nil {
		return errors.New("empty request")
	}
	if len(r.GetDeviceId()) == 0 {
		return errors.New("missing device_id in request")
	}
	return nil
}
