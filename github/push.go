package github

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	. "github.com/xiaosumay/server-code-mgr/utils"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type pushPayload struct {
	Ref     string      `json:"ref"`
	Before  string      `json:"before"`
	After   string      `json:"after"`
	Created bool        `json:"created"`
	Deleted bool        `json:"deleted"`
	Forced  bool        `json:"forced"`
	BaseRef interface{} `json:"base_ref"`
	Compare string      `json:"compare"`
	Commits []struct {
		ID        string    `json:"id"`
		TreeID    string    `json:"tree_id"`
		Distinct  bool      `json:"distinct"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		URL       string    `json:"url"`
		Author    struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"author"`
		Committer struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"committer"`
		Added    []interface{} `json:"added"`
		Removed  []interface{} `json:"removed"`
		Modified []string      `json:"modified"`
	} `json:"commits"`
	HeadCommit struct {
		ID        string    `json:"id"`
		TreeID    string    `json:"tree_id"`
		Distinct  bool      `json:"distinct"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		URL       string    `json:"url"`
		Author    struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"author"`
		Committer struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"committer"`
		Added    []interface{} `json:"added"`
		Removed  []interface{} `json:"removed"`
		Modified []string      `json:"modified"`
	} `json:"head_commit"`
	Repository struct {
		ID       int    `json:"id"`
		NodeID   string `json:"node_id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Private  bool   `json:"private"`
		Owner    struct {
			Name              string `json:"name"`
			Email             string `json:"email"`
			Login             string `json:"login"`
			ID                int    `json:"id"`
			NodeID            string `json:"node_id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"owner"`
		HTMLURL          string      `json:"html_url"`
		Description      interface{} `json:"description"`
		Fork             bool        `json:"fork"`
		URL              string      `json:"url"`
		ForksURL         string      `json:"forks_url"`
		KeysURL          string      `json:"keys_url"`
		CollaboratorsURL string      `json:"collaborators_url"`
		TeamsURL         string      `json:"teams_url"`
		HooksURL         string      `json:"hooks_url"`
		IssueEventsURL   string      `json:"issue_events_url"`
		EventsURL        string      `json:"events_url"`
		AssigneesURL     string      `json:"assignees_url"`
		BranchesURL      string      `json:"branches_url"`
		TagsURL          string      `json:"tags_url"`
		BlobsURL         string      `json:"blobs_url"`
		GitTagsURL       string      `json:"git_tags_url"`
		GitRefsURL       string      `json:"git_refs_url"`
		TreesURL         string      `json:"trees_url"`
		StatusesURL      string      `json:"statuses_url"`
		LanguagesURL     string      `json:"languages_url"`
		StargazersURL    string      `json:"stargazers_url"`
		ContributorsURL  string      `json:"contributors_url"`
		SubscribersURL   string      `json:"subscribers_url"`
		SubscriptionURL  string      `json:"subscription_url"`
		CommitsURL       string      `json:"commits_url"`
		GitCommitsURL    string      `json:"git_commits_url"`
		CommentsURL      string      `json:"comments_url"`
		IssueCommentURL  string      `json:"issue_comment_url"`
		ContentsURL      string      `json:"contents_url"`
		CompareURL       string      `json:"compare_url"`
		MergesURL        string      `json:"merges_url"`
		ArchiveURL       string      `json:"archive_url"`
		DownloadsURL     string      `json:"downloads_url"`
		IssuesURL        string      `json:"issues_url"`
		PullsURL         string      `json:"pulls_url"`
		MilestonesURL    string      `json:"milestones_url"`
		NotificationsURL string      `json:"notifications_url"`
		LabelsURL        string      `json:"labels_url"`
		ReleasesURL      string      `json:"releases_url"`
		DeploymentsURL   string      `json:"deployments_url"`
		CreatedAt        int         `json:"created_at"`
		UpdatedAt        time.Time   `json:"updated_at"`
		PushedAt         int         `json:"pushed_at"`
		GitURL           string      `json:"git_url"`
		SSHURL           string      `json:"ssh_url"`
		CloneURL         string      `json:"clone_url"`
		SvnURL           string      `json:"svn_url"`
		Homepage         interface{} `json:"homepage"`
		Size             int         `json:"size"`
		StargazersCount  int         `json:"stargazers_count"`
		WatchersCount    int         `json:"watchers_count"`
		Language         string      `json:"language"`
		HasIssues        bool        `json:"has_issues"`
		HasProjects      bool        `json:"has_projects"`
		HasDownloads     bool        `json:"has_downloads"`
		HasWiki          bool        `json:"has_wiki"`
		HasPages         bool        `json:"has_pages"`
		ForksCount       int         `json:"forks_count"`
		MirrorURL        interface{} `json:"mirror_url"`
		Archived         bool        `json:"archived"`
		OpenIssuesCount  int         `json:"open_issues_count"`
		License          struct {
			Key    string `json:"key"`
			Name   string `json:"name"`
			SpdxID string `json:"spdx_id"`
			URL    string `json:"url"`
			NodeID string `json:"node_id"`
		} `json:"license"`
		Forks         int    `json:"forks"`
		OpenIssues    int    `json:"open_issues"`
		Watchers      int    `json:"watchers"`
		DefaultBranch string `json:"default_branch"`
		Stargazers    int    `json:"stargazers"`
		MasterBranch  string `json:"master_branch"`
	} `json:"repository"`
	Pusher struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"pusher"`
	Sender struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"sender"`
}

