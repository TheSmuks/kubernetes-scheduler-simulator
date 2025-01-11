package plugin

import (
	"context"
	"fmt"
	simontype "github.com/hkust-adsl/kubernetes-scheduler-simulator/pkg/type"
	"github.com/hkust-adsl/kubernetes-scheduler-simulator/pkg/utils"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

type DummyScorePlugin struct {
	handle framework.Handle
}

var _ framework.ScorePlugin = &DummyScorePlugin{}

func NewDummyPlugin(configuration runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	return &DummyScorePlugin{
		handle: handle,
	}, nil
}

func (plugin *DummyScorePlugin) Name() string {
	return simontype.DummyScorePluginName
}

func (plugin *DummyScorePlugin) Score(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeName string) (int64, *framework.Status) {
	// < common procedure that prepares node, podRes, nodeRes>
	nodeResPtr := utils.GetNodeResourceViaHandleAndName(plugin.handle, nodeName)
	if nodeResPtr == nil {
		return framework.MinNodeScore, framework.NewStatus(framework.Error, fmt.Sprintf("failed to get nodeRes(%s)\n", nodeName))
	}
	nodeRes := *nodeResPtr

	podRes := utils.GetPodResource(pod)
	if !utils.IsNodeAccessibleToPod(nodeRes, podRes) {
		log.Errorf("Node (%s) %s does not match GPU type request of pod %s. Should be filtered by GpuSharePlugin", nodeName, nodeRes.Repr(), podRes.Repr())
		return framework.MinNodeScore, framework.NewStatus(framework.Error, fmt.Sprintf("Node (%s) %s does not match GPU type request of pod %s\n", nodeName, nodeRes.Repr(), podRes.Repr()))
	}
	// </common procedure that prepares node, podRes, nodeRes>

	score := getDummyScore(nodeRes, podRes)
	if score == -1 {
		return framework.MinNodeScore, framework.NewStatus(framework.Error, fmt.Sprintf("the score between node(%s) and pod(%s) is negative, should not happen\n", nodeName, utils.GeneratePodKey(pod)))
	}
	return score, framework.NewStatus(framework.Success)
}

func (plugin *DummyScorePlugin) ScoreExtensions() framework.ScoreExtensions {
	return plugin
}

func (plugin *DummyScorePlugin) NormalizeScore(ctx context.Context, state *framework.CycleState, p *corev1.Pod, scores framework.NodeScoreList) *framework.Status {
	return NormalizeScore(scores)
}

func getDummyScore(nodeRes simontype.NodeResource, podRes simontype.PodResource) int64 {
	return 0
}
