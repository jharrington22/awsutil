// This file contains the code to build the default loggers used by the project.

package logging

import (
	"github.com/sirupsen/logrus"

	"github.com/jharrington22/awsutil/pkg/debug"
)

// LoggerBuilder contains the information and logic needed to create the default loggers used by
// the project. Don't create instances of this type directly; use the NewLogger function instead.
type LoggerBuilder struct {
}

// NewLogger creates new builder that can then be used to configure and build an OCM logger that
// uses the logging framework of the project.
func NewLogger() *LoggerBuilder {
	return &LoggerBuilder{}
}

// Build uses the information stored in the builder to create a new logger.
func (b *LoggerBuilder) Build() (result *logrus.Logger, err error) {
	// Create the logger:
	result = logrus.New()
	result.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	// Enable the debug level if needed:
	if debug.Enabled() {
		result.SetLevel(logrus.DebugLevel)
	}

	return
}
