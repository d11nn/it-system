package processor

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/Alonza0314/it-system/controller/backend/constant"
	"github.com/Alonza0314/it-system/controller/backend/model"
)

func (p *Processor) GetGithubPRs(nf string) (*model.ResponseGetGithubPRs, *model.ErrorDetail) {
	prs, err := p.itContext.GetPrList(nf)
	if err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("Failed to get PR list: %v", err),
		}
	}
	p.GitLog.Debugf("Retrieved %d PRs", len(prs))
	p.GitLog.Tracef("PRs details: %+v", prs)

	response := &model.ResponseGetGithubPRs{
		Message: "PRs retrieved successfully",
		PRs:     make([]model.PR, len(prs)),
	}

	for i, pr := range prs {
		response.PRs[i] = model.PR{
			Number: pr.Number(),
			Title:  pr.Title(),
		}
	}

	return response, nil
}

func (p *Processor) SuggestLibraryPRs(req *model.RequestDependencySuggestions) (*model.ResponseDependencySuggestions, *model.ErrorDetail) {
	suggestions := make([]model.LibraryPrSuggestion, 0)
	seen := make(map[string]bool)

	for _, nfPr := range req.NFPrList {
		repo := normalizeGithubRepoName(nfPr.NfName)
		if repo == "" || slices.Contains(constant.LIBRARY_LIST, repo) {
			continue
		}

		detail, err := p.itContext.GetPrDetail(repo, nfPr.PR)
		if err != nil {
			return nil, &model.ErrorDetail{
				HttpStatus: http.StatusInternalServerError,
				Detail:     fmt.Sprintf("Failed to inspect %s PR #%d: %v", repo, nfPr.PR, err),
			}
		}

		text := strings.ToLower(detail.Title() + "\n" + detail.Body())
		mentionedLibraries := mentionedLibraryRepos(text)
		if len(mentionedLibraries) == 0 {
			continue
		}

		for _, library := range mentionedLibraries {
			candidates, err := p.itContext.GetPrList(library)
			if err != nil {
				return nil, &model.ErrorDetail{
					HttpStatus: http.StatusInternalServerError,
					Detail:     fmt.Sprintf("Failed to get %s PR list: %v", library, err),
				}
			}

			for _, candidate := range candidates {
				if !isLikelyLibraryDependencyPR(detail.Author(), candidate.Author(), detail.Title(), text, candidate.Title()) {
					continue
				}

				key := fmt.Sprintf("%s#%d", library, candidate.Number())
				if seen[key] {
					continue
				}
				seen[key] = true
				suggestions = append(suggestions, model.LibraryPrSuggestion{
					RepoName:   library,
					PR:         candidate.Number(),
					Title:      candidate.Title(),
					Reason:     fmt.Sprintf("%s PR #%d mentions %s and this open %s PR has a related title/author signal", repo, nfPr.PR, library, library),
					Confidence: "medium",
				})
			}
		}
	}

	return &model.ResponseDependencySuggestions{
		Message:     "Dependency suggestions generated",
		Suggestions: suggestions,
	}, nil
}

func normalizeGithubRepoName(name string) string {
	normalized := strings.ToLower(strings.TrimSpace(name))
	if normalized == constant.UPF {
		return constant.GO_UPF
	}
	if slices.Contains(constant.NF_LIST, normalized) || slices.Contains(constant.LIBRARY_LIST, normalized) || normalized == constant.GO_UPF {
		return normalized
	}
	return ""
}

func mentionedLibraryRepos(text string) []string {
	mentioned := make([]string, 0, len(constant.LIBRARY_LIST))
	for _, library := range constant.LIBRARY_LIST {
		if strings.Contains(text, library) || strings.Contains(text, "github.com/free5gc/"+library) {
			mentioned = append(mentioned, library)
		}
	}
	return mentioned
}

func isLikelyLibraryDependencyPR(nfAuthor, candidateAuthor, nfTitle, nfText, candidateTitle string) bool {
	candidate := strings.ToLower(candidateTitle)
	if nfAuthor != "" && strings.EqualFold(nfAuthor, candidateAuthor) {
		return true
	}

	nfTitleWords := significantWords(nfTitle)
	for _, word := range nfTitleWords {
		if strings.Contains(candidate, word) {
			return true
		}
	}

	for _, word := range significantWords(nfText) {
		if strings.Contains(candidate, word) {
			return true
		}
	}

	return false
}

func significantWords(text string) []string {
	fields := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return (r < 'a' || r > 'z') && (r < '0' || r > '9')
	})

	words := make([]string, 0, len(fields))
	for _, field := range fields {
		if len(field) < 4 {
			continue
		}
		if slices.Contains([]string{"github", "free5gc", "with", "from", "that", "this", "update", "fix", "feat"}, field) {
			continue
		}
		words = append(words, field)
		if len(words) >= 8 {
			break
		}
	}

	return words
}
