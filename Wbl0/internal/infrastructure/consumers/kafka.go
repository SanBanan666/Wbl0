package consumers

import (
	"context"
	"fmt"
	"log"
	"sync"

	"WbServis/Wbl0/internal/application/interfaces"

	"github.com/IBM/sarama"
)

type kafkaConsumer struct {
	consumer sarama.ConsumerGroup
	topics   []string
	handler  interfaces.OrderService
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func NewKafkaConsumer(brokers []string, groupID string, topics []string, handler interfaces.OrderService) (interfaces.MessageConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &kafkaConsumer{
		consumer: consumer,
		topics:   topics,
		handler:  handler,
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

func (k *kafkaConsumer) Start() error {
	log.Printf("Starting Kafka consumer for topics: %v", k.topics)

	k.wg.Add(1)
	go func() {
		defer k.wg.Done()
		for {
			select {
			case <-k.ctx.Done():
				return
			default:
				err := k.consumer.Consume(k.ctx, k.topics, k)
				if err != nil {
					log.Printf("Error from consumer: %v", err)
				}
			}
		}
	}()

	return nil
}

func (k *kafkaConsumer) Stop() error {
	log.Println("Stopping Kafka consumer...")
	k.cancel()
	k.wg.Wait()
	return nil
}

func (k *kafkaConsumer) Close() error {
	return k.consumer.Close()
}

func (k *kafkaConsumer) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Kafka consumer setup completed")
	return nil
}

func (k *kafkaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("Kafka consumer cleanup completed")
	return nil
}

func (k *kafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			log.Printf("Received message from topic %s, partition %d, offset %d",
				message.Topic, message.Partition, message.Offset)

			err := k.handler.ProcessMessage(message.Value)
			if err != nil {
				log.Printf("Failed to process message: %v", err)
			} else {
				session.MarkMessage(message, "")
			}

		case <-session.Context().Done():
			return nil
		}
	}
}
