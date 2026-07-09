package main

import (
	"bytes"
	"testing"

	pb "github.com/brotherlogic/devcontainer-manager/proto"
)

func TestFormatContainer(t *testing.T) {
	tests := []struct {
		name     string
		config   *pb.DevcontainerConfig
		expected string
	}{
		{
			name: "Ready state without issue number",
			config: &pb.DevcontainerConfig{
				Id: "my-container",
				State: pb.State_DCM_READY,
			},
			expected: "my-container  \n",
		},
		{
			name: "Creating state with issue number",
			config: &pb.DevcontainerConfig{
				Id: "123-container",
				Request: &pb.UpRequest{
					Identifier: &pb.Identifier{
						Id: &pb.Identifier_IssueNumber{
							IssueNumber: 123,
						},
					},
				},
				State: pb.State_DCM_CREATING,
			},
			expected: "123-container [123] creating\n",
		},
		{
			name: "Failed state with PR number",
			config: &pb.DevcontainerConfig{
				Id: "another-container",
				Request: &pb.UpRequest{
					Identifier: &pb.Identifier{
						Id: &pb.Identifier_PrNumber{
							PrNumber: 456,
						},
					},
				},
				State: pb.State_DCM_FAILED,
			},
			expected: "another-container  failed\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			formatContainer(&buf, tt.config)
			if got := buf.String(); got != tt.expected {
				t.Errorf("formatContainer() = %q, want %q", got, tt.expected)
			}
		})
	}
}
