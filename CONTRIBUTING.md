# Contributing

Contributions to this project are more than just welcome, we are in fact seeking active participation and collaboartion of the ad-tech community.

There are many way to get involved. It can range from feedsback, thoughts and ideas to actual code.

### Reporting Bugs

If you find an issue with the code or missing documentation, you can help us by submitting a [GitHub][github] issue. Even better, if obvious you can correct the problem and submit a Pull Request with a fix.

### Requesting Features

As with bugs, you can request new features via [GitHub][github] issue tracker, the same rules apply.

### Pull Requests

If you have never or rarely contributed code before, please read this excellent [Guide][guide]. Before submitting a PR, please review the list of existing pull requests to ensure that you don't sumit a duplicate. Once confirm, please follow the steps below:

* Fork the project, then clone the repo:

    ```shell
    git clone git@github.com:your-username/flood.git
    ```

* Make your changes in a new git branch:

    ```shell
    git checkout -b my-feature-branch master
    ```

* Apply your code change, please ensure to **include sufficient tests** and documentation
* Run the full test suite, as described in our [README][readme]. Ensure that all tests pass.
* Commit your changes using a descriptive commit message
* Push your feature branch back to [GitHub][github]:

    ```shell
    git push origin my-feature-branch
    ```

* In GitHub, send us a pull request to `flood:master`.
* When changes to your code are suggested:
  * Update your code in your feature-branch. Ensure that all tests (still) pass.
  * Rebase your branch and force push to your GitHub repository (this will update your Pull Request):

    ```shell
    git rebase master -i
    git push origin my-feature-branch -f
    ```

### Code Guidelines

[Effective Go][golang] is an excellent resource. Most importantly, please utilise [gofmt] or [goimports]. Additionally, we recommend following our [EditorConfig][editorcfg] standards.

[github]: https://github.com/bsm/flood
[readme]: https://github.com/bsm/flood/tree/master/README.md
[guide]: https://guides.github.com/activities/contributing-to-open-source/#contributing
[golang]: http://golang.org/doc/effective_go.html
[gofmt]: https://golang.org/cmd/gofmt/
[goimports]: https://godoc.org/golang.org/x/tools/cmd/goimports
[editorcfg]: http://editorconfig.org/
