package index

import (
	"log"
	"sort"
	"strings"
)

var (
	searches = make(chan *search, 10)
	updates  = make(chan []*Summary)
)

//Summary contains the key summary information for a Swagger file
type Summary struct {
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Version     string   `json:"version"`
}

//SearchResult contains a search result
type SearchResult struct {
	Relevance int      `json:"relevance"`
	Summary   *Summary `json:"summary"`
}

type search struct {
	term    string
	results chan *SearchResult
}

func init() {
	go func() {

		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic in index io: %s", err)
			}
		}()

		index := []*Summary{}

		for {
			select {
			case update := <-updates:
				index = update

				log.Printf("index of swagger summaries updated, %v currently available", len(index))

			case search := <-searches:
				term := strings.ToLower(search.term)

				if term == "" {
					for _, summary := range index {
						search.results <- &SearchResult{Relevance: 0, Summary: summary}
					}
				} else {
					for _, summary := range index {
						relevance := 0
						title := strings.ToLower(summary.Title)

						switch {
						case title == term:
							relevance = 1
						case strings.Contains(title, term):
							relevance = 2
						case strings.Contains(strings.ToLower(strings.Join(summary.Tags, ",")), term):
							relevance = 3
						case strings.Contains(strings.ToLower(summary.Description), term):
							relevance = 4
						}

						if relevance > 0 {
							search.results <- &SearchResult{Relevance: relevance, Summary: summary}
						}
					}
				}

				close(search.results)
			}
		}
	}()
}

//Search returns all summaries that match the specified search term
func Search(term string) []*SearchResult {
	c := make(chan *SearchResult, 10)
	searches <- &search{term: term, results: c}

	results := []*SearchResult{}

	for r := range c {
		results = append(results, r)
	}

	if term == "" {
		sort.Slice(results, func(i, j int) bool {
			return strings.Compare(results[i].Summary.Title, results[j].Summary.Title) < 0
		})
	} else {
		sort.Slice(results, func(i, j int) bool {
			return results[i].Relevance < results[j].Relevance
		})
	}

	return results
}
