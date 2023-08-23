**Migrated to https://codeberg.org/anaseto/gofrundis**

Harmonist: Dayoriah Clan Infiltration
-------------------------------------

Harmonist is a stealth coffee-break roguelike game.  The game has a heavy focus
on tactical positioning, light and noise mechanisms, making use of various
terrain types and cones of view for monsters.  Aiming for a replayable
streamlined experience, the game avoids complex inventory management and
character building, relying on items and player adaptability for character
progression.

*Your friend Shaedra got captured by nasty people from the Dayoriah Clan while
she was trying to retrieve a powerful magara artifact that was stolen from the
great magara-specialist Marevor Helith.*

*As a gawalt monkey, you don't understand much why people complicate so much
their lives caring about artifacts and the like, but one thing is clear: you
have to rescue your friend, somewhere to be found in this Underground area
controlled by the Dayoriah Clan.  If what you heard the guards say is true,
Shaedra's imprisoned on the eighth floor.*

*You are small and have good night vision, so you hope the infiltration
will go smoothly...*

Website
-------

[![Introduction Screen](https://download.tuxfamily.org/harmonist/intro-screen-tiles.png)](https://harmonist.tuxfamily.org/index.html)

You can visit the [game's
website](https://harmonist.tuxfamily.org/index.html)
for more informations, tips, screenshots and asciicasts. You will also be able
to play in the browser and download pre-built binaries for the latest release.

Install from Sources
--------------------

In all cases, you need first to perform the following preliminaries:

+ Install the [go compiler](https://golang.org/).
+ Add `$(go env GOPATH)/bin` to your `$PATH` (for example `export PATH="$PATH:$(go env GOPATH)/bin"`).

Harmonist uses the [gruid](https://github.com/anaseto/gruid) library for
grid-based user interfaces, which offers three different rendering drivers:
terminal, graphical SDL2, and browser.

### Terminal (ASCII)

You can build a native ASCII version from source by using the following
command:

	go install

Alternatively, you may use the `go build -o /path/to/bin/harmonist` to put the
resulting binary in a particular place.

The `harmonist` command should now be available (you may have to rename it to
remove the `.git` suffix).

This version uses the [tcell](https://github.com/gdamore/tcell) terminal
library.

### SDL2 (Tiles or ASCII)

You can build a graphical version depending on SDL2 by using the following
command:

	go install --tags sdl

Alternatively, you may use the `go build --tags sdl -o /path/to/bin/harmonist`
to put the resulting binary in a particular place.

This will install the [go-sdl2](https://github.com/veandco/go-sdl2/sdl) Go
bindings for SDL2. You need to install
[SDL2](https://libsdl.org/download-2.0.php) first.

### Browser (Tiles or ASCII)

You can also build a WebAssembly version with:

    GOOS=js GOARCH=wasm go build --tags js -o harmonist.wasm

You can then play by serving a directory containing the wasm file via http. The
directory should contain some other files that you can find in the main
website instance (some HTML and js).

Colors
------

If the default colors do not display nicely on your terminal emulator, you can
check the available options as documented in the manual page.

Check also the other color options.

Documentation
-------------

See the man page harmonist(6) for more information on command line options and use
of the replay file. For example:

    harmonist -r _

launches an auto-replay of your last game.
