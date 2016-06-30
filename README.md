restartd
=======

restartd allow systemd service control through unix socket.

Build
-----

    make

Big picture
-----------

_restartd_ run as root, talks to systemd via DBUS, listen one UNIX socket per user.

_restartctl_ send commands to _restartd_, just like _service_.

Licence
-------

    3 terms BSD licence, © 2016 Mathieu Lecarme, Wilfried Ollivier
