package sl

import (
	"log/slog"
	"reflect"
	"testing"
)

func TestErr(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want slog.Attr
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Err(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Err() = %v, want %v", got, tt.want)
			}
		})
	}
}
