lets do some refactoring
===

I like where this is going, however the go side of things has gotten pretty messy on the quest for features.

Lets start off with reducing complexity of a few things:

1) duplicate code all over the place - extract out the reused snippets for each endpoint.
2) remove view models / objects - they add another layer frustration getting the data to render.
3) evaluate file, task, and task lane go files.
- things that feel funny are files and tasks being so closely related.
- task lanes setting up watchers should likey be separated.