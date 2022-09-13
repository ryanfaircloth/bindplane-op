// Copyright  observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/observiq/bindplane-op/internal/otlp/record"
)

// HeaderSessionID is the name of the HTTP header where a BindPlane Session ID can be found. This will be used by
// Relayer to match up the request with an eventual response via OTLP HTTP POST.
const HeaderSessionID = "X-BindPlane-Session-ID"

// Relayers is a wrapper around multiple Relayer instances used for different types of results
type Relayers struct {
	Logs    *Relayer[[]*record.Log]
	Metrics *Relayer[[]*record.Metric]
	Traces  *Relayer[[]*record.Trace]
}

// NewRelayers returns a new set of Relayers
func NewRelayers(logger *zap.Logger) *Relayers {
	return &Relayers{
		Logs:    newRelayer[[]*record.Log](logger),
		Metrics: newRelayer[[]*record.Metric](logger),
		Traces:  newRelayer[[]*record.Trace](logger),
	}
}

// Relayer forwards results to consumers awaiting the results. It is intentionally generic and is used to support cases
// where the request for results is decoupled from the response with the results.
type Relayer[T any] struct {
	subscriptions map[string]subscription[T]
	mtx           sync.Mutex
	logger        *zap.Logger
}

// newRelayer will build a new Relayer for use with results of the specified type T.
func newRelayer[T any](logger *zap.Logger) *Relayer[T] {
	return &Relayer[T]{
		subscriptions: map[string]subscription[T]{},
		logger:        logger,
	}
}

type subscription[T any] struct {
	result chan T
}

// AwaitResult will create a new ID and channel where results will be received. It also returns a cancelFunc that must
// be called to cleanup the channel.
func (r *Relayer[T]) AwaitResult() (id string, result <-chan T, cancelFunc func()) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	id = uuid.NewString()
	resultChannel := make(chan T, 1)

	r.subscriptions[id] = subscription[T]{
		result: resultChannel,
	}

	cancelFunc = func() {
		r.mtx.Lock()
		defer r.mtx.Unlock()
		close(resultChannel)
		delete(r.subscriptions, id)
	}

	return id, resultChannel, cancelFunc
}

// SendResult will forward results to the channel returned by AwaitResult
func (r *Relayer[T]) SendResult(id string, result T) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if sub, ok := r.subscriptions[id]; ok {
		sub.result <- result
	} else {
		r.logger.Debug("received results but there is nothing waiting", zap.String("id", id))
	}
}
