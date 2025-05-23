package model

import "time"

type SampleTubeDetails struct {
	TubeType       string
	CollectionTime time.Time
	ReceiveCount   int32
}

type OrderTransmissionResult struct {
	SampleId int    `json:"sample_id,omitempty"`
	TubeType string `json:"tube_type,omitempty"`
	Status   string `json:"status,omitempty"`
}
