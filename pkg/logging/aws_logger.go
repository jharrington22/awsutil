// This file contains an AWS logger that uses the logging framework of the project.

package logging

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/sirupsen/logrus"
)

// AWSLoggerBuilder contains the information and logic needed to create an AWS logger that uses
// the logging framework of the project. Don't create instances of this type directly; use the
// NewAWSLogger function instead.
type AWSLoggerBuilder struct {
	logger *logrus.Logger
}

// AWSLogger is an implementation of the OCM logger interface that uses the logging framework of
// the project. Don't create instances of this type directly; use the NewAWSLogger function instead.
type AWSLogger struct {
	logger *logrus.Logger
}

// Make sure that we implement the OCM logger interface.
var _ aws.Logger = &AWSLogger{}

// NewAWSLogger creates new builder that can then be used to configure and build an OCM logger that
// uses the logging framework of the project.
func NewAWSLogger() *AWSLoggerBuilder {
	return &AWSLoggerBuilder{}
}

// Logger sets the underlying logger that will be used by the OCM logger to send the messages to the
// log.
func (b *AWSLoggerBuilder) Logger(value *logrus.Logger) *AWSLoggerBuilder {
	b.logger = value
	return b
}

// Build uses the information stored in the builder to create a new OCM logger that uses the logging
// framework of the project.
func (b *AWSLoggerBuilder) Build() (result *AWSLogger, err error) {
	// Check parameters:
	if b.logger == nil {
		err = fmt.Errorf("logger is mandatory")
		return
	}

	// Create and populate the object:
	result = &AWSLogger{
		logger: b.logger,
	}

	return
}

func (l *AWSLogger) Log(args ...interface{}) {
	l.logger.Info(args...)
}
