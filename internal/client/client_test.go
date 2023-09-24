package client

import (
	"testing"

	"github.com/go-faker/faker/v4"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    APIError
		want string
	}{
		{
			name: "empty error",
			e:    APIError{},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("APIError.Error() = %v, want %v", got, tt.want)
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
		wantErr bool
	}{
		{
			name:    "empty name",
			wantErr: true,
		},
		{
			name: "invalid url",
			args: args{
				apiURL: faker.URL(),
				name:   "Ivan",
			},
			wantErr: true,
		},
		{
			name: "not success http request",
			args: args{
				apiURL: "https://jsonplaceholder.typicode.com/ghghgh",
				name:   "Ivan",
			},
			wantErr: true,
		},
		{
			name: "nationalize api request",
			args: args{
				apiURL: nationalizeURL,
				name:   "Ivan",
			},
			wantErr: false,
		},
		{
			name: "genderize api request",
			args: args{
				apiURL: genderizeURL,
				name:   "Ivan",
			},
			wantErr: false,
		},
		{
			name: "agify api request",
			args: args{
				apiURL: agifyURL,
				name:   "Ivan",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := get(tt.args.apiURL, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
