# Contribution Guide

## Commitizenの導入

本プロジェクトでは、[Commitizen](https://github.com/commitizen/cz-cli)を利用してcommit messageの統一を行っております。

```
npm install commitizen
npm install cz-conventional-changelog-ja
echo '{ "path": "cz-conventional-changelog-ja" }' >> ~/.czrc
exec $SHELL -l
```

## Commitの方法

`npx cz`コマンドを利用してcommitをしてください。

```
git add <file>...
npx cz
```
