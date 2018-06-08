# pass

[Krypton](https://krypt.co) but for passwords

# Installation

If not already installed, install [Golang](https://golang.org), then run
`go install github.com/conradludgate/pass`

# Usage

## Pairing

To start, you must pair your machine with your phone. YOu can do this by running `pass pair`,
this will present a QR code, scan it with the pass phone app.

## Getting passwords

To get a password, run `pass [name]`. Alternatively, you can use the `--url` flag to choose a password by URL.
Your paired phone will prompt you, asking if you want to continue. 

If access is granted and the password is found, the password is returned, 
otherwise nothing is returned and the program exits with the following error codes:

*	13 Permission Denied
*	61 Password not found

If you recieve a prompt on your phone and you didn't request a password, deny access and check your keys.

## Creating new passwords

If you want to create a new password, use `pass --new`. This will load the new password tool on your phone.