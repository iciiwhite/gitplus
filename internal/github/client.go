package github

import (
	"context"
	"net/http"

	"github.com/google/go-github/v56/github"
)

type Client struct {
	gh *github.Client
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{
		gh: github.NewClient(httpClient),
	}
}

func (c *Client) ListRepos(ctx context.Context, user string) ([]*github.Repository, error) {
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	var allRepos []*github.Repository
	for {
		repos, resp, err := c.gh.Repositories.List(ctx, user, opt)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allRepos, nil
}

func (c *Client) CreatePullRequest(ctx context.Context, owner, repo, title, head, base, body string) (*github.PullRequest, error) {
	newPR := &github.NewPullRequest{
		Title: &title,
		Head:  &head,
		Base:  &base,
		Body:  &body,
	}
	pr, _, err := c.gh.PullRequests.Create(ctx, owner, repo, newPR)
	return pr, err
}

func (c *Client) GetIssues(ctx context.Context, owner, repo string) ([]*github.Issue, error) {
	opt := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	var allIssues []*github.Issue
	for {
		issues, resp, err := c.gh.Issues.ListByRepo(ctx, owner, repo, opt)
		if err != nil {
			return nil, err
		}
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allIssues, nil
}