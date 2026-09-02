/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package modelbuilder

import (
	"context"
	"testing"
)

func TestWithThinking(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		t.Run(map[bool]string{false: "disabled", true: "enabled"}[enabled], func(t *testing.T) {
			params := &LLMParams{}
			got := withThinking(params, enabled)

			if got != params {
				t.Fatal("withThinking should reuse the supplied params")
			}
			if got.EnableThinking == nil || *got.EnableThinking != enabled {
				t.Fatalf("unexpected thinking setting: %#v", got.EnableThinking)
			}
		})
	}
}

func TestWithThinkingAllocatesParams(t *testing.T) {
	params := withThinking(nil, false)
	if params == nil || params.EnableThinking == nil || *params.EnableThinking {
		t.Fatalf("expected allocated params with disabled thinking, got %#v", params)
	}
}

func TestBuildModelBySettingsWithThinkingRejectsNilSettings(t *testing.T) {
	_, _, err := BuildModelBySettingsWithThinking(context.Background(), nil, false)
	if err == nil {
		t.Fatal("expected nil model settings to be rejected")
	}
}
