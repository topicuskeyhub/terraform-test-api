# Updating and releasing Topicus KeyHub Terraform tester api/container

## 1. Updating

### 1.1 Dependencies

Update the go dependencies

```Shell
go get -u
go mod tidy
```

Update the Terraform version used in `docker/build.sh`.

### 1.2 Commit the results

```Shell
git add .
git commit -m "Upgrade dependencies"
git push
```

No need to release, this is only a tester container.
The project is set up to build master automatically on pushes and we just always the latest.
