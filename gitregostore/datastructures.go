package gitregostore

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"

	// "github.com/armosec/capacketsgo/opapolicy"
	"github.com/armosec/armoapi-go/armotypes"
	opapolicy "github.com/kubescape/opa-utils/reporthandling"
	"github.com/kubescape/opa-utils/reporthandling/attacktrack/v1alpha1"

	"github.com/go-gota/gota/dataframe"
)

type GitRegoStore struct {
	frameworksLock                     sync.RWMutex
	DefaultConfigInputsLock            sync.RWMutex
	rulesLock                          sync.RWMutex
	controlsLock                       sync.RWMutex
	attackTracksLock                   sync.RWMutex
	systemPostureExceptionPoliciesLock sync.RWMutex
	ControlRuleRelations               dataframe.DataFrame
	FrameworkControlRelations          dataframe.DataFrame
	httpClient                         *http.Client
	Tag                                string
	Owner                              string
	CurGitVersion                      string
	Branch                             string
	URL                                string
	Path                               string
	BaseUrl                            string
	Repository                         string
	DefaultConfigInputs                armotypes.CustomerConfig
	AttackTracks                       []v1alpha1.AttackTrack
	Frameworks                         []opapolicy.Framework
	Controls                           []opapolicy.Control
	Rules                              []opapolicy.PolicyRule
	SystemPostureExceptionPolicies     []armotypes.PostureExceptionPolicy
	FrequencyPullFromGitMinutes        int
	Watch                              bool
	StripFilesExtension                bool
}

func newGitRegoStore(baseUrl string, owner string, repository string, path string, tag string, branch string, frequency int) *GitRegoStore {
	var stripFilesExtension bool

	watch := false
	if frequency > 0 {
		watch = true
	}

	if strings.Contains(tag, "latest") || strings.Contains(tag, "download") {
		// TODO - This condition was added to avoid dependency on updating productions configs on deployment.
		// Once production configs are updated (branch set to ""), this condition can be removed.
		if strings.ToLower(branch) == "master" {
			branch = ""
		}
		stripFilesExtension = true
	} else {
		stripFilesExtension = false
	}

	gs := &GitRegoStore{httpClient: &http.Client{},
		BaseUrl:                     baseUrl,
		Owner:                       owner,
		Repository:                  repository,
		Path:                        path,
		Tag:                         tag,
		Branch:                      branch,
		FrequencyPullFromGitMinutes: frequency,
		Watch:                       watch,
	}

	gs.StripFilesExtension = stripFilesExtension

	return gs
}

// NewGitRegoStore return gitregostore obj with basic fields, before pulling from git
func NewGitRegoStore(baseUrl string, owner string, repository string, path string, tag string, branch string, frequency int) *GitRegoStore {

	gs := newGitRegoStore(baseUrl, owner, repository, path, tag, branch, frequency)
	gs.setURL()

	return gs
}

// SetRegoObjects pulls opa obj from git and stores in gitregostore
func (gs *GitRegoStore) SetRegoObjects() error {
	err := gs.setObjects()
	return err
}

// NewDefaultGitRegoStore - generates git store object for production regolibrary release files.
// Release files source: "https://github.com/kubescape/regolibrary/releases/latest/download"
func NewDefaultGitRegoStore(frequency int) *GitRegoStore {
	gs := NewGitRegoStore("https://github.com", "kubescape", "regolibrary", "releases", "latest/download", "", frequency)
	return gs
}

// NewDevGitRegoStore - generates git store object for dev regolibrary release files
// Release files source: "https://raw.githubusercontent.com/kubescape/regolibrary/dev/releaseDev"
func NewDevGitRegoStore(frequency int) *GitRegoStore {
	gs := NewGitRegoStore("https://raw.githubusercontent.com", "kubescape", "regolibrary", "releaseDev", "", "dev", frequency)
	return gs
}

// Deprecated
// if frequency < 0 will pull only once
func InitGitRegoStore(baseUrl string, owner string, repository string, path string, tag string, branch string, frequency int) *GitRegoStore {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("InitGitRegoStore failed: stacktrace from panic: \n" + string(debug.Stack()))
		}
	}()
	gs := newGitRegoStore(baseUrl, owner, repository, path, tag, branch, frequency)
	gs.setURL()
	gs.setObjects()
	return gs
}

// Deprecated
func InitDefaultGitRegoStore(frequency int) *GitRegoStore {
	return InitGitRegoStore("https://github.com", "kubescape", "regolibrary", "releases", "latest/download", "", frequency)
}
