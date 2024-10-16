# Contributing to eebus-go

First off, thanks for taking the time to contribute! ❤️

All types of contributions are encouraged and valued. See the [Table of Contents](#table-of-contents) for different ways to help and details about how this project handles them. Please make sure to read the relevant section before making your contribution. It will make it a lot easier for us maintainers and smooth out the experience for all involved. The community looks forward to your contributions. 🎉

## Table of Contents

- [Discussions and Questions](#discussions-and-questions)
- [Bug Reports](#bug-reports)
- [New Feature Requests](#new-feature-requests)
- [Issue Tracker](#issue-tracker)
- [Pull Requests](#pull-requests)
- [Styleguides](#styleguides)

## Discussions and Questions

For discussions, questions, feature requests, or ideas, [start a new discussion](https://github.com/enbility/eebus-go/discussions/new) in the eebus-go repository under the Discussions tab.

Before you ask a question, it is best to search for existing [Discussions](https://github.com/enbility/eebus-go/discussions) that might help you. In case you have found a suitable issue and still need clarification, you can write your question in this issue. It is also advisable to search the internet for answers first.

If you then still feel the need to ask a question and need clarification, we recommend the following:

- Open an [Discussion](https://github.com/enbility/eebus-go/discussions/new/choose).
- Provide as much context as you can about what you're running into.

## Bug Reports

### Before Submitting a Bug Report

A good bug report shouldn't leave others needing to chase you up for more information. Therefore, we ask you to investigate carefully, collect information and describe the issue in detail in your report. Please complete the following steps in advance to help us fix any potential bug as fast as possible.

- Make sure that you are using the latest version.
- Determine if your bug is really a bug and not an error on your side.
- To see if other users have experienced (and potentially already solved) the same issue you are having, check if there is not already a bug report existing for your bug or error in the [bug tracker](https://github.com/enbility/eebus-go/issues?q=label%3Abug).
- Collect information about the bug:
  - Stack trace
  - If possible and relevant, the `trace` log of the SHIP and SPINE communication
  - Possibly your input and the output
  - Can you reliably reproduce the issue? And can you also reproduce it with older versions?

### How Do I Submit a Good Bug Report?

> You must never report security related issues, vulnerabilities or bugs including sensitive information to the issue tracker, or elsewhere in public. Instead sensitive bugs must be sent by email to <mail@andreaslinde.de>.

We use GitHub issues to track bugs and errors. If you run into an issue with the project:

- Open an [Issue](https://github.com/enbility/eebus-go/issues/new). (Since we can't be sure at this point whether it is a bug or not, we ask you not to talk about a bug yet and not to label the issue.)
- Explain the behavior you would expect and the actual behavior.
- Please provide as much context as possible and describe the *reproduction steps* that someone else can follow to recreate the issue on their own. This usually includes your code. For good bug reports you should isolate the problem and create a reduced test case.

## New Feature Requests

This section guides you through submitting an enhancement suggestion for eebus-go, **including completely new features and minor improvements to existing functionality**. Following these guidelines will help maintainers and the community to understand your suggestion and find related suggestions.

### Before Submitting an Enhancement

- Make sure that you are using the latest version.
- Check the api interfaces carefully and find out if the functionality is already covered.
- Perform a [discussion search](https://github.com/enbility/eebus-go/discussions) to see if the enhancement has already been suggested. If it has, add a comment to the existing discussion instead of opening a new one.
- Find out whether your idea fits with the scope and aims of the project. It's up to you to make a strong case to convince the project's developers of the merits of this feature. Keep in mind that we want features that will be useful to the majority of our users and not just a small subset. If you're just targeting a minority of users, consider writing an add-on/plugin library.

### How Do I Submit a Good Enhancement Suggestion?

Enhancement suggestions are tracked as [Discussions](https://github.com/enbility/eebus-go/discussions).

- Use a **clear and descriptive title** for the issue to identify the suggestion.
- Provide a **step-by-step description of the suggested enhancement** in as many details as possible.
- **Describe the current behavior** and **explain which behavior you expected to see instead** and why. At this point you can also tell which alternatives do not work for you.
- **Explain why this enhancement would be useful**. You may also want to point out the other projects that solved it better and which could serve as inspiration.

## Issue Tracker

The [Issue Tracker](https://github.com/enbility/eebus-go/issues) is used to discuss bug fixes and details for improvements once they agreed on as [Discussions](https://github.com/enbility/eebus-go/discussions).

- Use a **clear and descriptive title** for the issue to identify the suggestion.
- The issue should describe the intent of the change.
- Provide a link to the discussion (if available) that this issue is based on
- Provide a **step-by-step description of the suggested enhancement** in as many details as possible.

## Pull Requests

> ### Legal Notice
>
> When contributing to this project, you must agree that you have authored 100% of the content, that you have the necessary rights to the content and that the content you contribute may be provided under the project license.

We recommend creating your pull-request as a "draft" and to commit early and often so the community can give you feedback at the beginning of the process as opposed to asking you to change hours of hard work at the end.

- Describe the contribution. First document which issue number was fixed. Then describe the contribution.
- Associated coverage unit tests should be provided.
- Provide the expected behavior changes of the pull request.
- Provide any additional context if applicable.
- Verify that the PR passes all workflow checks. If you expect some of these checks to fail. Please note it in the Pull Request text or comments.

### Fix for whitespace, format code, or make a purely cosmetic patch?

Changes that are cosmetic in nature and do not add anything substantial to the stability, functionality, or testability of eebus-go will generally not be accepted.

### Do you want to add a new feature or change an existing one?

- Suggest your change in the [Discussions](https://github.com/enbility/eebus-go/discussions)
- Do not open an issue on GitHub until you have collected positive feedback about the change. GitHub issues are primarily intended for bug reports, fixes, and enhancement detail discussions.

## Styleguides

- The project uses [golangci-lint](https://golangci-lint.run)
- It is a goal to cover as much code as possible with at least one test case and don't decrease test coverage noticably

## Attribution

This guide is based on the **contributing.md**. [Make your own](https://contributing.md/)!
