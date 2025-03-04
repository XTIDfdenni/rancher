package git

import (
	"testing"

	assertlib "github.com/stretchr/testify/assert"
)

func Test_isGitSSH(t *testing.T) {
	testCases := []struct {
		gitURL   string
		expected bool
	}{
		// True cases
		{"customusername@github.com:user/repo.git", true},
		{"customusername@gitlab.com:user/repo.git", true},
		{"customusername@gitlab.com:user/repo", true},
		{"customusername@gitlab.com:user/repo-with-dashes.git", true},
		{"git@github.com:user/repo.git", true},
		{"git@gitlab.com:user/repo-with-dashes.git", true},
		{"git@gitlab.com:user/repo", true},
		// False cases
		{"https://github.com/user/repo.git", false},
		{"http://gitlab.com/user/repo.git", false},
		{"http://gitlab.com/user/repo", false},
		{"http://gitlab.com", false},
		{"git@gitlab.com", false},
	}
	assert := assertlib.New(t)
	for _, tc := range testCases {
		actual, err := isGitSSH(tc.gitURL)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		assert.Equalf(tc.expected, actual, "testcase: %v", tc)
	}
}

func Test_gitDir(t *testing.T) {
	assert := assertlib.New(t)
	testCases := []struct {
		namespace string
		name      string
		gitURL    string
		expected  string
	}{
		{
			"namespace", "name", "https://git.rancher.io/charts",
			"management-state/git-repo/namespace/name/4b40cac650031b74776e87c1a726b0484d0877c3ec137da0872547ff9b73a721",
		},
		// NOTE(manno): cannot test the other cases without poluting the filesystem
	}
	for _, tc := range testCases {
		actual := gitDir(tc.namespace, tc.name, tc.gitURL)
		assert.Equalf(tc.expected, actual, "testcase: %v", tc)
	}
}
