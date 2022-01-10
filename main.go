package main

import (
    "github.com/libgit2/git2go"
    "fmt"
    "log"
)

func main(){
    repo, err := git.OpenRepository(".\Vulnerability-goapp")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(repo)
}