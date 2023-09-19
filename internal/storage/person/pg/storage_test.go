package pg

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/romankravchuk/effective-mobile-test-task/internal/models"
)

func TestNew(t *testing.T) {
	type args struct {
		db *sql.DB
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
			got, err := New(tt.args.db)
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

func TestStorage_FindByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		s       *Storage
		args    args
		want    *models.Person
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.FindByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Update(t *testing.T) {
	type args struct {
		ctx context.Context
		p   *models.Person
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
			if err := tt.s.Update(tt.args.ctx, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("Storage.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		p   *models.Person
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
			if err := tt.s.Create(tt.args.ctx, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("Storage.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_List(t *testing.T) {
	type args struct {
		ctx    context.Context
		filter *models.Filter
	}
	tests := []struct {
		name    string
		s       *Storage
		args    args
		want    []models.Person
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.List(tt.args.ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
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
			if err := tt.s.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Storage.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
