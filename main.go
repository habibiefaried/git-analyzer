package main

import (
	git "github.com/libgit2/git2go/v33"
	"log"
	"os"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	destBranchName := "master"

	repo, err := git.OpenRepository(getEnv("ANALYZE_REPO", "."))
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Free()

	currentConfig, err := repo.Config()
	if err != nil {
		log.Fatal(err)
	}
	defer currentConfig.Free()

	ci, err := currentConfig.NewIterator()
	if err != nil {
		log.Fatal(err)
	}
	defer ci.Free()

	for {
		ce, err := ci.Next()
		if err != nil || ce == nil {
			break
		}

		if ce.Name == "remote.origin.url" {
			log.Println("Analyzing repo: " + (ce.Value))
			break
		}
	}

	currentRef, err := repo.Head()
	if err != nil {
		log.Fatal(err)
	}
	defer currentRef.Free()

	sourceBranchName, err := currentRef.Branch().Name()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Comparing branch %s to %s ", sourceBranchName, destBranchName)

	sourceBranch, err := repo.LookupBranch(sourceBranchName, git.BranchLocal)
	if err != nil {
		log.Fatal("Failed lookup source branch " + sourceBranchName)
	}
	defer sourceBranch.Free()

	commitSrc, err := repo.LookupCommit(sourceBranch.Target())
	if err != nil {
		log.Fatal("Failed to find remote branch commit: " + sourceBranchName)
	}
	defer commitSrc.Free()

	destBranch, err := repo.LookupBranch(destBranchName, git.BranchLocal)
	if err != nil {
		log.Fatal("Failed lookup source branch " + destBranchName)
	}
	defer destBranch.Free()

	commitDest, err := repo.LookupCommit(destBranch.Target())
	if err != nil {
		log.Fatal("Failed to find remote branch commit: " + destBranchName)
	}
	defer commitDest.Free()

	commitSrcTree, err := commitSrc.Tree()
	if err != nil {
		log.Fatal(err)
	}

	commitDestTree, err := commitDest.Tree()
	if err != nil {
		log.Fatal(err)
	}

	options, err := git.DefaultDiffOptions()
	if err != nil {
		log.Fatal(err)
	}

	gitDiff, err := repo.DiffTreeToTree(commitSrcTree, commitDestTree, &options)
	if err != nil {
		log.Fatal(err)
	}

	gitDiff.ForEach(func(file git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
		log.Println(file.OldFile.Path)
		return nil, nil
	}, git.DiffDetailFiles)
}
