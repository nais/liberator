package service

import (
	"context"
	"testing"

	"github.com/aiven/aiven-go-client/v2"
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
		svc *aiven.Service
		err error
	}
	tests := []struct {
		name                 string
		shortNameLookup      bool
		shortNameReturnValue returnValue
		longNameLookup       bool
		longNameReturnValue  returnValue
		want                 string
		wantErr              bool
	}{
		{
			name:            "no such service",
			shortNameLookup: true,
			shortNameReturnValue: returnValue{
				svc: nil,
				err: aiven.Error{Message: "not found", Status: 404},
			},
			longNameLookup: true,
			longNameReturnValue: returnValue{
				svc: nil,
				err: aiven.Error{Message: "not found", Status: 404},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:            "long named service",
			shortNameLookup: true,
			shortNameReturnValue: returnValue{
				svc: nil,
				err: aiven.Error{Message: "not found", Status: 404},
			},
			longNameLookup: true,
			longNameReturnValue: returnValue{
				svc: makeService(LongService),
				err: nil,
			},
			want:    LongService,
			wantErr: false,
		},
		{
			name:            "short named service",
			shortNameLookup: true,
			shortNameReturnValue: returnValue{
				svc: makeService(ShortService),
				err: nil,
			},
			longNameLookup: false,
			want:           ShortService,
			wantErr:        false,
		},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			times := 1
			if tt.wantErr {
				times = 2
			}
			mockInterface := NewMockInterface(t)
			if tt.shortNameLookup {
				mockInterface.
					On("Get", ctx, Project, ShortService).
					Times(times).
					Return(tt.shortNameReturnValue.svc, tt.shortNameReturnValue.err)
			}
			if tt.longNameLookup {
				mockInterface.
					On("Get", ctx, Project, LongService).
					Times(times).
					Return(tt.longNameReturnValue.svc, tt.longNameReturnValue.err)
			}

			r := NewCachedNameResolver(mockInterface)
			got, err := r.ResolveKafkaServiceName(ctx, Project)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveKafkaServiceName() first call error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ResolveKafkaServiceName() first call got = %v, want %v", got, tt.want)
			}
			got, err = r.ResolveKafkaServiceName(ctx, Project)
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
