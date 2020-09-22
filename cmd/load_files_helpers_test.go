package cmd

import (
	"path/filepath"
	"testing"

	"github.com/damianoneill/h7t/dsl"
	"github.com/jarcoal/httpmock"
	"github.com/spf13/afero"
	"gopkg.in/resty.v1"
)

func TestPostHelperFileToResource(t *testing.T) {
	AppFs = afero.NewMemMapFs()
	matchFile := "something/anyfile.txt"

	_ = AppFs.Mkdir("something", 0777)
	_ = afero.WriteFile(AppFs, matchFile, []byte("junk"), 0777)

	// Get the underlying HTTP Client and set it to Mock
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://localhost:8080/api/v1/files/helper-files/anyfile.txt/",
		httpmock.NewStringResponder(200, ``))

	type args struct {
		rc           *resty.Client
		filename     string
		path         string
		paramName    string
		ci           dsl.ConnectionInfo
		shouldCommit bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid file upload",
			args: args{
				rc:        client,
				filename:  matchFile,
				paramName: "up_file",
				path:      "/api/v1/files/helper-files/" + filepath.Base(matchFile) + "/",
				ci: dsl.ConnectionInfo{
					Authority: "localhost:8080",
					Username:  "root",
					Password:  "changeme",
				},
				shouldCommit: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PostFileToResource(tt.args.rc,
				tt.args.filename,
				tt.args.path,
				tt.args.paramName,
				tt.args.ci,
				tt.args.shouldCommit); (err != nil) != tt.wantErr {
				t.Errorf("PostHelperFileToResource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
