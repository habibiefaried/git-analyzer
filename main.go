package main

import (
        git "github.com/libgit2/git2go/v33"
        "log"
)

func main() {
        var lines   []string
        sourceBranchName := "feat/testdoang"
        destBranchName := "master"

        repo, err := git.OpenRepository("/testrepo")
        if err != nil {
                log.Fatal(err)
        }
        defer repo.Free()

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

        // fmt.Println(commitSrc.Message())
        // fmt.Println(commitDest.Message())

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

        // Show all file patch diffs in a commit.
        log.Println("========================================================")
        numDeltas, err := gitDiff.NumDeltas()
        if err != nil {
            log.Fatal(err)
        }

        for d := 0; d < numDeltas; d++ {
            patch, err := gitDiff.Patch(d)
            if err != nil {
                log.Fatal(err)
            }
            patchString, err := patch.String()
            if err != nil {
                log.Fatal(err)
            }
            log.Printf("\n%s", patchString)
            patch.Free()
        }

        log.Println("========================================================")
        gitDiff.ForEach(func(file git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
            lines = append(lines, file.OldFile.Path)
            return nil, nil
        }, git.DiffDetailFiles)

        log.Println(lines)
}