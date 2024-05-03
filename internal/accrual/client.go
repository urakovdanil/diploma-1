package accrual

import (
	"context"
	"diploma-1/internal/closer"
	"diploma-1/internal/config"
	"diploma-1/internal/logger"
	"diploma-1/internal/storage"
	"diploma-1/internal/types"
	"encoding/json"
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/go-resty/resty/v2"
	"golang.org/x/time/rate"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	clientRateLimit              = 20
	clientBurst                  = 30
	defaultRetryAfterHeaderValue = "10"
	accrualURITemplate           = "/api/orders/%s"
)

type c struct {
	client            *resty.Client
	ordersToBeUpdated chan *order
	canSend           bool
	notificationChan  chan struct{}
	mu                sync.RWMutex
	requestsWg        sync.WaitGroup
}

func (cl *c) close() error {
	close(cl.ordersToBeUpdated)
	cl.requestsWg.Wait()
	return nil
}

func (cl *c) getCanSend() bool {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	return cl.canSend
}

func (cl *c) setCanSend(canSend bool) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.canSend = canSend
}

func (cl *c) freezeSending(ctx context.Context, retryAfterHeaderValue string) {
	go func() {
		if !cl.getCanSend() {
			return
		}
		if retryAfterHeaderValue == "" {
			logger.Warnf(ctx, "accrual system returned empty retry-after: set to %s", defaultRetryAfterHeaderValue)
			retryAfterHeaderValue = defaultRetryAfterHeaderValue
		}
		retryAfter, err := strconv.Atoi(retryAfterHeaderValue)
		if err != nil {
			logger.Warnf(ctx, "accrual system returned invalid retry-after %s: %v", retryAfterHeaderValue, err)
			retryAfter, _ = strconv.Atoi(defaultRetryAfterHeaderValue)
		}
		cl.setCanSend(false)
		time.Sleep(time.Duration(retryAfter) * time.Second)
		close(cl.notificationChan)
		cl.notificationChan = make(chan struct{})
		cl.setCanSend(true)
	}()
}

func (cl *c) checkUpdate(breakCtx context.Context, ord *order) {
	defer cl.requestsWg.Done()
	ctx := context.WithValue(context.Background(), types.CtxKeyRequestID, ord.requestID)
	bo := newBackoff()
mainLoop:
	for {
		select {
		case <-breakCtx.Done():
			fmt.Println("done")
			logger.Errorf(ctx, "tracking of order %s was canceled on shutdown", ord.Number)
			break mainLoop
		default:
		}
		nextBO := bo.NextBackOff()
		if nextBO == backoff.Stop {
			logger.Error(ctx, "accrual system did not meet deadline")
			break mainLoop
		}
		time.Sleep(nextBO)
		if !cl.getCanSend() {
			<-cl.notificationChan
		}
		ctx = context.WithValue(ctx, types.CtxUsedAccrualAddress, client.client.BaseURL+fmt.Sprintf(accrualURITemplate, ord.Number))
		req := cl.client.R().SetContext(ctx)
		resp, err := req.Get(fmt.Sprintf(accrualURITemplate, ord.Number))
		if err != nil {
			logger.Errorf(ctx, "unexpected error on request to accrual system: %v", err)
			continue mainLoop
		}
	statusCodeSwitch:
		switch resp.StatusCode() {
		case http.StatusInternalServerError:
			logger.Warn(ctx, "accrual system returned 500 error")
		case http.StatusTooManyRequests:
			cl.freezeSending(ctx, resp.Header().Get("Retry-After"))
		case http.StatusNoContent:
			logger.Debug(ctx, "accrual system returned 204")
		case http.StatusOK:
			var fromAccrual *types.OrderFromAccrual
			if err := json.Unmarshal(resp.Body(), &fromAccrual); err != nil {
				logger.Warnf(ctx, "unexpected error on unmarshal response from accrual system: %v", err)
				break statusCodeSwitch
			}
			if err := storage.UpdateOrderFromAccrual(ctx, fromAccrual); err != nil {
				logger.Warnf(ctx, err.Error())
				break statusCodeSwitch
			}
			if _, ok := types.FinalOrderStatuses[fromAccrual.Status]; ok {
				break mainLoop
			}
		default:
			logger.Warnf(ctx, "accrual system returned unexpected status code %d", resp.StatusCode())
		}
	}
}

func (cl *c) run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		fmt.Println("CANCELED")
		cancel()
	}()
	for ord := range cl.ordersToBeUpdated {
		logger.Debugf(ctx, "processing order %s from request %s", ord.Number, ord.requestID)
		cl.requestsWg.Add(1)
		go cl.checkUpdate(ctx, ord)
	}
}

var client *c

func New(ctx context.Context, su *config.StartUp) error {
	if _, err := strconv.Atoi(defaultRetryAfterHeaderValue); err != nil {
		return err
	}
	client = &c{
		ordersToBeUpdated: make(chan *order, 100),
		canSend:           true,
		notificationChan:  make(chan struct{}),
	}
	client.client = resty.New().
		SetBaseURL(su.GetAccrualSystemAddress()).
		SetRateLimiter(rate.NewLimiter(rate.Limit(clientRateLimit), clientBurst)).
		OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
			logger.Debugf(r.Context(), "sending '%v %v'", r.Method, r.URL)
			return nil
		}).
		OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
			logger.Debugf(context.Background(), "%v request %v (%v) took %v ms", time.Now().Format(time.DateTime), r.Request.URL, r.StatusCode(), r.Time().Milliseconds())
			return nil
		})
	go client.run()
	closer.Add(client.close)
	logger.Error(ctx, "initialized accrual client")
	return nil
}
