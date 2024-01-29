# tmix

A terminal UI for your favorite music provider*. 

*Currently only supports Spotify.


## Why another terminal Spotify player?

While tmix currently only supports Spotify, the goal is not to replicate Spotify's functionality 1:1. Instead, the aim is to provide a unified UI to support any number of music providers simultaneously 
and allow seamlessly switching between them. All without ever leaving your terminal.


## Spotify Setup

To use tmix, we need to connect to the Spotify API.


1. Go to the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Create a new app.
3. Set the redirect URI to http://localhost:3200/callback
4. Check the box for Web API in the section about which API/SDKs you're using.
5. Get the Client ID and Client SecretClient ID and Client Secret.
6. Place these in a Client ID and Client Secret.
7. Create a file at $HOME/.config/tmix/config.toml.
8. Add the following to your config, replacing the Client ID and Secret:


```toml
[providers.spotify]
client-id="MY_CLIENT_ID_HERE"
client-secret="MY_CLIENT_SECRET_HERE"
```
