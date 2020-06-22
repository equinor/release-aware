package main

import (
	"errors"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/levigross/grequests"
	"github.com/tidwall/gjson"
)

type Release struct {
	Name           string    `json:"name"`
	TagName        string    `json:"tag_name"`
	PublishedAt    time.Time `json:"published_at"`
	Severity       string    `json:"severity"`
	RepositoryName string    `json:"repository_name"`
	HtmlUrl        string    `json:"html_url"`
	Days           int       `json:"days"`
	Type           string    `json:"type"`
}

type HelmRelease struct {
	Name           string    `json:"name"`
	TagName        string    `json:"tag_name"`
	PublishedAt    time.Time `json:"published_at"`
	Severity       string    `json:"severity"`
	RepositoryName string    `json:"repository_name"`
	HtmlUrl        string    `json:"html_url"`
	Days           int       `json:"days"`
	Type           string    `json:"type"`
}

// DateSorter sorts releases by date.
type DateSorter []Release

func (a DateSorter) Len() int {
	return len(a)
}
func (a DateSorter) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a DateSorter) Less(i, j int) bool {
	return a[i].PublishedAt.After(a[j].PublishedAt)
}

// Don't know how to overload in Go, so duplicating..
type HelmDateSorter []HelmRelease

func (a HelmDateSorter) Len() int {
	return len(a)
}
func (a HelmDateSorter) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a HelmDateSorter) Less(i, j int) bool {
	return a[i].PublishedAt.After(a[j].PublishedAt)
}

func getSeverity(days int) string {
	if days < 3 {
		return "critical"
	} else if days < 7 {
		return "warning"
	} else {
		return "white"
	}
}

func getGithubRepositories() []string {
	repositories := strings.ReplaceAll(os.Getenv("REPOSITORIES"), "\n", "")
	return strings.Split(strings.ReplaceAll(repositories, " ", ""), ",")
}

func getGithubToken() string {
	return "token " + os.Getenv("GITHUB_TOKEN")
}

func parseRelease(resp string, prefix string) Release {
	var release Release
	isDraft := gjson.Get(resp, prefix+".isDraft")
	isPrerelease := gjson.Get(resp, prefix+".isPrerelease")
	releaseType := "Latest release"
	if isPrerelease.Bool() {
		releaseType = "Pre-release"
	}
	if isDraft.Bool() {
		releaseType = "Draft"
	}
	release.Type = releaseType
	release.HtmlUrl = gjson.Get(resp, prefix+".url").String()
	release.PublishedAt = gjson.Get(resp, prefix+".publishedAt").Time()
	release.TagName = gjson.Get(resp, prefix+".tagName").String()
	release.Name = gjson.Get(resp, prefix+".name").String()
	return release
}

func parseTag(resp string, repositoryName string, prefix string) Release {
	var release Release
	tagName := gjson.Get(resp, prefix+".name")

	// Need to select a date, not sure if this is completely correct.
	pushedDate := gjson.Get(resp, prefix+".target.pushedDate")
	date := gjson.Get(resp, prefix+".target.tagger.date")
	authoredDate := gjson.Get(resp, prefix+".target.authoredDate")
	if pushedDate.Value() != nil {
		release.PublishedAt = pushedDate.Time()
	} else if date.Value() != nil {
		release.PublishedAt = date.Time()
	} else if authoredDate.Value() != nil {
		release.PublishedAt = authoredDate.Time()
	}
	release.TagName = tagName.String()
	release.HtmlUrl = "https://github.com/" + repositoryName + "/releases/tag/" + tagName.String()
	release.Type = "Tag"
	return release
}

