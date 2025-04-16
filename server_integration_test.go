//go:build integration
// +build integration

package registry

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// +------------+
// | Setup      |
// +------------+

type myLogger struct {
	DebugLogs []string
	InfoLogs  []string
	ErrorLogs []string
}

func (l *myLogger) Debug(format string, args ...any) {
	l.DebugLogs = append(l.DebugLogs, fmt.Sprintf(format, args...))
}
func (l *myLogger) Info(format string, args ...any) {
	l.InfoLogs = append(l.InfoLogs, fmt.Sprintf(format, args...))
}
func (l *myLogger) Error(format string, args ...any) {
	l.ErrorLogs = append(l.ErrorLogs, fmt.Sprintf(format, args...))
}
func (l *myLogger) Reset() {
	l.ErrorLogs = []string{}
	l.InfoLogs = []string{}
	l.DebugLogs = []string{}
}

type testsuite struct {
	suite.Suite
	Logger *myLogger
	Server *Server
}

// SetupTest is executed before each test
func (s *testsuite) SetupTest() {
	var err error
	s.Logger = &myLogger{}
	s.Server, err = NewInsecureServer("http://registry.local:5000", "", "", s.Logger)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), s.Server)
	s.Logger.Reset()
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(testsuite))
}

// +------------+
// | Tests      |
// +------------+

func (s *testsuite) TestPing() {
	err := s.Server.Ping()
	require.NoError(s.T(), err)
	assert.Equal(s.T(), s.Logger.DebugLogs, []string{
		"registry.ping url=http://registry.local:5000/v2/",
	})
}

func (s *testsuite) TestCatalog() {
	res, err := s.Server.GetCatalog()
	require.NoError(s.T(), err)
	assert.Equal(s.T(), res, []string{
		"debian",
		"jc21/dnsrouter",
		"jc21/rpmbuild-centos6",
	})
	assert.Equal(s.T(), s.Logger.DebugLogs, []string{
		"registry.catalog url=http://registry.local:5000/v2/_catalog",
	})
}

func (s *testsuite) TestTags() {
	tags, err := s.Server.GetTags("jc21/rpmbuild-centos6")
	require.NoError(s.T(), err)
	assert.Equal(s.T(), tags, []string{"latest"})
	assert.Equal(s.T(), s.Logger.DebugLogs, []string{
		"registry.tags url=http://registry.local:5000/v2/jc21/rpmbuild-centos6/tags/list repository=jc21/rpmbuild-centos6",
	})
}
