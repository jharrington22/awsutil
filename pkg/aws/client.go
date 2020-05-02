package aws

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/iam/stsiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"gitlab.cee.redhat.com/service/moactl/pkg/logging"
)

type Client interface {
	GetRegion() string
	ValidateCredentials() (bool, error)
	EnsureOsdCcsAdminUser() (bool, error)
	CreateAccessKey(username string) (*AWSAccessKey, error)
	GetCreator() (*AWSCreator, error)
	TagUser(username string, clusterID string, clusterName string) error
	ValidateSCP() (bool, error)
}

// ClientBuilder contains the information and logic needed to build a new AWS client.
type ClientBuilder struct {
	logger *logrus.Logger
}

type awsClient struct {
	logger     *logrus.Logger
	iamClient  iamiface.IAMAPI
	stsClient  stsiface.STSAPI
	awsSession *session.Session
}

// NewClient creates a builder that can then be used to configure and build a new AWS client.
func NewClient() *ClientBuilder {
	return &ClientBuilder{}
}

// Logger sets the logger that the AWS client will use to send messages to the log.
func (b *ClientBuilder) Logger(value *logrus.Logger) *ClientBuilder {
	b.logger = value
	return b
}

// Build uses the information stored in the builder to build a new AWS client.
func (b *ClientBuilder) Build() (result Client, err error) {
	// Check parameters:
	if b.logger == nil {
		err = fmt.Errorf("logger is mandatory")
		return
	}

	// Create the AWS logger:
	logger, err := logging.NewAWSLogger().
		Logger(b.logger).
		Build()
	if err != nil {
		return
	}

	// Create the AWS session:
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	sess.Config.Logger = logger
	sess.Config.HTTPClient = &http.Client{
		Transport: http.DefaultTransport,
	}
	if b.logger.IsLevelEnabled(logrus.DebugLevel) {
		var dumper http.RoundTripper
		dumper, err = logging.NewRoundTripper().
			Logger(b.logger).
			Next(sess.Config.HTTPClient.Transport).
			Build()
		if err != nil {
			return
		}
		sess.Config.HTTPClient.Transport = dumper
	}
	if err != nil {
		return
	}

	// Check that the region is set:
	region := aws.StringValue(sess.Config.Region)
	if region == "" {
		err = fmt.Errorf("region is not set")
		return
	}

	// Check that the AWS credentials are available:
	_, err = sess.Config.Credentials.Get()
	if err != nil {
		err = fmt.Errorf("can't find credentials: %v", err)
		return
	}

	// Create and populate the object:
	result = &awsClient{
		logger:     b.logger,
		iamClient:  iam.New(sess),
		stsClient:  sts.New(sess),
		awsSession: sess,
	}

	return
}
