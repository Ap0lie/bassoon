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

package agentflow

import (
	"context"
	"errors"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/eino/schema"
	"github.com/stretchr/testify/assert"

	config "github.com/coze-dev/coze-studio/backend/api/model/admin/config"
	"github.com/coze-dev/coze-studio/backend/api/model/app/developer_api"
	"github.com/coze-dev/coze-studio/backend/bizpkg/config/modelmgr"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-studio/backend/pkg/urltobase64url"
)

func TestPreHandlerReqKeepsURLWhenBase64TransferDisabled(t *testing.T) {
	runner := &AgentRunner{modelInfo: &modelmgr.Model{Model: &config.Model{
		Capability: &developer_api.ModelAbility{
			ImageUnderstanding: ptr.Of(true),
		},
		EnableBase64URL: false,
	}}}

	preparedReq, err := runner.PreHandlerReq(context.Background(), &AgentRequest{
		Input: &schema.Message{MultiContent: []schema.ChatMessagePart{{
			Type:     schema.ChatMessagePartTypeImageURL,
			ImageURL: &schema.ChatMessageImageURL{URL: "https://example.com/image.png"},
		}}},
	})

	if assert.NoError(t, err) && assert.NotNil(t, preparedReq) && assert.NotNil(t, preparedReq.Input) &&
		assert.Len(t, preparedReq.Input.MultiContent, 1) {
		assert.Equal(t, "https://example.com/image.png", preparedReq.Input.MultiContent[0].ImageURL.URL)
	}
}

func TestMediaReferencesAreKeptWhenBase64TransferEnabled(t *testing.T) {
	t.Run("data URL", func(t *testing.T) {
		imageURL := &schema.ChatMessageImageURL{URL: "data:image/png;base64,base64encodedstring"}

		result, err := transImageURLToBase64(imageURL, true)

		assert.NoError(t, err)
		assert.Same(t, imageURL, result)
	})

	t.Run("Moonshot file reference", func(t *testing.T) {
		videoURL := &schema.ChatMessageVideoURL{URL: "ms://file-id"}

		result, err := transVideoURLToBase64(videoURL, true)

		assert.NoError(t, err)
		assert.Same(t, videoURL, result)
	})
}

func TestNonBase64DataURLIsRejectedWhenTransferEnabled(t *testing.T) {
	mockey.PatchConvey("TestNonBase64DataURLIsRejectedWhenTransferEnabled", t, func() {
		mockey.Mock(urltobase64url.URLToBase64).Return((*urltobase64url.FileData)(nil), errors.New("download failed")).Build()

		result, err := transImageURLToBase64(&schema.ChatMessageImageURL{
			URL: "data:image/png,not-base64",
		}, true)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, "convert image URL to base64")
	})
}

func TestPreHandlerReqRejectsUnconvertedHTTPMedia(t *testing.T) {
	mockey.PatchConvey("TestPreHandlerReqRejectsUnconvertedHTTPMedia", t, func() {
		mockey.Mock(urltobase64url.URLToBase64).Return((*urltobase64url.FileData)(nil), errors.New("download failed")).Build()

		tests := []struct {
			name       string
			capability *developer_api.ModelAbility
			part       schema.ChatMessagePart
		}{
			{
				name: "image",
				capability: &developer_api.ModelAbility{
					ImageUnderstanding: ptr.Of(true),
				},
				part: schema.ChatMessagePart{
					Type:     schema.ChatMessagePartTypeImageURL,
					ImageURL: &schema.ChatMessageImageURL{URL: "https://example.com/image.png"},
				},
			},
			{
				name: "video",
				capability: &developer_api.ModelAbility{
					VideoUnderstanding: ptr.Of(true),
				},
				part: schema.ChatMessagePart{
					Type:     schema.ChatMessagePartTypeVideoURL,
					VideoURL: &schema.ChatMessageVideoURL{URL: "https://example.com/video.mp4"},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				runner := &AgentRunner{modelInfo: &modelmgr.Model{Model: &config.Model{
					Capability:      tt.capability,
					EnableBase64URL: true,
				}}}

				_, err := runner.PreHandlerReq(context.Background(), &AgentRequest{
					Input: &schema.Message{MultiContent: []schema.ChatMessagePart{tt.part}},
				})

				assert.ErrorContains(t, err, "convert "+tt.name+" URL to base64")
			})
		}
	})
}
