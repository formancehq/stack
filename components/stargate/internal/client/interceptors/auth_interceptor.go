package interceptors

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	defaultWaitingTime = 10 * time.Second
)

type Config struct {
	refreshTokenDurationBeforeExpireTime time.Duration

	clientID     string
	clientSecret string
	endpoint     string
}

func NewConfig(
	endpoint string,
	refreshTokenDurationBeforeExpireTime time.Duration,
	clientID string,
	clientSecret string,
) Config {
	return Config{
		refreshTokenDurationBeforeExpireTime: refreshTokenDurationBeforeExpireTime,
		clientID:                             clientID,
		clientSecret:                         clientSecret,
		endpoint:                             endpoint,
	}
}

type AuthInterceptor struct {
	config Config

	accessToken string
	closeChan   chan struct{}
}

func NewAuthInterceptor(config Config) (*AuthInterceptor, error) {
	i := &AuthInterceptor{
		config:    config,
		closeChan: make(chan struct{}),
	}

	return i, nil
}

func (a *AuthInterceptor) Close() {
	close(a.closeChan)
}

func (a *AuthInterceptor) StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		return streamer(
			metadata.AppendToOutgoingContext(ctx, "authorization", a.accessToken),
			desc,
			cc,
			method,
			opts...,
		)
	}
}

func (a *AuthInterceptor) ScheduleRefreshToken() error {
	expire, err := a.refreshToken()
	if err != nil {
		return err
	}

	go func() {
		waitingTime := time.Until(expire.Add(-a.config.refreshTokenDurationBeforeExpireTime))
		if waitingTime < 0 {
			waitingTime = defaultWaitingTime
		}
		for {
			select {
			case <-a.closeChan:
				return
			case <-time.After(waitingTime):
				expire, err := a.refreshToken()
				if err != nil {
					// TODO(polo): add metrics + log
					waitingTime = time.Second
				} else {
					waitingTime = time.Until(expire.Add(-a.config.refreshTokenDurationBeforeExpireTime))
					if waitingTime < 0 {
						waitingTime = defaultWaitingTime
					}
				}
			}
		}
	}()

	return nil
}

func (a *AuthInterceptor) refreshToken() (time.Time, error) {
	config := clientcredentials.Config{
		ClientID:     a.config.clientID,
		ClientSecret: a.config.clientSecret,
		TokenURL:     a.config.endpoint,
	}

	token, err := config.Token(context.Background())
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "cannot fetch token")
	}

	a.accessToken = token.AccessToken

	return token.Expiry, nil
}
