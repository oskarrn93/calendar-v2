package main

import (
	"testing"

	"github.com/aws/aws-cdk-go/awscdk"
)

func TestInfraStack(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewInfraStack(app, "MyStack", nil)

	// THEN
	if stack == nil {
		t.Fail()
	}

}
