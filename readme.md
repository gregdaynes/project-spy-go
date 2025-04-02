Build & Install
---

Build the project to the local directory

```bash
make build
```

If you want to install the binary to your `$GOPATH/bin` directory
otherwise you can move the binary somewhere else and update the `$PATH`

```bash
make install
```


Setup
---

If you used the install command above, the binary in bin will be called `pspy`

Initialize the project
this will create the following directories
- `backlog`
- `in-progress`
- `done`

```bash
pspy -init
```

You can manually crreate the directories in a `.projectspy`

You can also create two additional directories without a configuration file
- `inbox`
- `blocked`

Run the app which will open a browser window

```bash
pspy
```

Advanced Configuration
---

Configuring lanes is done through a config file located at `.projectSpy/projectspy.json`
This file is not created by default. The following is an example of the config file that defines the default lanes

```json
{
  "lanes": [
    {
      "dir": "inbox",
      "name": "Inbox"
    },
    {
      "dir": "backlog",
      "name": "Backlog"
    },
    {
      "dir": "blocked",
      "name": "Blocked"
    },
    {
      "dir": "in-progress",
      "name": "In Progress"
    },
    {
      "dir": "done",
      "name": "Done"
    }
  ]
}
