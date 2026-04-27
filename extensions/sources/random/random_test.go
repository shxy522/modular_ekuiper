// Copyright 2021-2023 EMQ Technologies Co., Ltd.
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

package random

import (
	"math"
	"testing"
)

func TestRandomizePrecisionZeroKeepsInteger(t *testing.T) {
	got := randomize(map[string]interface{}{"count": 50}, 1, 0)

	v, ok := got["count"].(int)
	if !ok {
		t.Fatalf("expected int value, got %[1]T(%[1]v)", got["count"])
	}
	if v != 50 {
		t.Fatalf("expected 50, got %d", v)
	}
}

func TestRandomizePrecisionGeneratesFloat(t *testing.T) {
	tests := []struct {
		name      string
		precision int
		scale     float64
	}{
		{
			name:      "one decimal place",
			precision: 1,
			scale:     10,
		},
		{
			name:      "two decimal places",
			precision: 2,
			scale:     100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := randomize(map[string]interface{}{"count": 50}, 2, tt.precision)

			v, ok := got["count"].(float64)
			if !ok {
				t.Fatalf("expected float64 value, got %[1]T(%[1]v)", got["count"])
			}
			if v < 50 || v > 52 {
				t.Fatalf("expected value in [50, 52], got %v", v)
			}
			scaled := v * tt.scale
			if math.Abs(scaled-math.Round(scaled)) > 1e-9 {
				t.Fatalf("expected value with precision %d, got %v", tt.precision, v)
			}
		})
	}
}

func TestConfigureRejectsInvalidPrecision(t *testing.T) {
	r := &randomSource{}

	err := r.Configure("", map[string]interface{}{
		"interval":  1000,
		"seed":      1,
		"precision": -1,
		"pattern": map[string]interface{}{
			"count": 50,
		},
	})

	if err == nil {
		t.Fatal("expected error for negative precision")
	}
}
