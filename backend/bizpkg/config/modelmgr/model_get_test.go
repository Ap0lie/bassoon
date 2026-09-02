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

package modelmgr

import (
	"testing"

	config "github.com/coze-dev/coze-studio/backend/api/model/admin/config"
)

func TestNormalizeMoonshotKimiK26Connection(t *testing.T) {
	tests := []struct {
		name string
		conn *config.Connection
		want bool
	}{
		{
			name: "supported Moonshot China endpoint",
			conn: &config.Connection{BaseConnInfo: &config.BaseConnectionInfo{
				BaseURL: moonshotKimiK26BaseURL,
				Model:   moonshotKimiK26Model,
			}},
			want: true,
		},
		{
			name: "trailing slash is normalized",
			conn: &config.Connection{BaseConnInfo: &config.BaseConnectionInfo{
				BaseURL: moonshotKimiK26BaseURL + "/",
				Model:   moonshotKimiK26Model,
			}},
			want: true,
		},
		{
			name: "chat completions resource is normalized to the base endpoint",
			conn: &config.Connection{BaseConnInfo: &config.BaseConnectionInfo{
				BaseURL: moonshotKimiK26BaseURL + "/chat/completions",
				Model:   moonshotKimiK26Model,
			}},
			want: true,
		},
		{
			name: "another model is unchanged",
			conn: &config.Connection{BaseConnInfo: &config.BaseConnectionInfo{
				BaseURL: moonshotKimiK26BaseURL,
				Model:   "another-model",
			}},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalBaseURL := tt.conn.BaseConnInfo.BaseURL
			if got := normalizeMoonshotKimiK26Connection(tt.conn); got != tt.want {
				t.Fatalf("normalizeMoonshotKimiK26Connection() = %v, want %v", got, tt.want)
			}
			if tt.want && tt.conn.BaseConnInfo.BaseURL != moonshotKimiK26BaseURL {
				t.Fatalf("base URL = %q, want %q", tt.conn.BaseConnInfo.BaseURL, moonshotKimiK26BaseURL)
			}
			if !tt.want && tt.conn.BaseConnInfo.BaseURL != originalBaseURL {
				t.Fatalf("base URL = %q, want unchanged %q", tt.conn.BaseConnInfo.BaseURL, originalBaseURL)
			}
		})
	}
}
