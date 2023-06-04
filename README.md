# Spotify-Sync
Ever wanted to create a collaborative playlist on spotify with a friend, but you don't both have premium? Perhaps this could help. You can run this command line tool with ```spotify-sync <source-playlist-id> <dest-playlist-id> ...more pairs of playlist ids```. This will copy all new tracks form the source playlist (any public playlist) into the destination playlist (any of your playlists). 

If you have *playlist A* and your friend has *playlist B*, simply get your friend to follow *playlist A*, and add any new songs to *playlist B*. Then run spotify-sync as a cronjob to effectively turn *playlist A* into a collaborative playlist.

## Contributions
Pull requests are welcome!

## Setup

### Binary
TODO

### Spotify API
You will need to login to spotify for developers, and get your client id and client secret.
Then go to the dashboard and find the option to **Create App**. Most of the information here is not important, but makes sure to put ```localhost:9000/callback``` as the **Redirect URI**.

### Environment setup
You will need the following environment variables defined
- SPOTIFY_API_CLIENT_ID
- SPOTIFY_API_CLIENT_SECRET
- SPOTIFY_USERNAME
- SPOTIFY_PASSWORD

### Testing it Out
Move the binary to somewher in your PATH, such as usr/local/bin, then try invoking it from the command line with ```spotify-sync <source-playlist-id> <dest-playlist-id> ...more pairs of playlist ids```

### Setting up a CronJob
By default the crontab won't pick up your environment variables defined in ```.profile``` or equivalent, so you'll need to source these as part of the cronjob. This can be done as below. You can also supply a filepath to output logs to, so that you can check that everything is working.

```
30 19 * * * . <path to file with environment variables>; spotify-sync <source playlist id> <destination playlist id> >> <path to log file>
```
The above example will run the script with one pair of playlists at 19:30 every day and log to a file.


## In the future
Support + instructions for Windows and MacOS