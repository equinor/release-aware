package main

import (
	"errors"
	// "fmt"
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
	TagName        string `json:"tag_name"`
	PublishedAt    string `json:"published_at"`
	CreatedAt      string `json:"created_at"`
	Severity       string `json:"severity"`
	RepositoryName string `json:"repository_name"`
	HtmlUrl        string `json:"html_url"`
	Days           int    `json:"days"`
	Type           string `json:"type"`
	date           time.Time
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
	return a[i].date.After(a[j].date)
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

func parseRelease(resp string, repositoryName string, prefix string) Release {
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
	release.PublishedAt = gjson.Get(resp, prefix+".createdAt").String()
	release.TagName = gjson.Get(resp, prefix+".tagName").String()
	return release
}

func parseTag(resp string, repositoryName string, prefix string) Release {
	var release Release
	tagName := gjson.Get(resp, prefix+".name")
	pushedDate := gjson.Get(resp, prefix+".target.pushedDate")
	if !pushedDate.Exists() {
		pushedDate = gjson.Get(resp, prefix+".target.tagger.date")
	}
	release.PublishedAt = pushedDate.String()
	release.TagName = tagName.String()
	release.HtmlUrl = "https://github.com/" + repositoryName + "/releases/tag/" + tagName.String()
	release.Type = "Tag"
	return release
}

func getGithubRelease(repositoryName string) (Release, error) {
	splitted := strings.Split(repositoryName, "/")
	query := `{
		repository(owner: "` + splitted[0] + `", name: "` + splitted[1] + `") {
		  releases (last: 1) {
		     edges {
		  		node {
			  	  createdAt
                  url
	    		  isDraft
  		          isPrerelease
                  tagName
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

	prefix_release := "data.repository.releases.edges.0.node"
	prefix_tag := "data.repository.refs.edges.0.node"

	isRelease := gjson.Get(resp.String(), prefix_release).Exists()
	isTag := gjson.Get(resp.String(), prefix_tag).Exists()

	if isRelease {
		return parseRelease(resp.String(), repositoryName, prefix_release), nil
	} else if isTag {
		return parseTag(resp.String(), repositoryName, prefix_tag), nil
	}
	return Release{}, errors.New("Could not find any release for " + repositoryName)
}

func main() {
	r := gin.Default()

	r.GET("/api/releases", func(c *gin.Context) {
		releases := []Release{}

		for _, repositoryName := range getGithubRepositories() {
			release, err := getGithubRelease(repositoryName)
			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}

			// Get number of days since last release			
			date, err := time.Parse(time.RFC3339, release.PublishedAt)
			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}
			t2 := time.Now()
			days := int(t2.Sub(date).Hours() / 24)

			// date is used when sorting
			release.date = date
			release.Days = days
			release.Severity = getSeverity(days)
			release.RepositoryName = repositoryName
			releases = append(releases, release)
		}

		sort.Sort(DateSorter(releases))

		c.JSON(http.StatusOK, releases)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 by default
	// set environment variable PORT if you want to change port
}
