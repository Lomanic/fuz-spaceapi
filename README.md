# Fuz SpaceAPI endpoint

Simple project in Go retrieving opening status from [presence button API](https://wiki.fuz.re/doku.php?id=projets:fuz:presence_button) and serving JSON according to [SpaceAPI specs](https://spaceapi.io/docs/).

Read about it on [Fuz wiki](https://wiki.fuz.re/doku.php?id=projets:fuz:spaceapi).

## Configuration

The app gets its configuration from these environment variables:
* `PRESENCEAPI`: URL where the app should get the space opening status (e.g. `PRESENCEAPI=https://presence.fuz.re/api`)
* `SPACEAPI`: JSON string, static information to be served on the endpoint (e.g. `SPACEAPI='{"api":"0.13","space":"FUZ","logo":"https://fuz.re/WWW.FUZ.RE_fichiers/5c02b2a84373a.png","url":"https://fuz.re/","location":{"address":"11-15 rue dela Réunion, Paris 75020, FRANCE","lon":2.40308,"lat":48.85343},"contact":{"email":"","irc":"","ml":"fuz@fuz.re","twitter":"@fuz_re","matrix":"https://matrix.to/#/#fuz_general:matrix.fuz.re"},"issue_report_channels":["ml","twitter"],"state":{"icon":{"open":"https://presence.fuz.re/img","closed":"https://presence.fuz.re/img"},"message":"open under conditions: https://wiki.fuz.re/doku.php?id=map"},"projects":["https://wiki.fuz.re/doku.php?id=projets:fuz:start"]}'`)
