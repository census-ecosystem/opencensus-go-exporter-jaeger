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

func TestHTTPFormat(t *testing.T) {
	format := &HTTPFormat{}
	traceID := [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 66, 179, 103, 245, 105, 105, 242, 156}
	traceID2 := [16]byte{66, 179, 103, 245, 105, 105, 242, 156, 66, 179, 103, 245, 105, 105, 242, 156}
	spanID1 := [8]byte{104, 185, 184, 89, 243, 185, 19, 51}
	spanID2 := [8]byte{67, 211, 230, 84, 180, 39, 182, 139}
	tests := []struct {
		incoming        string
		wantSpanContext trace.SpanContext
	}{
		{
			incoming: "42b367f56969f29c:68b9b859f3b91333::1",
			wantSpanContext: trace.SpanContext{
				TraceID:      traceID,
				SpanID:       spanID1,
				TraceOptions: 1,
			},
		},
		{
			incoming: "42b367f56969f29c:43d3e654b427b68b::0",
			wantSpanContext: trace.SpanContext{
				TraceID:      traceID,
				SpanID:       spanID2,
				TraceOptions: 0,
			},
		},
		{
			incoming: "42b367f56969f29c42b367f56969f29c:43d3e654b427b68b::0",
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
			if !ok {
				t.Errorf("exporter.SpanContextFromRequest() = false; want true")
			}
			if got, want := sc, tt.wantSpanContext; !reflect.DeepEqual(got, want) {
				t.Errorf("exporter.SpanContextFromRequest() returned span context %v; want %v", got, want)
			}

			req, _ = http.NewRequest("GET", "http://example.com", nil)
			format.SpanContextToRequest(sc, req)

			if got, want := req.Header.Get(httpHeader), tt.incoming; got != want {
				t.Errorf("exporter.SpanContextToRequest() returned header %q; want %q", got, want)
			}
		})
	}
}
