# Release

This document explains the current release procedure for the Glide project.

A person that handles the release is called the release manager.

The release version should be picked based on the previous version and changes that are going to be released:
- If the previous version was `0.0.3` and changeset doesn't have breaking changes in API, the next version may be either `0.0.4-rc.1`.
- If the previous version was `0.0.3-rc.2` and changeset doesn't have breaking changes in API,  the next version may be either `0.0.3-rc.3` or `0.0.3`.
- if the previous version was `0.0.3` and changeset has breaking changes in API, the next version may be either `0.1.0-rc.1` or `0.1.0`.

It's encouraged to have a release candidate as an intermediate step in the releasing the "stable" version.

## The release process

The following is step-by-step instruction.

### The release PR

The release manager creates a new pull request from the `develop` branch into the `main`.

### Collect changelog

The changeset from the PR should be reviewed and represented as a changelog in the `CHANGELOG.md` file.

Git can be helpful to summarize changes:

```bash
git log --pretty=format:"%s (%an)" {LAST_VERSION_TAG}..HEAD
```

Sort out all changes to the categories they belong to.

### Update the release PR description

Use the created changelog to update the PR description.

### Merge the release PR

Once all checks are green on the Release PR, you can merge it into the main.
When squashing the pull request, please specify the same changelog you have added to the PR description.

### Tag the merge commit

Tag the merge commit in the main branch using the release version like `0.1.0-rc.1`.

That should trigger a deployment pipeline that can be checked under https://github.com/EinStack/glide/actions.

Wait until the pipeline is done or failed.

### Create a GitHub release

The release manager then creates a release in GitHub for the released version.
The GitHub release should contain the same changelog message we have used in the release PR description.

The release should be marked as a pre-release if the version is a release candidate or a stable release otherwise.

Feel free to announce the release via creating a discussion (there is a checkbox for that on the new GitHub release form).

### Spread a word about the release

The release manager should spread a word about the release in the Discord community space.

Please use the `#announcements` channel for announcing stable releases. 
You can use `#general` to announce release candidates.

Be sure to do shutouts for people who spend their time improving the project.

| Releases         | Release Manager |
|------------------|-----------------|
| 0.0.1-0.1.0-rc.1 | @roma-glushko   |
