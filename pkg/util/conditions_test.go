package util

import (
	"testing"

	"time"

	"github.com/harvester/vm-import-controller/pkg/apis/common"
	importjob "github.com/harvester/vm-import-controller/pkg/apis/importjob.harvesterhci.io/v1beta1"
	source "github.com/harvester/vm-import-controller/pkg/apis/source.harvesterhci.io/v1beta1"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_ConditionExists(t *testing.T) {
	conditions := []common.Condition{
		{
			Type:               source.ClusterReadyCondition,
			Status:             corev1.ConditionTrue,
			LastUpdateTime:     metav1.Now().Format(time.RFC3339),
			LastTransitionTime: metav1.Now().Format(time.RFC3339),
		},
		{
			Type:               source.ClusterErrorCondition,
			Status:             corev1.ConditionFalse,
			LastUpdateTime:     metav1.Now().Format(time.RFC3339),
			LastTransitionTime: metav1.Now().Format(time.RFC3339),
		},
	}

	assert := require.New(t)
	assert.True(ConditionExists(conditions, source.ClusterReadyCondition, corev1.ConditionTrue))
	assert.True(ConditionExists(conditions, source.ClusterErrorCondition, corev1.ConditionFalse))
}

func Test_AddOrUpdateCondition(t *testing.T) {
	conditions := []common.Condition{
		{
			Type:               source.ClusterReadyCondition,
			Status:             corev1.ConditionTrue,
			LastUpdateTime:     metav1.Now().Format(time.RFC3339),
			LastTransitionTime: metav1.Now().Format(time.RFC3339),
		},
		{
			Type:               source.ClusterErrorCondition,
			Status:             corev1.ConditionFalse,
			LastUpdateTime:     metav1.Now().Format(time.RFC3339),
			LastTransitionTime: metav1.Now().Format(time.RFC3339),
		},
	}

	extraCondition := common.Condition{

		Type:               importjob.VirtualMachinePoweringOff,
		Status:             corev1.ConditionTrue,
		LastUpdateTime:     metav1.Now().Format(time.RFC3339),
		LastTransitionTime: metav1.Now().Format(time.RFC3339),
	}

	newCond := AddOrUpdateCondition(conditions, extraCondition)
	assert := require.New(t)
	assert.True(ConditionExists(newCond, importjob.VirtualMachinePoweringOff, corev1.ConditionTrue))
	assert.True(ConditionExists(conditions, source.ClusterErrorCondition, corev1.ConditionFalse))
	assert.True(ConditionExists(conditions, source.ClusterReadyCondition, corev1.ConditionTrue))
}

func Test_MergeConditions(t *testing.T) {
	conditions := []common.Condition{
		{
			Type:               source.ClusterReadyCondition,
			Status:             corev1.ConditionTrue,
			LastUpdateTime:     metav1.Now().Format(time.RFC3339),
			LastTransitionTime: metav1.Now().Format(time.RFC3339),
		},
		{
			Type:               source.ClusterErrorCondition,
			Status:             corev1.ConditionFalse,
			LastUpdateTime:     metav1.Now().Format(time.RFC3339),
			LastTransitionTime: metav1.Now().Format(time.RFC3339),
		},
	}

	extraConditions := []common.Condition{
		{
			Type:               importjob.VirtualMachineExported,
			Status:             corev1.ConditionTrue,
			LastUpdateTime:     metav1.Now().Format(time.RFC3339),
			LastTransitionTime: metav1.Now().Format(time.RFC3339),
		},
		{
			Type:               importjob.VirtualMachineImageReady,
			Status:             corev1.ConditionTrue,
			LastUpdateTime:     metav1.Now().Format(time.RFC3339),
			LastTransitionTime: metav1.Now().Format(time.RFC3339),
		},
	}

	newConds := MergeConditions(conditions, extraConditions)
	assert := require.New(t)
	assert.Len(newConds, 4, "expected to find 4 conditions in the merged conditions")
}

func Test_RemoveCondition(t *testing.T) {
	conditions := []common.Condition{
		{
			Type:               source.ClusterReadyCondition,
			Status:             corev1.ConditionTrue,
			LastUpdateTime:     metav1.Now().Format(time.RFC3339),
			LastTransitionTime: metav1.Now().Format(time.RFC3339),
		},
		{
			Type:               source.ClusterErrorCondition,
			Status:             corev1.ConditionFalse,
			LastUpdateTime:     metav1.Now().Format(time.RFC3339),
			LastTransitionTime: metav1.Now().Format(time.RFC3339),
		},
	}

	noRemoveCond := RemoveCondition(conditions, source.ClusterErrorCondition, corev1.ConditionTrue)
	assert := require.New(t)
	assert.True(ConditionExists(noRemoveCond, source.ClusterErrorCondition, corev1.ConditionFalse))
	removeCond := RemoveCondition(conditions, source.ClusterErrorCondition, corev1.ConditionFalse)
	assert.False(ConditionExists(removeCond, source.ClusterErrorCondition, corev1.ConditionFalse))
}
