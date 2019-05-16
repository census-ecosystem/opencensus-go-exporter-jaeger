// Copyright 2019, OpenCensus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package propagation implement uber-trace-id header propagation used
// by Jaeger.
package propagation

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"
)

const (
	httpHeader   = `uber-trace-id`
	maxLenHeader = 200
)

var _ propagation.HTTPFormat = (*HTTPFormat)(nil)

// HTTPFormat implements propagation.HTTPFormat to propagate
// traces in HTTP headers for Jaeger traces.
type HTTPFormat struct{}

// SpanContextFromRequest extracts a Jaeger Trace span context from incoming requests.
func (f *HTTPFormat) SpanContextFromRequest(req *http.Request) (sc trace.SpanContext, ok bool) {
	h := req.Header.Get(httpHeader)

	if h == "" || len(h) > maxLenHeader {
		return trace.SpanContext{}, false
	}

	// Parse the trace id field.
	traceHeaderParts := strings.Split(h, `:`)
	if len(traceHeaderParts) != 4 {
		return trace.SpanContext{}, false
	}

	traceID, err := hex.DecodeString(traceHeaderParts[0])
	if err != nil {
		return trace.SpanContext{}, false
	}
	if len(traceID) == 8 {
		copy(sc.TraceID[8:16], traceID)
	} else {
		copy(sc.TraceID[:], traceID)
	}

	spanID, err := hex.DecodeString(traceHeaderParts[1])
	if err != nil {
		return trace.SpanContext{}, false
	}
	copy(sc.SpanID[:], spanID)

	opt, err := strconv.Atoi(traceHeaderParts[3])

	if err != nil {
		return trace.SpanContext{}, false
	}

	sc.TraceOptions = trace.TraceOptions(opt)

	return sc, true
}

// SpanContextToRequest modifies the given request to include a Jaeger Trace header.
func (f *HTTPFormat) SpanContextToRequest(sc trace.SpanContext, req *http.Request) {
	header := fmt.Sprintf("%s:%s:%s:%d",
		strings.Replace(sc.TraceID.String(), "0000000000000000", "", 1), //Replacing 0 if string is 8bit
		sc.SpanID.String(),
		"", //Parent span deprecated and will therefore be ignored.
		int64(sc.TraceOptions))
	req.Header.Set(httpHeader, header)
}
