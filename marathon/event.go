package marathon

import (
	"fmt"
	"github.com/ContainX/depcon/pkg/encoding"
	"github.com/donovanhide/eventsource"
	"net/http"
)

func (c *MarathonClient) CreateEventStreamListener(channel EventsChannel, filter int) error {
	c.Lock()
	defer c.Unlock()

	// no-op if already listening on a stream
	if c.eventStreamState != nil {
		return nil
	}

	if err := c.setupSSEStream(); err != nil {
		return err
	}

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

func (c *MarathonClient) setupSSEStream() error {
	request, err := c.http.CreateHttpRequest(http.MethodGet, c.marathonUrl(API_EVENTS), nil)
	if err != nil {
		return err
	}

	stream, err := eventsource.SubscribeWith("", http.DefaultClient, request)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case ev := <-stream.Events:
				if err := c.handleStreamEvent(ev.Data()); err != nil {
					log.Errorf("Error setting up SSE Stream: %v", err)
				}
			case err := <-stream.Errors:
				log.Errorf("Error setting up SSE Stream: %v", err)
			}
		}
	}()
	return nil
}

func (c *MarathonClient) handleStreamEvent(data string) error {
	if data == "" {
		return nil
	}

	eventType := new(EventType)

	if err := encoding.DefaultJSONEncoder().UnMarshalStr(data, eventType); err != nil {
		return fmt.Errorf("Failed to decode event, content: %s, error: %s", data, err)
	}

	event, err := c.GetEvent(eventType.EventType)
	if err != nil {
		return fmt.Errorf("Unable to handle event type, type: %s, error: %s", eventType.EventType, err)
	}

	if err := encoding.DefaultJSONEncoder().UnMarshalStr(data, event.Event); err != nil {
		return fmt.Errorf("Failed to decode event, id: %d, error: %s", event.ID, err)
	}

	if event.ID&c.eventStreamState.filter != 0 {
		go func(ch EventsChannel, e *Event) {
			ch <- e
		}(c.eventStreamState.channel, event)
	}
	return nil
}