func PushEvent(data []byte) bool {
	var push pushPayload
	err := json.Unmarshal(data, &push)
	if err != nil {
		log.Println(err)
		return false
	}

	repoName := push.Repository.Name

	if repo, ok := Repositories[repoName]; ok {
		ref := "refs/heads/" + DefaultValue(repo.Branch, "master")
		if push.Ref == ref {
			go doReposUpdate(repoName, repo)
			return true
		}
	}

	log.Println(repoName + " 不存在！")
	return false
}

func runCommand(repoName string, rep Repo) bool {
	if _, err := os.Stat(rep.Script); err != nil {
		log.Println(err)
		return false
	}

	cmd := exec.Command("bash", rep.Script)
	cmd.Env = append(cmd.Env, "BRANCH="+Quote(rep.Branch), "WORK_PATH="+Quote(rep.Path), "REPOS="+Quote(repoName))
	if 0 != len(rep.Key) {
		key := rep.Key
		if _, err := os.Stat(rep.Key); err != nil {
			key = strings.Join([]string{os.Getenv("HOME"), ".ssh", rep.Key}, string(os.PathSeparator))
		}

		cmd.Env = append(cmd.Env, "GIT_SSH_COMMAND=ssh -v -i "+Quote(key))
	}

	log.Println(strings.Join(cmd.Env, " "))

	data, err := cmd.CombinedOutput()

	if err != nil {
		log.Println(err)
		return false
	}

	log.Println(string(data))
	return true
}

func doReposUpdate(repoName string, rep Repo) {

	for {
		if len(rep.Script) == 0 {
			break
		}

		log.Printf("启用自定义脚本: %s\n", rep.Script)

		if runCommand(repoName, rep) {
			return
		}

		break
	}

	localPath := DefaultValue(rep.Path, "/home/wwwroot/"+repoName)

	if _, err := os.Stat(localPath); err != nil {
		cloneRepos(repoName, rep.RemotePath, rep)
	}

	r, err := git.PlainOpenWithOptions(localPath, &git.PlainOpenOptions{
		DetectDotGit: false,
	})

	if err != nil {
		log.Println(err)
		return
	}

	remoteRef, err := r.Reference(
		plumbing.NewRemoteReferenceName("origin", DefaultValue(rep.Branch, "master")),
		true,
	)
	if err != nil {
		log.Println(err)
		return
	}

	localRef, err := r.Reference(plumbing.ReferenceName("HEAD"), true)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(remoteRef)
	log.Println(localRef)

	if remoteRef.Hash() == localRef.Hash() {
		log.Println("已经是最新的了！")
		return
	}

	auth, err := getAuth(rep.Key)
	if err != nil {
		log.Println(err)
		return
	}

	err = r.Fetch(&git.FetchOptions{
		Auth:     auth,
		Force:    true,
		Progress: os.Stdout,
		Tags:     git.AllTags,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		log.Println(err)
		return
	}

	log.Println("强制拉去完成")

	w, err := r.Worktree()
	if err != nil {
		log.Println(err)
		return
	}

	err = w.Reset(&git.ResetOptions{
		Commit: remoteRef.Hash(),
	})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("更新完成")
}
