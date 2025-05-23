package subscriber

import "github.com/segmentio/kafka-go"

type GeneralEventHandler struct {
	GeneralEventReader *kafka.Reader
	MembershipEventHandler
	SampleOrderGeneralEventHandler
}
