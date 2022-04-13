# Weave Gitops Release Process

To release a new version of Weave Gitops, you need to:
- Decide on a new release number
- Trigger the github action that prepares for that release
- Merge the PR
- Check out the merge commit, and run `git tag -a $version_number -m $version_number && git push tag $version_number`
- Add a record of the new version in the checkpoint system

# Record the new version
- Add a record in the [checkpoint system](https://checkpoint-api.weave.works/admin) to inform users of the new version.  The CLI checks for a more recent version and informs the user where to download it based on this data.
  - Record must match this template:
     ```
    Name: weave-gitops
    Version: N.N.N
    Release date: (current date in UTC. i.e.: 2021-09-08 12:41:00 )
    Download URL: https://github.com/weaveworks/weave-gitops/releases/tag/vN.N.N
    Changelog URL: https://github.com/weaveworks/weave-gitops/releases/tag/vN.N.N
    Project Website: https://www.weave.works/product/gitops-core/
    ```
  - _note: A Weaveworks employee must perform this step_

That's it!
