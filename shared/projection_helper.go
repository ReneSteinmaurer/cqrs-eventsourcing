package shared

import (
	"context"
	"log"
)

func ListenToEvent(ctx context.Context, kafkaService *KafkaService, eventType string, applyFunc func([]byte)) {
	consumer := kafkaService.NewConsumerOffsetNewest(eventType)
	defer func() {
		log.Printf("Closing consumer for %s...\n", eventType)
		_ = consumer.Close()
	}()

	msgs := consumer.Messages()

	for {
		select {
		case <-ctx.Done():
			log.Printf("%s listener stopped\n", eventType)
			return
		case msg := <-msgs:
			applyFunc(msg.Value)
		}
	}
}
