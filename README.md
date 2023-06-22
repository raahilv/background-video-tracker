# background-video-tracker
Scans for videos playing and updates a tracker.

Currently using MPRIS and updates MyAnimeList

Instructions:
Run the authorizer first (MAL_updater_authorizer will create the MAL_token.txt file)
Run it in the background

Goals:
Fully rewrite in Go to not need the python intrepreter
Add support for more sites (e.g. AniList)
Allow blacklisting/whitelisting players
Allowing absolute numbers
Add Windows support?
Add a GUI?
