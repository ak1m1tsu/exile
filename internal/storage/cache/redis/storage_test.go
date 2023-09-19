package redis

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestNew(t *testing.T) {
	type args struct {
		client *redis.Client
	}
	tests := []struct {
		name    string
		args    args
		want    *Storage
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Set(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		value []byte
		ttl   time.Duration
	}
	tests := []struct {
		name    string
		s       *Storage
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Set(tt.args.ctx, tt.args.key, tt.args.value, tt.args.ttl); (err != nil) != tt.wantErr {
				t.Errorf("Storage.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		s       *Storage
		args    args
		want    []byte
		want1   bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.s.Get(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Storage.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestStorage_Del(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		s       *Storage
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Del(tt.args.ctx, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Storage.Del() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
