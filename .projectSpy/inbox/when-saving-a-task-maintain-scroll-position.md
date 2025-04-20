When saving a task, maintain scroll position
===



If the task body is longer than the dialog, saving/or pressing ctrl-s causes a reload, which then reopens the dialog task, highlighting the close button. Because the change log is automatically appended to the contents, we'd need a way to set the cursor where it last was - maybe session storage? and number of bytes the cursor is positioned at (that's gonna be quite the JS to write - also needing to take into account that the header may be added on a new task.

---

2025-04-11 21:47	Created task
2025-04-11 21:50	Updated task
2025-04-11 21:50	Updated task