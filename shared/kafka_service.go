package shared

import (
	"cqrs-playground/config"
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

type KafkaService struct {
	brokerList    []string
	SyncProducers []sarama.SyncProducer
	Consumers     []sarama.Consumer
}

func NewKafkaService() *KafkaService {
	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	return &KafkaService{
		brokerList: []string{conf.Kafka.HostAndPortString()},
	}
}

func (k *KafkaService) NewSyncProducer() sarama.SyncProducer {
	c := sarama.NewConfig()
	c.Producer.Return.Successes = true
	c.Producer.RequiredAcks = sarama.WaitForAll
	c.Producer.Retry.Max = 3

	producer, err := sarama.NewSyncProducer(k.brokerList, c)
	if err != nil {
		log.Fatalf("error while creating kafka producer:: %v", err)
	}
	k.SyncProducers = append(k.SyncProducers, producer)
	return producer
}

func (k *KafkaService) SendEvent(producer sarama.SyncProducer, topic string, message []byte) error {
	if producer == nil {
		return fmt.Errorf("error: event could not be sent due to producer being null")
	}
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder("Key"),
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := (producer).SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("message sent successfully! partition=%d, offset=%d\n", partition, offset)
	return nil
}

func (k *KafkaService) NewConsumerOffsetNewest(topic string) sarama.PartitionConsumer {
	return k.createConsumer(topic, 0, sarama.OffsetNewest)
}

func (k *KafkaService) NewConsumerOffsetOldest(topic string) sarama.PartitionConsumer {
	return k.createConsumer(topic, 0, sarama.OffsetOldest)
}

func (k *KafkaService) createConsumer(topic string, partition int32, offset int64) sarama.PartitionConsumer {
	c := sarama.NewConfig()
	c.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumer(k.brokerList, c)
	if err != nil {
		log.Fatalf("error while creating kafka consumer: %v", err)
	}

	p, err := consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		log.Fatalf("error while creating kafka consumer: %v", err)
	}

	k.Consumers = append(k.Consumers, consumer)
	return p
}

func (k *KafkaService) Close() {
	for _, producer := range k.SyncProducers {
		if producer != nil {
			producer.Close()
		}
	}
	for _, consumer := range k.Consumers {
		if consumer != nil {
			consumer.Close()
		}
	}
}
