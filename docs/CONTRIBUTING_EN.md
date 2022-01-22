# Contribution Guide

## Install the Commitizen

We use [Commitizen](https://github.com/commitizen/cz-cli) to unify commit messages in this project.

```
npm install commitizen
npm install cz-conventional-changelog
echo '{ "path": "cz-conventional-changelog" }' >> ~/.czrc
exec $SHELL -l
```

## How to Commit

Use `npx cz` command when you try to commit.

```
git add <file>...
npx cz
```
