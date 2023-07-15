# Design

From a high level, the purpose of this library is to provide an additional
abstraction over the Twitch IRC handling implemented in `gempir/go-twitch-irc`.
That gives us message content and metadata, but another layer proves useful
for bot implementations that rely on specific commands/messages to run.

The foundational work for this library was done in `glotchimo/gotato` which
is a chat bot for facilitating hot potato among chatters using a handful of
commands. It has two goroutines, one that listens to IRC for commands, and
another that reacts to those commands. A similar model will be used here,
but in a way that can be extracted and extended to implement arbitrary
functionality.
