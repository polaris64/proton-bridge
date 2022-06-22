# Fork of ProtonMail Bridge to add support for `mu4e` tags as Labels

This is a custom fork of ProtonMail Bridge. It adds support for tracking message tags added by `mu` and `mu4e` as ProtonMail Labels.

**NOTE**: This is an experiment at this stage and should be treated as early alpha quality. Only very minor testing has been done and no optimisation has been performed. Therefore please only use this fork at your own risk; I am not responsible for any data loss or corruption that occurs due to its use.

## Explanation of `mu4e` tags and ProtonMail labels

When a message is tagged in `mu4e` a new "X-Keywords" header is added by default (although this can be changed). `mu` can use these headers to filter messages easily by those tags.

With normal IMAP servers, this header can be saved as part of the message which means that tags are stored on the server. However, once a message is stored in ProtonMail it can no longer be modified; if the message with "X-Keywords" is pushed to ProtonMail via IMAP then this header will not be saved. Please see [this issue](https://protonmail.uservoice.com/forums/284483-protonmail/suggestions/43440678-support-for-custom-imap-headers-e-g-x-keywords) for further details.

ProtonMail however offers support for Labels, which function very similarly to these tags, however these are not directly visible per message via ProtonMail Bridge. Instead, Labels act like IMAP folders, meaning that messages appear to be duplicated in one or more of these Label folders.

## Changes in this fork

This forked version of ProtonMail bridge currently has two modifications: -

  - When fetching a message via IMAP, it will now create an "X-Keywords" header automatically with the names of all Labels assigned to the message; and
  - When storing a message, it will extract the tags from an "X-Keywords" header (if present), convert them to Label IDs (creating new Labels if they are missing), and then assign those Labels to the message.

The benefit of this is that messages can be tagged in `mu4e` and those tags will appear as Labels in ProtonMail. Furthermore, those tags/Labels will be present on all machines which are using this forked version of ProtonMail Bridge to fetch messages.

# Proton Mail Bridge and Import Export app
Copyright (c) 2022 Proton AG

This repository holds the Proton Mail Bridge and the Proton Mail Import-Export applications.
For a detailed build information see [BUILDS](./BUILDS.md).
The license can be found in [LICENSE](./LICENSE) file, for more licensing information see [COPYING_NOTES](./COPYING_NOTES.md).
For contribution policy see [CONTRIBUTING](./CONTRIBUTING.md).


## Description Bridge
Proton Mail Bridge for e-mail clients.

When launched, Bridge will initialize local IMAP/SMTP servers and render 
its GUI.

To configure an e-mail client, firstly log in using your Proton Mail credentials. 
Open your e-mail client and add a new account using the settings which are 
located in the Bridge GUI. The client will only be able to sync with 
your Proton Mail account when the Bridge is running, thus the option 
to start Bridge on startup is enabled by default.

When the main window is closed, Bridge will continue to run in the
background.

More details [on the public website](https://protonmail.com/bridge).

## Description Import-Export app
Proton Mail Import-Export app for importing and exporting messages.

To transfer messages, firstly log in using your Proton Mail credentials.
For import, expand your account, and pick the address to which to import
messages from IMAP server or local EML or MBOX files. For export, pick
the whole account or only a specific address. Then, in both cases,
configure transfer rules (match source and target mailboxes, set time
range limits and so on) and hit start. Once the transfer is complete,
check the results.

More details [on the public website](https://protonmail.com/import-export).

The Import-Export app is developed in separate branch `master-ie`.

## Launchers
Launchers are binaries used to run the Proton Mail Bridge or Import-Export apps.

Official distributions of the Proton Mail Bridge and Import-Export apps contain
both a launcher and the app itself. The launcher is installed in a protected
area of the system (i.e. an area accessible only with admin privileges) and is
used to run the app. The launcher ensures that nobody tampered with the app's
files by verifying their signature using a hardcoded public key. App files are
placed in regular userspace and are signed by Proton's private key. This
feature enables the app to securely update itself automatically without asking
the user for a password.

## Keychain
You need to have a keychain in order to run the Proton Mail Bridge. On Mac or
Windows, Bridge uses native credential managers. On Linux, use `secret-service` freedesktop.org API
(e.g. [Gnome keyring](https://wiki.gnome.org/Projects/GnomeKeyring/))
or
[pass](https://www.passwordstore.org/). We are working on allowing other secret
services (e.g. KeepassXC), but for now only gnome-keyring is usable without
major problems.


## Environment Variables

### Bridge application
- `BRIDGESTRICTMODE`: tells bridge to turn on `bbolt`'s "strict mode" which checks the database after every `Commit`. Set to `1` to enable.

### Dev build or run
- `APP_VERSION`: set the bridge app version used during testing or building
- `PROTONMAIL_ENV`: when set to `dev` it is not using Sentry to report crashes
- `VERBOSITY`: set log level used during test time and by the makefile

### Integration testing
- `TEST_ENV`: set which env to use (fake or live)
- `TEST_ACCOUNTS`: set JSON file with configured accounts
- `TAGS`: set build tags for tests
- `FEATURES`: set feature dir, file or scenario to test


## Files
### Database
The database stores metadata necessary for presenting messages and mailboxes to an email client:
- Linux: `~/.cache/protonmail/bridge/<cacheVersion>/mailbox-<userID>.db` (unless `XDG_CACHE_HOME` is set, in which case that is used as your `~`)
- macOS: `~/Library/Caches/protonmail/bridge/<cacheVersion>/mailbox-<userID>.db`
- Windows: `%LOCALAPPDATA%\protonmail\bridge\<cacheVersion>\mailbox-<userID>.db`

### Preferences
User preferences are stored in json at the following location:
- Linux: `~/.config/protonmail/bridge/prefs.json`
- macOS: `~/Library/ApplicationSupport/protonmail/bridge/prefs.json`
- Windows: `%APPDATA%\protonmail\bridge\prefs.json`

### IMAP Cache
The currently subscribed mailboxes are held in a json file:
- Linux: `~/.cache/protonmail/bridge/<cacheVersion>/user_info.json` (unless `XDG_CACHE_HOME` is set, in which case that is used as your `~`)
- macOS: `~/Library/Caches/protonmail/bridge/<cacheVersion>/user_info.json`
- Windows: `%LOCALAPPDATA%\protonmail\bridge\<cacheVersion>\user_info.json`

### Lock file
Bridge utilises an on-disk lock to ensure only one instance is run at once. The lock file is here: 
- Linux: `~/.cache/protonmail/bridge/<cacheVersion>/bridge.lock` (unless `XDG_CACHE_HOME` is set, in which case that is used as your `~`)
- macOS: `~/Library/Caches/protonmail/bridge/<cacheVersion>/bridge.lock`
- Windows: `%LOCALAPPDATA%\protonmail\bridge\<cacheVersion>\bridge.lock`

### TLS Certificate and Key
When bridge first starts, it generates a unique TLS certificate and key file at the following locations:
- Linux: `~/.config/protonmail/bridge/{cert,key}.pem` (unless `XDG_CONFIG_HOME` is set, in which case that is used as your `~/.config`)
- macOS: `~/Library/ApplicationSupport/protonmail/bridge/{cert,key}.pem`
- Windows: `%APPDATA%\protonmail\bridge\{cert,key}.pem`

