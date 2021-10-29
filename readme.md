# Lastpass-Search
This is a simple X11-Tool to search and copy a password from the lastpass-vault.

## Usage
1. rename the `auth.example.json` to `auth.json` and insert your credentials.
2. execute `lastpass-search update` to update your local accounts-cache. This will ask you for your OTP-Code if configured
3. execute `lastpass-search search` to open the X11-UI for search.

In the Search you can enter a string that will be used as a filter for the accounts, then select the account using the arrow `up` and `down` keys and copy the password using the `Enter`-key. Pressing `Escape` will close the application.

## Config
You can adjust the colors of the application by setting the values in your xrdb like so:
```
.Xresources
background: #1D2025
background-alt: #2E3933
foreground: #DEE2E5
```
