# vatValidator
Microservice to validate whether a VAT number is a valid German Vat number.

Steps to run vatValidator

- Install go on your machine https://golang.org/doc/install

- Clone this repo in $GOPATH/src. Default GOPATH is your $HOME/go. If you have not exported the GOPATH in bachrc or zshrc, you can see your GOPATH by simply entering ```go env GOPATH``` command in the terminal.

- Go into vatValidator directory and run ```go run main.go```. This will expose api to validate vat number on the endpoint `http://localhost:8001/validate_vat` on your localhost.

- Run index.html in any browser.

