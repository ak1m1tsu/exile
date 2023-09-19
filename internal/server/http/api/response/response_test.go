package response

import (
	"reflect"
	"testing"
)

func TestOK(t *testing.T) {
	tests := []struct {
		name string
		want Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OK(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OK() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError(t *testing.T) {
	type args struct {
		err string
	}
	tests := []struct {
		name string
		args args
		want Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Error(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
