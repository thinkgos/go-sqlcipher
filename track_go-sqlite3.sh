#!/bin/sh -e

if [ $# -ne 1 ]; then
  echo "Usage: $0 go-sqlite3_dir" >&2
  echo "Copy tracked source files from go-sqlite3 to current directory." >&2
  exit 1
fi

ltd=$1

# copy C files
cp -f $ltd/sqlite3_opt_unlock_notify.c .

# copy Go files
cp -f $ltd/*.go .
rm -rf _example
cp -r $ltd/_example .
rm -rf upgrade
cp -r $ltd/upgrade .
cp -f $ltd/README.md README_GO_SQLITE3.md

# 定义需要替换的字符串
OLD_PACKAGE="github.com/mattn/go-sqlite3"
NEW_PACKAGE="github.com/thinkgos/go-sqlcipher"

# 查找并替换文件中的字符串
find "_example" -type f | while read -r file; do
  sed -i "s|$OLD_PACKAGE|$NEW_PACKAGE|g" "$file"
done
sed -i "s|$OLD_PACKAGE|$NEW_PACKAGE|g" "sqlite3.go"
sed -i "s|$OLD_PACKAGE|$NEW_PACKAGE|g" "doc.go"

echo "make sure to adjust sqlite3.go with sqlcipher pragmas!!!"