func getGithubRelease(repositoryName string) (Release, error) {
	splitted := strings.Split(repositoryName, "/")
	query := `{
		repository(owner: "` + splitted[0] + `", name: "` + splitted[1] + `") {
		  releases (last: 1, orderBy: {field: CREATED_AT, direction: ASC}) {
		     edges {
		  		node {
			  	  createdAt
                  publishedAt
                  url
	    		  isDraft
  		          isPrerelease
                  tagName
                  name
                }
		     }
		  }
		}
		repository(owner: "` + splitted[0] + `", name: "` + splitted[1] + `") {
			refs(refPrefix: "refs/tags/", last: 1, orderBy: {field: TAG_COMMIT_DATE, direction: ASC}) {
			edges {
				node {
					name
					target {
					oid
					... on Commit {
						pushedDate
						authoredDate

					}
					... on Tag {
						tagger {
						date
						}
                     }
					}
				}
				}
			}
		}
	}`
	url := "https://api.github.com/graphql"
	requestOptions := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization": getGithubToken()},
		JSON:    map[string]string{"query": query},
	}
	resp, err := grequests.Post(url, requestOptions)
	if err != nil {
		return Release{}, errors.New(err.Error())
	}
	if !resp.Ok {
		return Release{}, errors.New(string(resp.Bytes()))
	}

	prefixRelease := "data.repository.releases.edges.0.node"
	prefixTag := "data.repository.refs.edges.0.node"

	isRelease := gjson.Get(resp.String(), prefixRelease).Exists()
	isTag := gjson.Get(resp.String(), prefixTag).Exists()

	if isRelease && isTag {
		release := parseRelease(resp.String(), prefixRelease)
		tag := parseTag(resp.String(), repositoryName, prefixTag)
		// Need to select the latest release
		if release.PublishedAt.After(tag.PublishedAt) {
			return release, nil
		} else {
			return tag, nil
		}
	} else if isRelease {
		return parseRelease(resp.String(), prefixRelease), nil
	} else if isTag {
		return parseTag(resp.String(), repositoryName, prefixTag), nil
	}
	return Release{}, errors.New("Could not find any release for " + repositoryName)
}

func getHelmhubRepositories() []string {
	repositories := strings.ReplaceAll(os.Getenv("HELM_REPOS"), "\n", "")
	return strings.Split(strings.ReplaceAll(repositories, " ", ""), ",")
}

func getHelmhubRelease(repositoryName string) (HelmRelease, error) {
	//	curl https://hub.helm.sh/api/chartsvc/v1/charts/stable/sealed-secrets | jq '.data.relationships.latestChartVersion.data.version'

	baseUrl := "https://hub.helm.sh/api/chartsvc/v1/charts/"
	url := (baseUrl + repositoryName)

	requestOptions := &grequests.RequestOptions{}

	resp, err := grequests.Get(url, requestOptions)
	if err != nil {
		return HelmRelease{}, errors.New(err.Error())
	}
	if !resp.Ok {
		return HelmRelease{}, errors.New(string(resp.Bytes()))
	}

	release := parseHelmhubRelease(resp.String())
	if !release.PublishedAt.IsZero() {
		return release, nil
	}

	return HelmRelease{}, errors.New("Could not find any helmrelease for " + repositoryName)
}

func parseHelmhubRelease(resp string) HelmRelease {
	var helmRelease HelmRelease

	helmRelease.Name = gjson.Get(resp, "data.attributes.name").String()
	helmRelease.TagName = "chart: " + gjson.Get(resp, "data.relationships.latestChartVersion.data.version").String()
	helmRelease.PublishedAt = gjson.Get(resp, "data.relationships.latestChartVersion.data.created").Time()
	helmRelease.RepositoryName = gjson.Get(resp, "data.attributes.repo.name").String() + "/" + helmRelease.Name
	helmRelease.HtmlUrl = "https://hub.helm.sh/charts/" + helmRelease.RepositoryName + "/" + gjson.Get(resp, "data.relationships.latestChartVersion.data.version").String()
	helmRelease.Type = "Helm chart"

	return helmRelease
}

func main() {
	r := gin.Default()

	r.GET("/api/releases", func(c *gin.Context) {
		var releases []Release

		for _, repositoryName := range getGithubRepositories() {
			release, err := getGithubRelease(repositoryName)
			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}

			days := int(time.Now().Sub(release.PublishedAt).Hours() / 24)
			release.Days = days
			release.Severity = getSeverity(days)
			release.RepositoryName = repositoryName
			releases = append(releases, release)
		}

		sort.Sort(DateSorter(releases))

		c.JSON(http.StatusOK, releases)
	})

	r.GET("api/helmreleases", func(c *gin.Context) {
		var helmReleases []HelmRelease

		for _, repositoryName := range getHelmhubRepositories() {
			release, err := getHelmhubRelease(repositoryName)
			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}

			days := int(time.Now().Sub(release.PublishedAt).Hours() / 24)
			release.Days = days
			release.Severity = getSeverity(days)
			release.RepositoryName = repositoryName
			helmReleases = append(helmReleases, release)
		}

		sort.Sort(HelmDateSorter(helmReleases))

		c.JSON(http.StatusOK, helmReleases)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 by default
	// set environment variable PORT if you want to change port
}
