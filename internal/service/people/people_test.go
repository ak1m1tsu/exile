package people

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/romankravchuk/effective-mobile-test-task/internal/models"
	brokermocks "github.com/romankravchuk/effective-mobile-test-task/internal/storage/broker/mocks"
	cachemocks "github.com/romankravchuk/effective-mobile-test-task/internal/storage/cache/mocks"
	storagemocks "github.com/romankravchuk/effective-mobile-test-task/internal/storage/person/mocks"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "service with nil producer",
			args:    args{[]Option{WithProducer(nil, "")}},
			wantErr: true,
		},
		{
			name:    "service with producer",
			args:    args{[]Option{WithProducer(brokermocks.NewProducer(t), "")}},
			wantErr: false,
		},
		{
			name:    "service with nil storage",
			args:    args{[]Option{WithPersonStorage(nil)}},
			wantErr: true,
		},
		{
			name:    "service with storage",
			args:    args{[]Option{WithPersonStorage(storagemocks.NewStorage(t))}},
			wantErr: false,
		},
		{
			name:    "service with kafka producer with nil cfg map",
			args:    args{[]Option{WithKafkaProducer(nil, "")}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := New(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_Save(t *testing.T) {
	producer := brokermocks.NewProducer(t)

	svc, err := New(WithProducer(producer, ""))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	p := models.Person{}
	msg, _ := json.Marshal(&p)

	producer.On("Produce", msg).
		Once().
		Return(nil)

	err = svc.Save(context.Background(), p)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}
}

func TestService_Get(t *testing.T) {
	storage := storagemocks.NewStorage(t)
	cache := cachemocks.NewCache(t)

	svc, err := New(
		WithPersonStorage(storage),
		WithCache(cache, time.Minute),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	ctx := context.Background()
	uuid := "uuid"
	p := &models.Person{ID: uuid}
	data, _ := json.Marshal(p)

	cache.On("Get", ctx, uuid).
		Once().
		Return(nil, false, nil)

	storage.On("FindByID", ctx, uuid).
		Once().
		Return(p, nil)

	cache.On("Set", ctx, uuid, data, time.Minute).
		Once().
		Return(nil)

	_, err = svc.Get(ctx, uuid)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
}

func TestService_List(t *testing.T) {
	storage := storagemocks.NewStorage(t)
	cache := cachemocks.NewCache(t)

	svc, err := New(
		WithPersonStorage(storage),
		WithCache(cache, time.Minute),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	ctx := context.Background()
	filter := &models.Filter{Limit: 10, Skip: 0}
	query := filter.String()
	people := make([]models.Person, 0)
	data, _ := json.Marshal(people)

	cache.On("Get", ctx, query).Once().Return(nil, false, nil)
	storage.On("List", ctx, filter).Once().Return(people, nil)
	cache.On("Set", ctx, query, data, time.Minute).Once().Return(nil)

	_, err = svc.List(ctx, filter, query)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
}

func TestService_Update(t *testing.T) {
	storage := storagemocks.NewStorage(t)
	cache := cachemocks.NewCache(t)

	svc, err := New(
		WithPersonStorage(storage),
		WithCache(cache, time.Minute),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	ctx := context.Background()
	p := &models.Person{ID: "uuid", Name: "Ivan"}

	storage.On("Update", ctx, p).Once().Return(nil)
	cache.On("Del", ctx, p.ID).Once().Return(nil)

	err = svc.Update(ctx, p)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
}

func TestService_Delete(t *testing.T) {
	storage := storagemocks.NewStorage(t)
	cache := cachemocks.NewCache(t)

	svc, err := New(
		WithPersonStorage(storage),
		WithCache(cache, time.Minute),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	ctx := context.Background()
	uuid := "uuid"

	storage.On("Delete", ctx, uuid).Once().Return(nil)
	cache.On("Del", ctx, uuid).Once().Return(nil)

	if err = svc.Delete(ctx, uuid); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
}

func TestService_Close(t *testing.T) {
	producer := brokermocks.NewProducer(t)

	svc, err := New(WithProducer(producer, ""))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	producer.On("Close").Once().Return(nil)

	svc.Close()
}
