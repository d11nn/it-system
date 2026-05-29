package context

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Alonza0314/it-system/controller/backend/constant"

	"github.com/free-ran-ue/util"
)

type pr struct {
	Num       int       `json:"number"`
	Tit       string    `json:"title"`
	BodyText  string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      struct {
		Login string `json:"login"`
	} `json:"user"`
	Head struct {
		Sha  string `json:"sha"`
		Repo struct {
			FullName string `json:"full_name"`
		} `json:"repo"`
	} `json:"head"`
}

func (p *pr) Number() int {
	return p.Num
}

func (p *pr) Title() string {
	return p.Tit
}

func (p *pr) Body() string {
	return p.BodyText
}

func (p *pr) Author() string {
	return p.User.Login
}

func (p *pr) HeadSHA() string {
	return p.Head.Sha
}

func (p *pr) HeadRepoFullName() string {
	return p.Head.Repo.FullName
}

type githubContext struct{}

func newGithubContext() *githubContext {
	return &githubContext{}
}

func (ctx *githubContext) getPrList(nf string) ([]pr, error) {
	apiUrl := fmt.Sprintf(constant.GITHUB_FREE5GC_BASE_API_URL, nf)
	responseRaw, err := util.SendHttpRequest(apiUrl, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}

	if responseRaw.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get PRs from %s: status code %d", apiUrl, responseRaw.StatusCode)
	}

	var prList []pr
	if err := json.Unmarshal(responseRaw.Body, &prList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal PR response from %s: %v", apiUrl, err)
	}

	return prList, nil
}

func (ctx *githubContext) getPrDetail(repo string, prNumber int) (*pr, error) {
	apiUrl := fmt.Sprintf("%s/%d", fmt.Sprintf(constant.GITHUB_FREE5GC_BASE_API_URL, repo), prNumber)
	responseRaw, err := util.SendHttpRequest(apiUrl, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}

	if responseRaw.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get PR from %s: status code %d", apiUrl, responseRaw.StatusCode)
	}

	var detail pr
	if err := json.Unmarshal(responseRaw.Body, &detail); err != nil {
		return nil, fmt.Errorf("failed to unmarshal PR response from %s: %v", apiUrl, err)
	}

	return &detail, nil
}
