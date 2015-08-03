# Quobar -- X11 status bar

> The Quo is not Status

Quobar is a fairly minimalist status bar. It's greatest difference
from `dzen2`, `xmobar` and such is that it is completely and from the
ground up explicitly graphical, and aims to take advantage of this
fact.

## Why use raw X11?

Because QML wouldn't let me set strut hints.

The widgets are still abstracted away from the details, and use Go
idioms for drawing.
