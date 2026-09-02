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

package coze

import (
	"testing"

	"github.com/coze-dev/coze-studio/backend/api/model/admin/config"
)

func TestNormalizeMoonshotKimiK26Connection(t *testing.T) {
	t.Run("normalizes the supported China endpoint", func(t *testing.T) {
		conn := &config.Connection{BaseConnInfo: &config.BaseConnectionInfo{
			BaseURL: moonshotKimiK26BaseURL + "/",
			Model:   moonshotKimiK26Model,
		}}

		matched, err := normalizeMoonshotKimiK26Connection(conn)
		if err != nil || !matched {
			t.Fatalf("normalizeMoonshotKimiK26Connection() = (%v, %v), want (true, nil)", matched, err)
		}
		if conn.BaseConnInfo.BaseURL != moonshotKimiK26BaseURL {
			t.Fatalf("base URL = %q, want %q", conn.BaseConnInfo.BaseURL, moonshotKimiK26BaseURL)
		}
	})

	t.Run("rejects a chat completions resource URL", func(t *testing.T) {
		conn := &config.Connection{BaseConnInfo: &config.BaseConnectionInfo{
			BaseURL: moonshotKimiK26BaseURL + "/chat/completions",
			Model:   moonshotKimiK26Model,
		}}

		matched, err := normalizeMoonshotKimiK26Connection(conn)
		if !matched || err == nil {
			t.Fatalf("normalizeMoonshotKimiK26Connection() = (%v, %v), want (true, error)", matched, err)
		}
	})

	t.Run("does not change another model", func(t *testing.T) {
		conn := &config.Connection{BaseConnInfo: &config.BaseConnectionInfo{
			BaseURL: "https://example.com/v1",
			Model:   "another-model",
		}}

		matched, err := normalizeMoonshotKimiK26Connection(conn)
		if err != nil || matched {
			t.Fatalf("normalizeMoonshotKimiK26Connection() = (%v, %v), want (false, nil)", matched, err)
		}
	})
}
