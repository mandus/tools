package main

import (
	"fmt"
	"os"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

const maxTraversal = 1000

func main() {
	repo, err := gogit.PlainOpenWithOptions(".", &gogit.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		os.Exit(0)
	}

	head, err := repo.Head()
	if err != nil {
		// No commits yet (initial repo)
		fmt.Print("\033[;32m(HEAD)\033[0m")
		return
	}

	var branch string
	if head.Name().IsBranch() {
		branch = head.Name().Short()
	} else {
		branch = "HEAD" // detached
	}

	sync := computeSync(repo, head)

	tag := ""
	if os.Getenv("GIT_PROMPT_DISABLE_TAGS") != "1" {
		tag = nearestTag(repo, head.Hash())
		if tag != "" {
			tag = "(" + tag + ")"
		}
	}

	fmt.Printf("\033[;32m(%s%s%s)\033[0m", branch, tag, sync)
}

func computeSync(repo *gogit.Repository, head *plumbing.Reference) string {
	if !head.Name().IsBranch() {
		return ""
	}

	cfg, err := repo.Config()
	if err != nil {
		return ""
	}

	branchName := head.Name().Short()
	branchCfg, ok := cfg.Branches[branchName]
	if !ok || branchCfg.Remote == "" || branchCfg.Merge == "" {
		return ""
	}

	remoteRefName := plumbing.NewRemoteReferenceName(branchCfg.Remote, branchCfg.Merge.Short())
	remoteRef, err := repo.Reference(remoteRefName, true)
	if err != nil {
		return ""
	}

	localHash := head.Hash()
	remoteHash := remoteRef.Hash()
	if localHash == remoteHash {
		return "="
	}

	ahead, behind := aheadBehind(repo, localHash, remoteHash)
	switch {
	case ahead > 0 && behind > 0:
		return "<>"
	case ahead > 0:
		return ">"
	case behind > 0:
		return "<"
	default:
		return "="
	}
}

func aheadBehind(repo *gogit.Repository, local, remote plumbing.Hash) (int, int) {
	localSet := reachable(repo, local)
	remoteSet := reachable(repo, remote)

	ahead := 0
	for h := range localSet {
		if _, ok := remoteSet[h]; !ok {
			ahead++
		}
	}
	behind := 0
	for h := range remoteSet {
		if _, ok := localSet[h]; !ok {
			behind++
		}
	}
	return ahead, behind
}

func reachable(repo *gogit.Repository, start plumbing.Hash) map[plumbing.Hash]struct{} {
	seen := make(map[plumbing.Hash]struct{})
	queue := []plumbing.Hash{start}
	for len(queue) > 0 && len(seen) < maxTraversal {
		h := queue[0]
		queue = queue[1:]
		if _, ok := seen[h]; ok {
			continue
		}
		seen[h] = struct{}{}
		c, err := repo.CommitObject(h)
		if err != nil {
			continue
		}
		queue = append(queue, c.ParentHashes...)
	}
	return seen
}

func nearestTag(repo *gogit.Repository, head plumbing.Hash) string {
	tagMap := make(map[plumbing.Hash]string)
	tags, err := repo.Tags()
	if err != nil {
		return ""
	}
	_ = tags.ForEach(func(ref *plumbing.Reference) error {
		h := ref.Hash()
		name := ref.Name().Short()
		if to, err := repo.TagObject(h); err == nil {
			h = to.Target
		}
		if _, exists := tagMap[h]; !exists {
			tagMap[h] = name
		}
		return nil
	})

	if len(tagMap) == 0 {
		return ""
	}

	seen := make(map[plumbing.Hash]struct{})
	queue := []plumbing.Hash{head}
	for len(queue) > 0 {
		h := queue[0]
		queue = queue[1:]
		if _, ok := seen[h]; ok {
			continue
		}
		seen[h] = struct{}{}
		if name, ok := tagMap[h]; ok {
			return name
		}
		c, err := repo.CommitObject(h)
		if err != nil {
			continue
		}
		queue = append(queue, c.ParentHashes...)
	}
	return ""
}
