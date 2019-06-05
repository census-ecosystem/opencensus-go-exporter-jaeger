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

package propagation

import (
	"net/http"
	"reflect"
	"testing"

	"go.opencensus.io/trace"
)

func TestHTTPFormat_SpanContextFromRequest(t *testing.T) {
	format := &HTTPFormat{}
	traceID := [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 66, 179, 103, 245, 105, 105, 242, 156}
	traceID2 := [16]byte{66, 179, 103, 245, 105, 105, 242, 156, 66, 179, 103, 245, 105, 105, 242, 156}
	spanID1 := [8]byte{104, 185, 184, 89, 243, 185, 19, 51}
	spanID2 := [8]byte{67, 211, 230, 84, 180, 39, 182, 139}
	tests := []struct {
		incoming        string
		isPresent       bool
		wantSpanContext trace.SpanContext
	}{
		{
			incoming:  "42b367f56969f29c:68b9b859f3b91333::1",
			isPresent: true,
			wantSpanContext: trace.SpanContext{
				TraceID:      traceID,
				SpanID:       spanID1,
				TraceOptions: 1,
			},
		},
		{
			incoming:  "42b367f56969f29c:68b9b859f3b91333:1",
			isPresent: true,
			wantSpanContext: trace.SpanContext{
				TraceID:      traceID,
				SpanID:       spanID1,
				TraceOptions: 1,
			},
		},
		{
			incoming:  "42b367f56969f29c42b367f56969f29c:43d3e654b427b68b::1",
			isPresent: true,
			wantSpanContext: trace.SpanContext{
				TraceID:      traceID2,
				SpanID:       spanID2,
				TraceOptions: 1,
			},
		},
		{
			incoming:  "42b367f56969f29c42b367f56969f29c:43d3e654b427b68b::0",
			isPresent: true,
			wantSpanContext: trace.SpanContext{
				TraceID:      traceID2,
				SpanID:       spanID2,
				TraceOptions: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.incoming, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			req.Header.Add(httpHeader, tt.incoming)
			sc, ok := format.SpanContextFromRequest(req)
			if ok != tt.isPresent {
				t.Errorf("exporter.SpanContextFromRequest() = %v; want %v", ok, tt.isPresent)
			}
			if !ok {
				return
			}
			if got, want := sc, tt.wantSpanContext; !reflect.DeepEqual(got, want) {
				t.Errorf("exporter.SpanContextFromRequest() returned span context %v; want %v", got, want)
			}
		})
	}
}

func TestHTTPFormat_SpanContextToRequest(t *testing.T) {
	format := &HTTPFormat{}
	traceID := [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 66, 179, 103, 245, 105, 105, 242, 156}
	traceID2 := [16]byte{66, 179, 103, 245, 105, 105, 242, 156, 66, 179, 103, 245, 105, 105, 242, 156}
	spanID1 := [8]byte{104, 185, 184, 89, 243, 185, 19, 51}
	spanID2 := [8]byte{67, 211, 230, 84, 180, 39, 182, 139}
	tests := []struct {
		spanContext trace.SpanContext
		outbound    string
	}{
		{
			spanContext: trace.SpanContext{
				TraceID:      traceID,
				SpanID:       spanID1,
				TraceOptions: 1,
			},
			outbound: "42b367f56969f29c:68b9b859f3b91333::1",
		},
		{
			spanContext: trace.SpanContext{
				TraceID:      traceID,
				SpanID:       spanID1,
				TraceOptions: 0,
			},
			outbound: "42b367f56969f29c:68b9b859f3b91333::0",
		},
		{
			spanContext: trace.SpanContext{
				TraceID:      traceID2,
				SpanID:       spanID2,
				TraceOptions: 1,
			},
			outbound: "42b367f56969f29c42b367f56969f29c:43d3e654b427b68b::1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.outbound, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			format.SpanContextToRequest(tt.spanContext, req)

			if got, want := req.Header.Get(httpHeader), tt.outbound; got != want {
				t.Errorf("exporter.SpanContextToRequest() returned header %q; want %q", got, want)
			}
		})
	}
}
