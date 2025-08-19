package chat

import "context"

type ServicePort interface {
	Send(ctx context.Context, from string, req SendRequest)(string,error) //this is for sending message to persistnat databse or data stire
	Subscribe(ctx context.Context, userID string, write func([]byte)) (func(),error)
	Replay(ctx context.Context,userID string, write func([]byte)) (error)
}
