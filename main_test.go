package main

import (
	"os"
	"path"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
)

func TestBucketCompany(t *testing.T) {
	awsRegion := "us-east-2"
	cwd, _ := os.Getwd()
	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Quick:       true,
		SkipRefresh: true,
		Dir:         path.Join(cwd),
		Config: map[string]string{
			"aws:region": awsRegion,
		},
	})
}
