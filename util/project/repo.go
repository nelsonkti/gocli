package project

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/mod/modfile"
	"log"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	scpSyntaxRe = regexp.MustCompile(`^(\w+)@([\w.-]+):(.*)$`)
	scheme      = []string{"git", "https", "http", "git+ssh", "ssh", "file", "ftp", "ftps"}
)

var unExpandVarPath = []string{"~", ".", ".."}

// Repo is git repository manager.
type Repo struct {
	url    string
	home   string
	branch string
}

// repoDir 返回存储库的目录名。
func repoDir(url string) string {
	vcsURL, err := ParseVCSUrl(url)
	if err != nil {
		return url
	}
	host, _, err := net.SplitHostPort(vcsURL.Host)
	if err != nil {
		host = vcsURL.Host
	}
	for _, p := range unExpandVarPath {
		host = strings.TrimLeft(host, p)
	}
	dir := path.Base(path.Dir(vcsURL.Path))
	return fmt.Sprintf("%s/%s", host, dir)
}

// NewRepo 创建一个新的存储库管理器。
func NewRepo(url string, branch string) *Repo {
	return &Repo{
		url:    url,
		home:   HomeWithDir("repo/" + repoDir(url)),
		branch: branch,
	}
}

// Path 返回存储库缓存路径。
func (r *Repo) Path() string {
	start := strings.LastIndex(r.url, "/")
	end := strings.LastIndex(r.url, ".git")
	if end == -1 {
		end = len(r.url)
	}
	branch := "@main"
	if r.branch != "" {
		branch = "@" + r.branch
	}
	return path.Join(r.home, r.url[start+1:end]+branch)
}

// Pull 从远程URL获取存储库。
func (r *Repo) Pull(ctx context.Context) error {
	// 保存当前工作目录
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// 切换到目标目录
	if err := os.Chdir(r.Path()); err != nil {
		return fmt.Errorf("failed to switch to repo directory: %w", err)
	}
	defer func() {
		// 恢复到原工作目录
		if chdirErr := os.Chdir(originalDir); chdirErr != nil {
			log.Printf("failed to revert to original directory: %v", chdirErr)
		}
	}()
	cmd := exec.CommandContext(ctx, "git", "symbolic-ref", "HEAD")
	cmd.Dir = r.Path()
	if _, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to check out branch: %w", err)
	}
	cmd = exec.CommandContext(ctx, "git", "pull")
	cmd.Dir = r.Path()
	out, err := cmd.CombinedOutput()
	log.Printf("git pull output: %s", string(out))
	return err
}

// Clone 将存储库克隆到缓存路径。
func (r *Repo) Clone(ctx context.Context) error {
	if _, err := os.Stat(r.Path()); !os.IsNotExist(err) {
		fmt.Println("repo already exists, skipping clone, path:", r.Path())
		return r.Pull(ctx)
	}
	var cmd *exec.Cmd
	if r.branch == "" {
		cmd = exec.CommandContext(ctx, "git", "clone", r.url, r.Path())
	} else {
		cmd = exec.CommandContext(ctx, "git", "clone", "-b", r.branch, r.url, r.Path())
	}
	out, err := cmd.CombinedOutput()
	log.Printf("git clone output: %s", string(out))
	return err
}

// CopyTo 将存储库复制到项目路径。
func (r *Repo) CopyTo(ctx context.Context, to string, modPath string, ignores []string) error {
	if err := r.Clone(ctx); err != nil {
		return err
	}
	mod, err := ModulePath(filepath.Join(r.Path(), "go.mod"))
	if err != nil {
		return err
	}
	return copyDir(r.Path(), to, []string{mod, modPath}, ignores)
}

// CopyToV2 将存储库复制到项目路径。
func (r *Repo) CopyToV2(ctx context.Context, to string, modPath string, ignores, replaces []string) error {
	if err := r.Clone(ctx); err != nil {
		return err
	}
	mod, err := ModulePath(filepath.Join(r.Path(), "go.mod"))
	if err != nil {
		return err
	}
	replaces = append([]string{mod, modPath}, replaces...)
	return copyDir(r.Path(), to, replaces, ignores)
}

func ParseVCSUrl(repo string) (*url.URL, error) {
	var (
		repoURL *url.URL
		err     error
	)

	if m := scpSyntaxRe.FindStringSubmatch(repo); m != nil {
		repoURL = &url.URL{
			Scheme: "ssh",
			User:   url.User(m[1]),
			Host:   m[2],
			Path:   m[3],
		}
	} else {
		if !strings.Contains(repo, "//") {
			repo = "//" + repo
		}
		if strings.HasPrefix(repo, "//git@") {
			repo = "ssh:" + repo
		} else if strings.HasPrefix(repo, "//") {
			repo = "https:" + repo
		}
		repoURL, err = url.Parse(repo)
		if err != nil {
			return nil, err
		}
	}

	for _, s := range scheme {
		if repoURL.Scheme == s {
			return repoURL, nil
		}
	}
	return nil, errors.New("unable to parse repo url")
}

// ModulePath 返回go模块路径
func ModulePath(filename string) (string, error) {
	modBytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return modfile.ModulePath(modBytes), nil
}
