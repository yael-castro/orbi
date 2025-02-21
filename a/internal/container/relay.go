//go:build relay

package container

import (
	"context"
	"database/sql"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/yael-castro/orbi/a/internal/app/business"
	"github.com/yael-castro/orbi/a/internal/app/input/command"
	"github.com/yael-castro/orbi/a/internal/app/output/composite"
	userskafka "github.com/yael-castro/orbi/a/internal/app/output/kafka"
	"github.com/yael-castro/orbi/a/internal/app/output/postgres"
	"github.com/yael-castro/orbi/a/internal/app/output/rpc"
	"github.com/yael-castro/orbi/a/pkg/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

func New() Container {
	logger := log.New(os.Stdout, "[CONTAINER] ", log.LstdFlags)

	return &usersRelay{
		logger: logger,
	}
}

type usersRelay struct {
	logger     *log.Logger
	producer   *kafka.Producer
	clientConn *grpc.ClientConn
	container
}

func (r *usersRelay) Inject(ctx context.Context, a any) (err error) {
	switch a := a.(type) {
	case *func(context.Context, ...string) int:
		return r.injectCommand(ctx, a)
	case **kafka.Producer:
		return r.injectProducer(ctx, a)
	case **grpc.ClientConn:

		// "github.com/yael-castro/orbi/b/pkg/pb"
		return r.injectClientConn(ctx, a)
	}

	return r.container.Inject(ctx, a)
}

func (r *usersRelay) injectCommand(ctx context.Context, cmd *func(context.Context, ...string) int) (err error) {
	// Getting environment variables
	createUserTopic, err := env.Get("CREATE_USER_TOPIC")
	if err != nil {
		return err
	}

	updateUserTopic, err := env.Get("UPDATE_USER_TOPIC")
	if err != nil {
		return err
	}

	// External dependencies
	errLogger := log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile)
	infoLogger := log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile)

	var db *sql.DB
	if err = r.Inject(ctx, &db); err != nil {
		return
	}

	var producer *kafka.Producer
	if err = r.Inject(ctx, &producer); err != nil {
		return
	}

	var conn *grpc.ClientConn
	if err = r.Inject(ctx, &conn); err != nil {
		return
	}

	// Secondary adapters
	reader := postgres.NewMessagesReader(db, infoLogger)

	grpcSender, err := rpc.NewMessageSender(conn, infoLogger)
	if err != nil {
		return err
	}

	kafkaSender := userskafka.NewMessageSender(userskafka.MessageSenderConfig{
		Info:     infoLogger,
		Error:    errLogger,
		Producer: producer,
	})

	sender := composite.NewMessageSender(
		composite.WithTopic(createUserTopic, grpcSender),
		composite.WithTopic(updateUserTopic, kafkaSender),
	)

	confirmer := postgres.NewMessageDeliveryConfirmer(db)

	// Business logic
	messagesRelay, err := business.NewMessagesRelay(business.MessagesRelayConfig{
		Reader:      reader,
		Sender:      sender,
		Confirmer:   confirmer,
		InfoLogger:  infoLogger,
		ErrorLogger: errLogger,
	})
	if err != nil {
		return
	}

	// Primary adapters
	cmdRelay, err := command.Relay(messagesRelay, errLogger)
	if err != nil {
		return
	}

	*cmd = cmdRelay
	return
}

func (r *usersRelay) injectClientConn(ctx context.Context, conn **grpc.ClientConn) (err error) {
	err = r.initClientConn(ctx)
	if err != nil {
		return
	}

	*conn = r.clientConn
	return
}

func (r *usersRelay) initClientConn(context.Context) (err error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	target, err := env.Get("TARGET")
	if err != nil {
		return err
	}

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}

	r.clientConn = conn
	return
}

func (r *usersRelay) injectProducer(ctx context.Context, producer **kafka.Producer) error {
	if err := r.initProducer(ctx); err != nil {
		return err
	}

	*producer = r.producer
	return nil
}

func (r *usersRelay) initProducer(_ context.Context) (err error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	kafkaServers, err := env.Get("KAFKA_SERVERS")
	if err != nil {
		return err
	}

	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaServers,
		"acks":              "all",
	})
	if err != nil {
		return
	}

	r.producer = kafkaProducer
	return
}

func (r *usersRelay) Close(ctx context.Context) (err error) {
	if r.producer != nil {
		r.producer.Close()
		r.logger.Println("Kafka producer is closed")
	}

	if r.clientConn != nil {
		err := r.clientConn.Close()
		if err != nil {
			r.logger.Println("Error trying to close gRPC client", err)
		} else {
			r.logger.Println("gRPC client is closed")
		}
	}

	err = r.container.Close(ctx)
	if err != nil {
		r.logger.Println("Error trying to close container", err)
		return
	}

	r.logger.Println("Container is closed")
	return
}
