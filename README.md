![](https://github.com/ottodono/task-cli/actions/workflows/go.yml/badge.svg)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)
[![Go Report Card](https://goreportcard.com/badge/github.com/ottodono/task-cli)](https://goreportcard.com/report/github.com/ottodono/task-cli)
___
# task-cli
Simple CLI task manager

# Branch Naming Convention
The Git Branching Naming Convention article is an excellent base.
However, you can simplify even more.

## Category

A git branch should start with a category. Pick one of these: `feature`, `bugfix`, `hotfix`, or `test`.

* `feature` is for adding, refactoring or removing a feature
* `bugfix` is for fixing a bug
* `hotfix` is for changing code with a temporary solution and/or without following the usual process (usually because of an emergency)
* `test` is for experimenting outside an issue/ticket 

## Reference
After the category, there should be a `"/"` followed by the reference of the issue/ticket you are working on. If there's no reference, just add `no-ref`.

## Description

After the reference, there should be another "/" followed by a description which sums up the purpose of this specific branch. This description should be short and "kebab-cased".

By default, you can use the title of the issue/ticket you are working on. Just replace any special character by `"-"`.

## To sum up, follow this pattern when branching:
```
git branch <category/reference/description-in-kebab-case>
```

## Examples:

* You need to add, refactor or remove a feature: git branch `feature/issue-42/create-new-button-component`
* You need to fix a bug: git branch `bugfix/issue-342/button-overlap-form-on-mobile`
* You need to fix a bug really fast (possibly with a temporary solution): git branch `hotfix/no-ref/registration-form-not-working`
* You need to experiment outside an issue/ticket: git branch `test/no-ref/refactor-components-with-atomic-design`

# Conventional Commit Messages <sub><img src="https://img.icons8.com/dusk/1600/commit-git.png" height="64" /></sub>
> [!TIP]
> Have a look at **[git-conventional-commits](https://github.com/qoomon/git-conventional-commits)** , a CLI util to ensure these conventions, determine version and generate changelogs

## Commit Message Formats

### Default
<pre>
<b><a href="#types">&lt;type&gt;</a></b></font>(<b><a href="#scopes">&lt;optional scope&gt;</a></b>): <b><a href="#description">&lt;description&gt;</a></b>
<sub>empty separator line</sub>
<b><a href="#body">&lt;optional body&gt;</a></b>
<sub>empty separator line</sub>
<b><a href="#footer">&lt;optional footer&gt;</a></b>
</pre>

### Merge Commit
<pre>
Merge branch '<b>&lt;branch name&gt;</b>'
</pre>
<sup>Follows default git merge message</sup>

### Revert Commit
<pre>
Revert "<b>&lt;reverted commit subject line&gt;</b>"
</pre>
<sup>Follows default git revert message</sup>

### Inital Commit
```
chore: init
```

### Types
* API relevant changes
    * `feat` Commits, that adds or remove a new feature
    * `fix` Commits, that fixes a bug
* `refactor` Commits, that rewrite/restructure your code, however does not change any API behaviour
    * `perf` Commits are special `refactor` commits, that improve performance
* `style` Commits, that do not affect the meaning (white-space, formatting, missing semi-colons, etc)
* `test` Commits, that add missing tests or correcting existing tests
* `docs` Commits, that affect documentation only
* `build` Commits, that affect build components like build tool, ci pipeline, dependencies, project version, ...
* `ops` Commits, that affect operational components like infrastructure, deployment, backup, recovery, ...
* `chore` Miscellaneous commits e.g. modifying `.gitignore`

### Scopes
The `scope` provides additional contextual information.
* Is an **optional** part of the format
* Allowed Scopes depends on the specific project
* Don't use issue identifiers as scopes

### Breaking Changes Indicator
Breaking changes should be indicated by an `!` before the `:` in the subject line e.g. `feat(api)!: remove status endpoint`
* Is an **optional** part of the format

### Description
The `description` contains a concise description of the change.
* Is a **mandatory** part of the format
* Use the imperative, present tense: "change" not "changed" nor "changes"
    * Think of `This commit will...` or `This commit should...`
* Don't capitalize the first letter
* No dot (`.`) at the end

### Body
The `body` should include the motivation for the change and contrast this with previous behavior.
* Is an **optional** part of the format
* Use the imperative, present tense: "change" not "changed" nor "changes"
* This is the place to mention issue identifiers and their relations

### Footer
The `footer` should contain any information about **Breaking Changes** and is also the place to **reference Issues** that this commit refers to.
* Is an **optional** part of the format
* **optionally** reference an issue by its id.
* **Breaking Changes** should start with the word `BREAKING CHANGES:` followed by space or two newlines. The rest of the commit message is then used for this.


### Examples
* ```
  feat: add email notifications on new direct messages
  ```
* ```
  feat(shopping cart): add the amazing button
  ```
* ```
  feat!: remove ticket list endpoint

  refers to JIRA-1337

  BREAKING CHANGES: ticket enpoints no longer supports list all entites.
  ```
* ```
  fix(shopping-cart): prevent order an empty shopping cart
  ```
* ```
  fix(api): fix wrong calculation of request body checksum
  ```
* ```
  fix: add missing parameter to service call

  The error occurred because of <reasons>.
  ```
* ```
  perf: decrease memory footprint for determine uniqe visitors by using HyperLogLog
  ```
* ```
  build: update dependencies
  ```
* ```
  build(release): bump version to 1.0.0
  ```
* ```
  refactor: implement fibonacci number calculation as recursion
  ```
* ```
  style: remove empty line
  ```

---

## Git Hook Scripts to ensure commit message header format
<details>
<summary>Click to expand</summary>

### commit-msg Hook (local)
* Create a commit-msg hook using [git-conventional-commits cli](https://github.com/qoomon/git-conventional-commits?tab=readme-ov-file#automatically-validate-commit-message-convention-before-commit)

### pre-receive Hook (server side)
* create following file in your repository folder `.git/hooks/pre-receive`
  ```shell
  #!/usr/bin/env bash

  # Pre-receive hook that will block commits with messges that do not follow regex rule

  commit_msg_type_regex='feat|fix|refactor|style|test|docs|build'
  commit_msg_scope_regex='.{1,20}'
  commit_msg_description_regex='.{1,100}'
  commit_msg_regex="^(${commit_msg_type_regex})(\(${commit_msg_scope_regex}\))?: (${commit_msg_description_regex})\$"
  merge_msg_regex="^Merge branch '.+'\$"

  zero_commit="0000000000000000000000000000000000000000"

  # Do not traverse over commits that are already in the repository
  excludeExisting="--not --all"

  error=""
  while read oldrev newrev refname; do
    # branch or tag get deleted
    if [ "$newrev" = "$zero_commit" ]; then
      continue
    fi

    # Check for new branch or tag
    if [ "$oldrev" = "$zero_commit" ]; then
      rev_span=`git rev-list $newrev $excludeExisting`
    else
      rev_span=`git rev-list $oldrev..$newrev $excludeExisting`
    fi

    for commit in $rev_span; do
      commit_msg_header=$(git show -s --format=%s $commit)
      if ! [[ "$commit_msg_header" =~ (${commit_msg_regex})|(${merge_msg_regex}) ]]; then
        echo "$commit" >&2
        echo "ERROR: Invalid commit message format" >&2
        echo "$commit_msg_header" >&2
        error="true"
      fi
    done
  done

  if [ -n "$error" ]; then
    exit 1
  fi
  ```
* ⚠ make `.git/hooks/pre-receive` executable (unix: `chmod +x '.git/hooks/pre-receive'`)

</details>

-----
## References
* https://www.conventionalcommits.org/
* https://github.com/angular/angular/blob/master/CONTRIBUTING.md
* http://karma-runner.github.io/1.0/dev/git-commit-msg.html
  <br>

* https://github.com/github/platform-samples/tree/master/pre-receive-hooks
* https://github.community/t5/GitHub-Enterprise-Best-Practices/Using-pre-receive-hooks-in-GitHub-Enterprise/ba-p/13863
