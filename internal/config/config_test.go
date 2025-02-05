package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name       string
		envVars    map[string]string
		want       *Config
		wantErr    bool
		wantErrMsg string // Add a field for specific error message check
	}{
		{
			name: "Success",
			envVars: map[string]string{
				"S3_BUCKET":     "test-bucket",
				"SQS_QUEUE_URL": "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
			},
			want: &Config{
				S3:  S3{Bucket: "test-bucket"},
				SQS: SQS{QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"},
			},
			wantErr: false,
		},
		{
			name: "MissingS3Bucket",
			envVars: map[string]string{
				"SQS_QUEUE_URL": "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
			},
			want: &Config{
				SQS: SQS{QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"}, // S3 should be empty
			},
			wantErr: false, // env package doesn't return error if some vars are missing
		},
		{
			name: "MissingSQSQueueURL",
			envVars: map[string]string{
				"S3_BUCKET": "test-bucket",
			},
			want: &Config{
				S3: S3{Bucket: "test-bucket"}, // SQS should be empty
			},
			wantErr: false, // env package doesn't return error if some vars are missing
		},
		{
			name:    "EmptyEnvVars",
			envVars: map[string]string{},
			want:    &Config{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables for the test
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k) // Ensure cleanup after the test
			}

			got, err := Load()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrMsg != "" {
					assert.Contains(t, err.Error(), tt.wantErrMsg) // check specific error message
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// Example of a test for a function that uses the config
func TestSomethingThatUsesConfig(t *testing.T) {
	// Set necessary environment variables
	os.Setenv("S3_BUCKET", "test-bucket")
	os.Setenv("SQS_QUEUE_URL", "test-queue")
	defer func() {
		os.Unsetenv("S3_BUCKET")
		os.Unsetenv("SQS_QUEUE_URL")
	}()

	cfg, err := Load()
	assert.NoError(t, err)

	// Now you can use cfg in your assertions or test logic
	assert.Equal(t, "test-bucket", cfg.S3.Bucket)
	assert.Equal(t, "test-queue", cfg.SQS.QueueURL)

	// ... your test logic here ...
}
