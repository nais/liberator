package kafka_nais_io_v1

import (
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"testing"
)

var (
	successTrue = AivenApplicationCondition{
		Type:   AivenApplicationSucceeded,
		Status: v1.ConditionTrue,
	}
	successFalse = AivenApplicationCondition{
		Type:   AivenApplicationSucceeded,
		Status: v1.ConditionFalse,
	}
	localFailTrue = AivenApplicationCondition{
		Type:   AivenApplicationLocalFailure,
		Status: v1.ConditionTrue,
	}
	localFailFalse = AivenApplicationCondition{
		Type:   AivenApplicationLocalFailure,
		Status: v1.ConditionFalse,
	}
)

func TestAivenApplicationStatus_AddCondition(t *testing.T) {
	type args struct {
		preExisting  []AivenApplicationCondition
		newCondition AivenApplicationCondition
		wanted       []AivenApplicationCondition
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
			in.AddCondition(tt.args.newCondition)
			assert.Len(t, in.Conditions, len(tt.args.wanted))
			for i, c := range in.Conditions {
				assertConditionsMatch(tt.args.wanted[i], c)
			}
		})
	}
}
