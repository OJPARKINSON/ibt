package headers

import "fmt"

// Header contains all sub-headers present in the ibt file.
type Header struct {
	TelemetryHeader *TelemetryHeader
	DiskHeader      *DiskHeader
	VarHeader       map[string]VarHeader
	SessionInfo     *Session
	VarBuffers      []VarBuffer
}

// ParseHeader parses each of the required sub-headers of the ibt file in sequence.
func ParseHeaders(r Reader) (*Header, error) {
	telemetryHeader, err := ReadTelemetryHeader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse telemetry header: %w", err)
	}

	diskHeader, err := ReadDiskHeader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse disk header: %w", err)
	}

	varHeader, err := ReadVarHeader(r, telemetryHeader.NumVars, telemetryHeader.VarHeaderOffset)
	if err != nil {
		return nil, fmt.Errorf("failed to parse variable header: %w", err)
	}

	varBuffers, err := ReadVarBufferHeaders(r, telemetryHeader.NumBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to parse var buffer header: %w", err)
	}

	sessionInfo, err := ReadSessionInfo(r, telemetryHeader.SessionInfoOffset, telemetryHeader.SessionInfoLength)
	if err != nil {
		return nil, fmt.Errorf("failed to parse session info: %w", err)
	}

	return &Header{
		TelemetryHeader: telemetryHeader,
		DiskHeader:      diskHeader,
		VarHeader:       varHeader,
		SessionInfo:     sessionInfo,
		VarBuffers:      varBuffers,
	}, nil
}

func (h *Header) UpdateVarBuffer(r Reader) error {
	varBuffers, err := ReadVarBufferHeaders(r, h.TelemetryHeader.NumBuf)
	if err != nil {
		return fmt.Errorf("failed to parse var buffer header: %w", err)
	}

	h.VarBuffers = varBuffers

	return nil
}
