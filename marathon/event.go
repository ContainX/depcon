package marathon

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ContainX/depcon/pkg/encoding"
	"github.com/donovanhide/eventsource"
)

func (c *MarathonClient) CreateEventStreamListener(channel EventsChannel, filter int) error {
	c.Lock()
	defer c.Unlock()

	// no-op if already listening on a stream
	if c.eventStreamState != nil {
		return nil
	}

	go func() {
		for {
			stream, err := c.setupSSEStream()
			if err != nil {
				log.Debugf("error connecting to SSE subscription: %s", err.Error())
				<-time.After(5 * time.Second)
				continue
			}
			err = c.listenToSSE(stream)
			if err != nil {
				log.Errorf("error on SSE subscription: %s", err)
			}
			stream.Close()
		}
	}()

	c.eventStreamState = &EventStreamState{
		channel: channel,
		filter:  filter,
	}

	return nil
}

func (c *MarathonClient) CloseEventStreamListener(channel EventsChannel) {
	c.Lock()
	defer c.Unlock()

	c.eventStreamState = nil
}

func (c *MarathonClient) setupSSEStream() (*eventsource.Stream, error) {
	request, err := c.http.CreateHttpRequest(http.MethodGet, c.marathonUrl(API_EVENTS), nil)
	if err != nil {
		return nil, err
	}

	stream, err := eventsource.SubscribeWith("", http.DefaultClient, request)
	if err != nil {
		return nil, err
	}

	return stream, nil
}

func (c *MarathonClient) listenToSSE(stream *eventsource.Stream) error {
	for {
		select {
		case ev := <-stream.Events:
			if err := c.handleStreamEvent(ev.Data()); err != nil {
				log.Errorf("error setting up SSE Stream: %v", err)
			}
		case err := <-stream.Errors:
			log.Errorf("error setting up SSE Stream: %v", err)
		}
	}
}

func (c *MarathonClient) handleStreamEvent(data string) error {
	if data == "" {
		return nil
	}

	eventType := new(EventType)

	if err := encoding.DefaultJSONEncoder().UnMarshalStr(data, eventType); err != nil {
		return fmt.Errorf("failed to decode event, content: %s, error: %s", data, err)
	}

	event, err := c.GetEvent(eventType.EventType)
	if err != nil {
		return fmt.Errorf("unable to handle event type, type: %s, error: %s", eventType.EventType, err)
	}

	if err := encoding.DefaultJSONEncoder().UnMarshalStr(data, event.Event); err != nil {
		return fmt.Errorf("failed to decode event, id: %d, error: %s", event.ID, err)
	}

	if event.ID&c.eventStreamState.filter != 0 {
		go func(ch EventsChannel, e *Event) {
			ch <- e
		}(c.eventStreamState.channel, event)
	}
	return nil
}
