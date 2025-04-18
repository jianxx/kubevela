/*
Copyright 2022 The KubeVela Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package view

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	types "github.com/oam-dev/kubevela/apis/core.oam.dev/common"
	"github.com/oam-dev/kubevela/pkg/utils/common"
	"github.com/oam-dev/kubevela/references/cli/top/model"
)

func TestApplicationView(t *testing.T) {
	testEnv := &envtest.Environment{
		ControlPlaneStartTimeout: time.Minute * 3,
		ControlPlaneStopTimeout:  time.Minute,
		UseExistingCluster:       ptr.To(false),
	}
	cfg, err := testEnv.Start()
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, testEnv.Stop())
	}()

	testClient, err := client.New(cfg, client.Options{Scheme: common.Scheme})
	assert.NoError(t, err)
	app := NewApp(testClient, cfg, "")
	assert.Equal(t, len(app.Components()), 4)

	ctx := context.Background()
	ctx = context.WithValue(ctx, &model.CtxKeyNamespace, "")

	appView := new(ApplicationView)

	t.Run("init view", func(t *testing.T) {
		assert.Empty(t, appView.CommonResourceView)
		appView.InitView(ctx, app)
		assert.NotEmpty(t, appView.CommonResourceView)
	})

	t.Run("init", func(t *testing.T) {
		appView.Init()
		assert.Equal(t, appView.Table.GetTitle(), "[ Application (all) ]")
	})

	t.Run("refresh", func(t *testing.T) {
		keyEvent := appView.Refresh(nil)
		assert.Empty(t, keyEvent)
	})

	t.Run("start", func(t *testing.T) {
		appView.Start()
		assert.Equal(t, appView.GetCell(0, 0).Text, "Name")
	})

	t.Run("stop", func(t *testing.T) {
		appView.Stop()
		assert.Equal(t, appView.GetCell(0, 0).Text, "")
	})

	t.Run("colorize text", func(t *testing.T) {
		testData := [][]string{
			{"app", "ns", "starting", ""},
			{"app", "ns", "rendering", ""},
			{"app", "ns", "generatingPolicy", ""},
			{"app", "ns", "runningWorkflow", ""},
			{"app", "ns", "workflowSuspending", ""},
			{"app", "ns", "workflowTerminated", ""},
			{"app", "ns", "workflowFailed", ""},
			{"app", "ns", "unhealthy", ""},
			{"app", "ns", "deleting", ""},
			{"app", "ns", "running", ""},
		}
		for i := 0; i < len(testData); i++ {
			for j := 0; j < 4; j++ {
				appView.Table.SetCell(1+i, j, tview.NewTableCell(testData[i][j]))
			}
		}
		appView.ColorizeStatusText(10)
		assert.Equal(t, appView.GetCell(1, 2).Text, fmt.Sprintf("[%s::]%s", appView.app.config.Theme.Status.Starting.String(), types.ApplicationStarting))
		assert.Equal(t, appView.GetCell(2, 2).Text, fmt.Sprintf("[%s::]%s", appView.app.config.Theme.Status.Waiting.String(), types.ApplicationRendering))
		assert.Equal(t, appView.GetCell(3, 2).Text, fmt.Sprintf("[%s::]%s", appView.app.config.Theme.Status.Waiting.String(), types.ApplicationPolicyGenerating))
		assert.Equal(t, appView.GetCell(4, 2).Text, fmt.Sprintf("[%s::]%s", appView.app.config.Theme.Status.Waiting.String(), types.ApplicationRunningWorkflow))
		assert.Equal(t, appView.GetCell(5, 2).Text, fmt.Sprintf("[%s::]%s", appView.app.config.Theme.Status.Waiting.String(), types.ApplicationWorkflowSuspending))
		assert.Equal(t, appView.GetCell(6, 2).Text, fmt.Sprintf("[%s::]%s", appView.app.config.Theme.Status.Failed.String(), types.ApplicationWorkflowTerminated))
		assert.Equal(t, appView.GetCell(7, 2).Text, fmt.Sprintf("[%s::]%s", appView.app.config.Theme.Status.Failed.String(), types.ApplicationWorkflowFailed))
		assert.Equal(t, appView.GetCell(8, 2).Text, fmt.Sprintf("[%s::]%s", appView.app.config.Theme.Status.Failed.String(), types.ApplicationUnhealthy))
		assert.Equal(t, appView.GetCell(9, 2).Text, fmt.Sprintf("[%s::]%s", appView.app.config.Theme.Status.Failed.String(), types.ApplicationDeleting))
		assert.Equal(t, appView.GetCell(10, 2).Text, fmt.Sprintf("[%s::]%s", appView.app.config.Theme.Status.Healthy.String(), types.ApplicationRunning))
	})

	t.Run("hint", func(t *testing.T) {
		assert.Equal(t, len(appView.Hint()), 9)
	})

	t.Run("managed resource view", func(t *testing.T) {
		appView.Table.Table = appView.Table.Select(1, 1)
		assert.Empty(t, appView.managedResourceView(nil))
	})

	t.Run("namespace view", func(t *testing.T) {
		appView.Table.Table = appView.Table.Select(1, 1)
		assert.Empty(t, appView.namespaceView(nil))
	})

	t.Run("yaml view", func(t *testing.T) {
		appView.Table.Table = appView.Table.Select(1, 1)
		assert.Empty(t, appView.managedResourceView(nil))
	})

}
