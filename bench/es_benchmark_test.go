package main_test

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type RelationSuite struct {
	suite.Suite
}

func (s *RelationSuite) SetupTest() {
}

func BenchmarkStr(b *testing.B) {
	s := new(RelationSuite)
	s.SetT(&testing.T{})
	s.SetupTest()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

	}
}

// mysql 常用场合
func (s *RelationSuite) TestMySQL() {

}

func TestConvBench(t *testing.T) {
	suite.Run(t, new(RelationSuite))

}
