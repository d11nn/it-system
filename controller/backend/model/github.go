package model

type ResponseGetGithubPRs struct {
	Message string `json:"message" binding:"required"`
	PRs     []PR   `json:"prs,omitempty"`
}

type RequestDependencySuggestions struct {
	NFPrList []NfPr `json:"nfPrList" binding:"required"`
}

type ResponseDependencySuggestions struct {
	Message     string                `json:"message" binding:"required"`
	Suggestions []LibraryPrSuggestion `json:"suggestions,omitempty"`
}

type LibraryPrSuggestion struct {
	RepoName   string `json:"repoName" binding:"required"`
	PR         int    `json:"pr" binding:"required"`
	Title      string `json:"title,omitempty"`
	Reason     string `json:"reason,omitempty"`
	Confidence string `json:"confidence,omitempty"`
}

type PR struct {
	Number int    `json:"number" binding:"required"`
	Title  string `json:"title" binding:"required"`
}
