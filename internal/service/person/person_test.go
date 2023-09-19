package person

import (
	"reflect"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/romankravchuk/effective-mobile-test-task/internal/storage/person/pg"
)

func TestNew(t *testing.T) {
	type args struct {
		options []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "with nil people storage",
			args:    args{options: []Option{WithPeopleStorage(nil)}},
			wantErr: true,
		},
		{
			name:    "with nil consumer",
			args:    args{options: []Option{WithConsumer(nil)}},
			wantErr: true,
		},
		{
			name:    "with nil producer",
			args:    args{options: []Option{WithProducer(nil, "")}},
			wantErr: true,
		},
		{
			name:    "with not nil people storage",
			args:    args{options: []Option{WithPeopleStorage(&pg.Storage{})}},
			wantErr: false,
		},
		{
			name:    "with not nil consumer",
			args:    args{options: []Option{WithConsumer(&kafka.Consumer{})}},
			wantErr: false,
		},
		{
			name:    "with not nil producer",
			args:    args{options: []Option{WithProducer(&kafka.Producer{}, "")}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_Save(t *testing.T) {
	tests := []struct {
		name    string
		s       *Service
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Save()
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.Save() = %v, want %v", got, tt.want)
			}
		})
	}
}
