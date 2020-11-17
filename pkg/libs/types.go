package libs

import topupService "grpc-mirror/pkg/protobufs/destinationpb"

// ChannelBeacon this is the message that traverse through recipients and publisher
type ChannelBeacon struct {
	Response chan interface{}
	Message  *topupService.TopupRequestType
	done     bool
}

// Wait let observers know if they have to wait in the message queue or not
func (c *ChannelBeacon) Wait() bool {
	return !c.done
}

// Done tells the dispatcher to close the channel
func (c *ChannelBeacon) Done() {
	c.done = true
	close(c.Response)
}
