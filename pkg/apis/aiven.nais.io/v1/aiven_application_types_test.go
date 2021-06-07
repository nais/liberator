package aiven_nais_io_v1

import (
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

var (
	successTrue = AivenApplicationCondition{
		Type:   AivenApplicationSucceeded,
		Status: corev1.ConditionTrue,
	}
	successFalse = AivenApplicationCondition{
		Type:   AivenApplicationSucceeded,
		Status: corev1.ConditionFalse,
	}
	localFailTrue = AivenApplicationCondition{
		Type:   AivenApplicationLocalFailure,
		Status: corev1.ConditionTrue,
	}
	localFailFalse = AivenApplicationCondition{
		Type:   AivenApplicationLocalFailure,
		Status: corev1.ConditionFalse,
	}
	aivenFailTrue = AivenApplicationCondition{
		Type:   AivenApplicationAivenFailure,
		Status: corev1.ConditionTrue,
	}
	aivenFailFalse = AivenApplicationCondition{
		Type:   AivenApplicationAivenFailure,
		Status: corev1.ConditionFalse,
	}
)

func TestAivenApplicationStatus_AddCondition(t *testing.T) {
	type args struct {
		preExisting  []AivenApplicationCondition
		newCondition AivenApplicationCondition
		wanted       []AivenApplicationCondition
		dropTypes    []AivenApplicationConditionType
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "EmptyAddSuccessFalse", args: args{
			preExisting:  nil,
			newCondition: successFalse,
			wanted:       []AivenApplicationCondition{successFalse},
		}},
		{name: "ToggleSuccess", args: args{
			preExisting:  []AivenApplicationCondition{successTrue},
			newCondition: successFalse,
			wanted:       []AivenApplicationCondition{successFalse},
		}},
		{name: "ExistingAddSuccessFalse", args: args{
			preExisting:  []AivenApplicationCondition{localFailTrue},
			newCondition: successFalse,
			wanted:       []AivenApplicationCondition{localFailTrue, successFalse},
		}},
		{name: "FilledAddSuccessTrue", args: args{
			preExisting:  []AivenApplicationCondition{localFailFalse, successFalse},
			newCondition: successTrue,
			wanted:       []AivenApplicationCondition{localFailFalse, successTrue},
		}},
		{name: "FilledAddSuccessTrueBackwards", args: args{
			preExisting:  []AivenApplicationCondition{successFalse, localFailFalse},
			newCondition: successTrue,
			wanted:       []AivenApplicationCondition{localFailFalse, successTrue},
		}},
		{name: "FilledAddFailureRemoveSuccess", args: args{
			preExisting:  []AivenApplicationCondition{localFailFalse, successTrue},
			newCondition: localFailTrue,
			wanted:       []AivenApplicationCondition{localFailTrue},
			dropTypes:    []AivenApplicationConditionType{AivenApplicationSucceeded},
		}},
		{name: "FilledFullRemoveMost", args: args{
			preExisting:  []AivenApplicationCondition{localFailTrue, aivenFailTrue, successFalse},
			newCondition: successTrue,
			wanted:       []AivenApplicationCondition{successTrue},
			dropTypes:    []AivenApplicationConditionType{AivenApplicationAivenFailure, AivenApplicationLocalFailure},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertConditionsMatch := func(expected, actual AivenApplicationCondition) {
				a := assert.New(t)
				a.Equal(expected.Type, actual.Type)
				a.Equal(expected.Status, actual.Status)
			}
			in := &AivenApplicationStatus{
				Conditions: tt.args.preExisting,
			}
			in.AddCondition(tt.args.newCondition, tt.args.dropTypes...)
			assert.Len(t, in.Conditions, len(tt.args.wanted))
			for i, c := range in.Conditions {
				assertConditionsMatch(tt.args.wanted[i], c)
			}
		})
	}
}
