package service

import (
	"testing"

	"github.com/aiven/aiven-go-client"
)

const (
	Project        = "project"
	ShortService   = "kafka"
	LongService    = "project-kafka"
	NoMatchService = "kafka-project"
)

func makeService(name string) *aiven.Service {
	return &aiven.Service{
		Name: name,
	}
}

func TestCachedNameResolver_ResolveKafkaServiceName(t *testing.T) {
	type returnValue struct {
		svcs []*aiven.Service
		err  error
	}
	tests := []struct {
		name          string
		timesResolved int
		returnValue   returnValue
		want          string
		wantErr       bool
	}{
		{
			name:          "no such service",
			timesResolved: 2,
			returnValue: returnValue{
				svcs: make([]*aiven.Service, 0),
				err:  nil,
			},
			want:    "",
			wantErr: true,
		},
		{
			name:          "long named service",
			timesResolved: 1,
			returnValue: returnValue{
				svcs: []*aiven.Service{makeService(LongService)},
				err:  nil,
			},
			want:    LongService,
			wantErr: false,
		},
		{
			name:          "short named service",
			timesResolved: 1,
			returnValue: returnValue{
				svcs: []*aiven.Service{makeService(ShortService)},
				err:  nil,
			},
			want:    ShortService,
			wantErr: false,
		},
		{
			name:          "unwanted kafka service",
			timesResolved: 2,
			returnValue: returnValue{
				svcs: []*aiven.Service{makeService(NoMatchService)},
				err:  nil,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockInterface := NewMockInterface(t)
			mockInterface.
				On("List", Project).
				Times(tt.timesResolved).
				Return(tt.returnValue.svcs, tt.returnValue.err)

			r := NewCachedNameResolver(mockInterface)
			got, err := r.ResolveKafkaServiceName(Project)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveKafkaServiceName() first call error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ResolveKafkaServiceName() first call got = %v, want %v", got, tt.want)
			}
			got, err = r.ResolveKafkaServiceName(Project)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveKafkaServiceName() second call error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ResolveKafkaServiceName() second call got = %v, want %v", got, tt.want)
			}
			mockInterface.AssertExpectations(t)
		})
	}
}
