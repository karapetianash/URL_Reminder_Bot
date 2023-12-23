package eventConsumer

import (
	"log"
	"sync"
	"time"

	"URLReminderBot/events"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, bs int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: bs,
	}
}

func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERROR] consumer: %s\n", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		if err = c.handleEvents(gotEvents); err != nil {
			log.Printf("[ERROR] event handler: %s\n", err.Error())

			continue
		}
	}
}

// TODO: retry for lost connection
func (c *Consumer) handleEvents(curEvents []events.Event) error {
	wg := sync.WaitGroup{}
	wg.Add(len(curEvents))

	for _, event := range curEvents {
		go func(event events.Event) {
			log.Printf("[INFO] got new event %s", event.Text)

			if err := c.processor.Process(event); err != nil {
				log.Printf("[ERROR] can't handle event: %s", err.Error())
			}

			wg.Done()
		}(event)
	}

	wg.Wait()

	return nil
}
