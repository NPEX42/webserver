package main

import "github.com/go-git/go-git/v5"

func Pull() error {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = worktree.Pull(&git.PullOptions{RemoteName: "gh"})
	if err != nil {
		return err
	}

	return nil
}
