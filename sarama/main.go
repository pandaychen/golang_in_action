package main

// SIGUSR1 toggle the pause/resume consumption
import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

func main() {

	endpoint := "cn-beijing.log.aliyuncs.com"
	port := "10012"
	version := "2.1.0"
	project := "test-project"
	topics := "your sls logstore"

	accessId := os.Getenv("SLS_ACCESS_KEY_ID")
	accessKey := os.Getenv("SLS_ACCESS_KEY_SECRET")
	group := "test-groupId"

	keepRunning := true
	log.Println("Starting a new Sarama consumer")

	version, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	/**
	 	 * 构建一个新的Sarama配置。
		 * 在初始化消费者/生产者之前，必须定义Kafka集群版本。
	*/
	brokers := []string{fmt.Sprintf("%s.%s:%s", project, endpoint, port)}

	conf := sarama.NewConfig()
	conf.Version = version

	conf.Net.TLS.Enable = true
	conf.Net.SASL.Enable = true
	conf.Net.SASL.User = project
	conf.Net.SASL.Password = fmt.Sprintf("%s#%s", accessId, accessKey)
	conf.Net.SASL.Mechanism = "PLAIN"

	conf.Consumer.Fetch.Min = 1
	conf.Consumer.Fetch.Default = 1024 * 1024
	conf.Consumer.Retry.Backoff = 2 * time.Second
	conf.Consumer.MaxWaitTime = 250 * time.Millisecond
	conf.Consumer.MaxProcessingTime = 100 * time.Millisecond
	conf.Consumer.Return.Errors = false
	conf.Consumer.Offsets.AutoCommit.Enable = true
	conf.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	conf.Consumer.Offsets.Initial = sarama.OffsetNewest
	conf.Consumer.Offsets.Retry.Max = 3

	/**
	 * 设置一个新的Sarama消费者组
	 */
	consumer := Consumer{
		ready: make(chan bool),
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(brokers, group, conf)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume`应该在一个无限循环内调用，当服务器端重新平衡时，消费者会话将需要重新创建以获取新的声明
			if err := client.Consume(ctx, strings.Split(topics, ","), &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// 检查上下文是否被取消，表示消费者应该停止
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // 等待消费者设置完成
	log.Println("Sarama consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		log.Println("Resuming consumption")
	} else {
		client.PauseAll()
		log.Println("Pausing consumption")
	}

	*isPaused = !*isPaused
}

// Consumer表示Sarama消费者组消费者
type Consumer struct {
	ready chan bool
}

// Setup在新会话开始时运行，在ConsumeClaim之前
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// 将消费者标记为已准备好
	close(consumer.ready)
	return nil
}

// Cleanup在会话结束时运行，一旦所有ConsumeClaim goroutine退出
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim必须启动ConsumerGroupClaim的Messages()的消费者循环
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// 注意：
	// 不要把下面的代码移到goroutine里。
	// ConsumeClaim本身在goroutine中调用，参见：https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message := <-claim.Messages():
			realUnixTimeSeconds := message.Timestamp.Unix()
			if realUnixTimeSeconds < 2000000 {
				realUnixTimeSeconds = message.Timestamp.UnixMicro() / 1000
			}

			log.Printf("Message claimed: value = %s, timestamp = %d, topic = %s", string(message.Value), realUnixTimeSeconds, message.Topic)
			session.MarkMessage(message, "")

		// 当session.Context()完成时应返回。
		// 如果不这样做，当kafka重新平衡时，会引发ErrRebalanceInProgress或read tcp <ip>:<port>: i/o timeout错误。参见：https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}
