package cmd

import (
	"reflect"
	"testing"

	"github.com/spf13/afero"
)

const matchFile = "something/anyfile.txt"

func Test_getDirectoryContents(t *testing.T) {
	AppFs = afero.NewMemMapFs()

	_ = AppFs.Mkdir("nothing", 0777)
	_ = AppFs.Mkdir("something", 0777)
	_ = afero.WriteFile(AppFs, matchFile, []byte("junk"), 0777)

	type args struct {
		inputDirectory string
	}
	tests := []struct {
		name        string
		args        args
		wantMatches []string
		wantErr     bool
	}{
		{
			name: "no files, return nil",
			args: args{
				inputDirectory: "nothing",
			},
			wantMatches: nil,
			wantErr:     false,
		},
		{
			name: "sample file, return 1",
			args: args{
				inputDirectory: "something",
			},
			wantMatches: []string{matchFile},
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatches, err := getDirectoryContents(tt.args.inputDirectory)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDirectoryContents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMatches, tt.wantMatches) {
				t.Errorf("getDirectoryContents() = %v, want %v", gotMatches, tt.wantMatches)
			}
		})
	}
}
