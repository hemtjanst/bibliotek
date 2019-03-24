package mqtt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTopicTest(t *testing.T) {

	assert.Equal(t, true, TopicTest("asdf/asdf/asdf", "#"))
	assert.Equal(t, true, TopicTest("foo/asdf/asdf", "foo/#"))
	assert.Equal(t, false, TopicTest("bar/asdf/asdf", "foo/#"))
	assert.Equal(t, true, TopicTest("a/b/c", "a/b/c"))
	assert.Equal(t, false, TopicTest("a/b/c", "a/b/c/d"))
	assert.Equal(t, false, TopicTest("a/b/c/d", "a/b/c"))
	assert.Equal(t, true, TopicTest("a/b/c/d", "a/+/+/d"))

}
