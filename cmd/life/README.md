# Life
This is an example command line program written using the library. It utilises `gocui` (https://github.com/jroimartin/gocui/) to create a CUI.

I never want to see this code again.

## Controls
```bash
life [optional filename] [character used for active cells] [character used for inactive cells]
# eg
life examples/gosper_glider_gun O .
```

| **Command**                | *Description*                                               |
| -------------------------- | ----------------------------------------------------------- |
| `play`                     | Play the animation                                          |
| `pause`                    | Pause the animation                                         |
| `step`                     | Move forward one generation while paused                    |
| `templates`                | View a list of templates, defined in the `examples` folder  |
| `set [template name]`      | View the template                                           |
| `clear`/`reset`            | Clear the grid                                              |
| `on/off [x] [y]`           | Sets that x & y position to 1 or 0, starting from top.      |
| `active/inactive [string]` | Controls the string used to print active and inactive cells |
| `quit`                     | Quit the program                                            |



