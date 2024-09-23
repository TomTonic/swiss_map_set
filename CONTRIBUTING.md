# Contributing

(This text is derived from https://github.com/coreinfrastructure/best-practices-badge/blob/main/CONTRIBUTING.md) 

**Feedback and contributions are very welcome!**

Here's help on how to make contributions, divided into the following sections:

* general information,
* vulnerability reporting,
* documentation changes,
* code changes,
* how to check proposed changes before submitting them,
* reuse (supply chain for third-party components, including updating them)

## General information

For specific proposals, please provide them as
[pull requests](https://github.com/TomTonic/Set3/pulls)
or
[issues](https://github.com/TomTonic/Set3/issues)
via our
[GitHub site](https://github.com/coreinfrastructure/best-practices-badge).

You are welcome aboard!

### Pull requests and different branches recommended

Pull requests are preferred, since they are specific.
For more about how to create a pull request, see
<https://help.github.com/articles/using-pull-requests/>.

We recommend creating different branches for different (logical)
changes, and creating a pull request when you're done into the main branch.
See the GitHub documentation on
[creating branches](https://help.github.com/articles/creating-and-deleting-branches-within-your-repository/)
and
[using pull requests](https://help.github.com/articles/using-pull-requests/).

### How we handle proposals

We use GitHub to track proposed changes via its
[issue tracker](https://github.com/TomTonic/Set3/issues) and
[pull requests](https://github.com/TomTonic/Set3/pulls).
Specific changes are proposed using those mechanisms.
Issues are assigned to an individual, who works it and then marks it complete.
If there are questions or objections, the conversation area of that
issue or pull request is used to resolve it.

### Two-person review

Our policy is that as many as possible proposed modifications will be reviewed
before release by a person other than the author,
to determine if it is a worthwhile modification and free of known issues
which would argue against its inclusion.

We achieve this by splitting proposals into two kinds:

1. Low-risk modifications.  These modifications are being proposed by
   people authorized to commit directly, pass all tests, and are unlikely
   to have problems.  These include documentation/text updates and/or updates
   to existing functions (especially minor updates) where no risk (such as a security risk)
   have been identified.  The project lead can decide that any particular
   modification is low-risk.
2. Other modifications.  These other modifications need to be
   reviewed by someone else or the project lead can decide to accept
   the modification.  Typically this is done by creating a branch and a
   pull request so that it can be reviewed before accepting it.

### Developer Certificate of Origin (DCO)

All contributions (including pull requests) must agree to
the Linux kernel developers'
[Developer Certificate of Origin (DCO) version 1.1](https://developercertificate.org).
This is a developer's certification that he or she has the right to
submit the patch for inclusion into the project.

Simply submitting a contribution implies this agreement, however,
please include a "Signed-off-by" tag in every patch
(this tag is a conventional way to confirm that you agree to the DCO).
You can do this with <tt>git commit --signoff</tt> (the <tt>-s</tt> flag
is a synonym for <tt>--signoff</tt>).

Another way to do this is to write the following at the end of the commit
message, on a line by itself separated by a blank line from the body of
the commit:

````
Signed-off-by: YOUR NAME <YOUR.EMAIL@EXAMPLE.COM>
````

You can sign-off by default in this project by creating a file
(say "git-template") that contains
some blank lines and the signed-off-by text above;
then configure git to use that as a commit template.  For example:

````sh
git config commit.template ~/best-practices-badge/git-template
````

It's not practical to fix old contributions in git, so if one is forgotten,
do not try to fix them.  We presume that if someone sometimes used a DCO,
a commit without a DCO is an accident and the DCO still applies.

### We are proactive

In general we try to be proactive to detect and eliminate
mistakes and vulnerabilities as soon as possible,
and to reduce their impact when they do happen.
We use a defensive design and coding style to reduce the likelihood of mistakes,
a variety of tools that try to detect mistakes early,
and an automatic test suite with significant coverage.
We also release the software as open source software so others can review it.

Since early detection and impact reduction can never be perfect, we also try to
detect and repair problems during deployment as quickly as possible.
This is *especially* true for security issues.

## <span id="how_to_report_vulnerabilities">Vulnerability reporting (security issues)</a>

Please privately report vulnerabilities you find, so we can fix them!

See [SECURITY.md](./SECURITY.md) for information on how to privately report vulnerabilities.

## Documentation changes

Most of the documentation is in "markdown" format.
All markdown files use the .md filename extension.

Where reasonable, limit yourself to Markdown
that will be accepted by different markdown processors
(e.g., what is specified by CommonMark or the original Markdown).

## Code changes

The code should strive to be DRY (don't repeat yourself),
clear, and obviously correct.
Some technical debt is inevitable, just don't bankrupt us with it.
Improved refactorizations are welcome.

### Automated tests

When adding or changing functionality, please include new tests for them as
part of your contribution.

We require the Go code to have at least 98% statement coverage;
please ensure your contributions do not lower the coverage below that minimum.
Additional tests are very welcome.

We encourage tests to be created first, run to ensure they fail, and
then add code to implement the test (aka test driven development).
However, each git commit should have both
the test and improvement in the *same* commit,
because 'git bisect' will then work well.

*WARNING*: It is possible that some tests may intermittently fail, even though
the software works fine.
If tests fail, restart to see if it's a problem with the software
or the tests.
Where possible, try to find and fix the problem; we have worked to
eliminate this, and at this point believe we have fixed it.

### Security, privacy, and performance

Pay attention to security, and work *with* (not against) our
security hardening mechanisms.

Protect private information, in particular passwords and email addresses.
Avoid mechanisms that could be used for tracking where possible
(we do need to verify people are logged in for some operations),
and ensure that third parties can't use interactions for tracking.
When sending an email to an existing account, use the original account
email not the claimed email address sent now; for why, see
[Hacking GitHub with Unicode's dotless 'i'](https://eng.getwisdom.io/hacking-github-with-unicode-dotless-i/).

We want the software to have decent performance for typical users.
Use benchmark testing, e.g., see https://github.com/TomTonic/Set3/blob/main/run_benchmark.txt

### Testing during continuous integration

Note that we use
[various Github Actions](https://github.com/TomTonic/Set3/actions)
for continuous integration tools to check changes
after they are checked into GitHub; if they find problems, please fix them.

## Git commit messages

When writing git commit messages, try to follow the guidelines in
[How to Write a Git Commit Message](https://chris.beams.io/posts/git-commit/):

1.  Separate subject from body with a blank line
2.  Limit the subject line to 50 characters.
    (We're flexible on this, but *do* limit it to 72 characters or less.)
3.  Capitalize the subject line
4.  Do not end the subject line with a period
5.  Use the imperative mood in the subject line (*command* form)
6.  Wrap the body at 72 characters ("<tt>fmt -w 72</tt>")
7.  Use the body to explain what and why vs. how
    (git tracks how it was changed in detail, don't repeat that)

## Reuse (supply chain)

### Requirements for reused components

We prefer reusing components instead of writing lots of code,
but please evaluate all new components before adding them
(including whether or not you need them).
We want to reduce our risks of depending on software that is poorly
maintained or has vulnerabilities (intentional or unintentional).
Furthermore, Set3 provides a concise set of functionality for a
fundamental [abstract data type](https://en.wikipedia.org/wiki/Abstract_data_type).
Consequently, Set3 shall not require a Go project using Set3 to pull a significant
tree of dependencies into the project.

#### License requirements for reused components

All *required* reused software *must* be open source software (OSS).
It's okay to *optionally* use proprietary software and add
portability fixes.

### Updating reused components

Please update only one or few components in each commit, instead of
"everything at once".  This makes debugging problems much easier.
In particular, if we find a problem later, we can
use "git bisect" to easily and quickly find the cause.

