package utils

import (
	"context"
	"firebase.google.com/go"
	"fmt"
	"github.com/labstack/gommon/log"
	"google.golang.org/api/option"
	"os"
)

const (
	CHATALERT   = "ChatAlert"
	EVENTALERT  = "EventAlert"
	FRIENDALERT = "FriendAlert"
	MATCHALERT  = "MatchAlert"
)

type AppAdapterInterface interface {
	SetSubscribe(ctx context.Context, isSubscribe bool, token []string, topic string) (int, error)
	subscribe(ctx context.Context, token []string, topic string) (int, error)
	unSubscribe(ctx context.Context, token []string, topic string) (int, error)
}

type FireBaseAppAdapter struct {
	app *firebase.App
}

func NewFireBaseApp() (AppAdapterInterface, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	var opt option.ClientOption
	if os.Getenv("RELEASE_SYSTEM") == "kubernetes" {
		opt = option.WithCredentialsFile("/config/galmal-8f900-firebase-adminsdk-zhjsl-f6d034ad3b.json")
	} else {
		opt = option.WithCredentialsFile(dir + "/kubernetes/galmal-8f900-firebase-adminsdk-zhjsl-f6d034ad3b.json")
	}
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}
	return &FireBaseAppAdapter{
		app: app,
	}, nil
}

func (f *FireBaseAppAdapter) SetSubscribe(ctx context.Context, isSubscribe bool, token []string, topic string) (int, error) {
	if isSubscribe {
		return f.subscribe(ctx, token, topic)
	} else {
		return f.unSubscribe(ctx, token, topic)
	}
}

func (f *FireBaseAppAdapter) subscribe(ctx context.Context, token []string, topic string) (int, error) {
	client, err := f.app.Messaging(ctx)
	if err != nil {
		return 0, err
	}

	res, err := client.SubscribeToTopic(ctx, token, topic)
	if err != nil {
		return res.SuccessCount, err
	}

	if len(res.Errors) != 0 {
		for _, err := range res.Errors {
			log.Error(err.Reason)
		}
	}

	return res.SuccessCount, nil
}

func (f *FireBaseAppAdapter) unSubscribe(ctx context.Context, token []string, topic string) (int, error) {
	client, err := f.app.Messaging(ctx)
	if err != nil {
		return 0, err
	}
	res, err := client.UnsubscribeFromTopic(ctx, token, topic)
	if err != nil {
		return res.SuccessCount, err
	}
	if len(res.Errors) != 0 {
		for _, err := range res.Errors {
			log.Error(err.Reason)
		}
	}
	return res.SuccessCount, nil
}

//
//type WorkerLauncher interface {
//	LaunchWorker(in chan Work)
//}
//
//type Worker struct {
//	*messaging.Client
//}
//
//func (w *Worker) LaunchWorker(in chan Work) {
//	badge := 10
//	for work := range in {
//
//		message := &messaging.Message{
//			Android: nil,
//			APNS: &messaging.APNSConfig{
//				Headers: map[string]string{
//					"apns-priority": "10",
//				},
//				Payload: &messaging.APNSPayload{
//					Aps: &messaging.Aps{
//						AlertString: "test",
//						Alert:       nil,
//						Badge:       &badge,
//					},
//					CustomData: nil,
//				},
//			},
//			Token: work.Token,
//		}
//
//		result, err := w.Send(context.Background(), message)
//		if err != nil {
//			log.Error(err)
//		}
//		fmt.Println("success", result)
//	}
//}
