package comm

import (
	"testing"
	"time"

	"github.com/mongodb/grip/send"
	"github.com/mongodb/grip/slogger"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHeartbeat(t *testing.T) {

	Convey("With a simple heartbeat ticker", t, func() {
		sigChan := make(chan Signal)
		mockCommunicator := &MockCommunicator{}
		hbTicker := &HeartbeatTicker{
			MaxFailedHeartbeats: 10,
			Interval:            10 * time.Millisecond,
			SignalChan:          sigChan,
			TaskCommunicator:    mockCommunicator,
			Logger: &slogger.Logger{
				Appenders: []send.Sender{slogger.StdOutAppender()},
			},
		}

		Convey("abort signals detected by heartbeat are sent on sigChan", func() {
			mockCommunicator.setShouldFail(false)
			hbTicker.StartHeartbeating()
			go func() {
				time.Sleep(2 * time.Second)
				mockCommunicator.setAbort(true)
			}()
			signal := <-sigChan
			So(signal, ShouldEqual, AbortedByUser)
		})

		Convey("failed heartbeats must signal failure on sigChan", func() {
			mockCommunicator.setAbort(false)
			mockCommunicator.setShouldFail(true)
			hbTicker.StartHeartbeating()
			signal := <-sigChan
			So(signal, ShouldEqual, HeartbeatMaxFailed)
		})

	})

}
