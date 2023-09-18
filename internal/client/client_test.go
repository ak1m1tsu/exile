package client

import (
	"reflect"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    APIError
		want string
	}{
		{
			name: "base test",
			e:    APIError{Message: "some error"},
			want: "some error",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("APIError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFetchAge(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Good input",
			args:    args{"Ivan"},
			wantErr: false,
		},
		{
			name:    "Bad input",
			args:    args{"1234567890"},
			wantErr: true,
		},
		{
			name:    "Empty input",
			args:    args{""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			age, err := FetchAge(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchAge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(age)
		})
	}
}

func TestFetchGender(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Good input",
			args:    args{"Ivan"},
			wantErr: false,
		},
		{
			name:    "Bad input",
			args:    args{"1234567890"},
			wantErr: true,
		},
		{
			name:    "Empty input",
			args:    args{""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := FetchGender(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchGender() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFetchNationality(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Good input",
			args:    args{"Ivan"},
			wantErr: false,
		},
		{
			name:    "Bad input",
			args:    args{"1234567890"},
			wantErr: true,
		},
		{
			name:    "Empty input",
			args:    args{""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := FetchNationality(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchNationality() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_get(t *testing.T) {
	type args struct {
		apiURL string
		name   string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "empty name",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := get(tt.args.apiURL, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("get() = %v, want %v", got, tt.want)
			}
		})
	}
}
