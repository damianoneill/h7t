package cmd

import (
	"bytes"
	"testing"
)

func Test_version(t *testing.T) {
	type args struct {
		bi buildInfo
	}
	tests := []struct {
		name  string
		args  args
		wantW string
	}{
		{"valid version output",
			args{
				buildInfo{
					"0.3.0",
					"c26cfaca0e38465935c48b13ae99d12fbf5d7cb1",
					"2019-11-05T21:16:06Z",
				},
			},
			"0.3.0, commit c26cfaca0e38465935c48b13ae99d12fbf5d7cb1, built at 2019-11-05T21:16:06Z \n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			version(w, tt.args.bi)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("version() = --%v--, want --%v--", gotW, tt.wantW)
			}
		})
	}
}
