Stats
===

I really want a stats page that shows some statistics of tasks.

Maybe these could be line charts with x being day and y being number of tickets in done last modified on that day:
Tickets completed per day of week
Tickets completed over the last 7 days

These might be tough, what does created mean? is it affected when a task moves directories?
Tickets created over the last 7 days
Tickets created per day of week

A pie chart for ticket directory distribution would need to figure out a fitting minimal style

Using git support, can possibly gather more-proper statistics. Something along the lines of pulling the last 30 days of commits, finding task changes, and plotting them.
Maybe lines of code changed per task? That might be interesting, but should probably only measure with the task-changing commit (Though if we are writing git commits automatically, this wont work - see git integration task).

```git
git log --since=30.days --follow ./.projectSpy
```
This gets us the last 30 days of commits with changes to the projectSpy directory.
Now we need to determine what changed in each commit. track created files, files that were updated, and count of files in done.

```git
git diff --name-only --diff-filter=ACMRTUXB <commit-hash>
```

This will give us a list of files that were changed in the commit.

If we aggregate the list of files changed in each commit, based on the task directory, and the type of change (create, update, move/copy/rename?)
A data structure like this might be useful:

```json
{
"2025-03-28": {
  "task-directory-1": {
    "create": 1,
    "update": 2,
    "move": 3
  },
  "task-directory-2": {
    "create": 1,
    "update": 2,
    "move": 3
  }
}
```

this works as a line chart with x being day and y being number of task files - with color code for type?




changelog
:2025-03-28 21:58	Created task
:2025-03-28 21:58	Updated task
:2025-03-28 21:59	Updated task
:2025-03-29 23:01	Updated task
:2025-03-29 23:01	Moved task from backlog to inbox
:2025-04-11 21:46	Updated task
: