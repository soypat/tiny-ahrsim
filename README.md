# tiny-ahrsim
TinyGo attitude estimation simulation applet.


## Instructions

### Requirements
* Go installed ([golang.org](https://golang.org/))
* git installed ([git-scm.com](https://git-scm.com/downloads))
* TinyGo installed ([tinygo.org](https://tinygo.org/getting-started/install/))
* [gopherjs](https://github.com/gopherjs/gopherjs) installed: once Go installed, run `go get -u github.com/gopherjs/gopherjs` in console

### Steps

1. Clone repository to local computer
    ```shell
    git clone https://github.com/soypat/tiny-ahrsim.git
    ```

2. Change directory to this repo and generate frontend app with gopherjs

    ```shell
    gopherjs build ./graphics/
    ```
    This should create two files: `graphics.js` and `graphics.js.map`.

3. Run the tinygo program (under `tinygo` directory) on your microcontroller of choice, make sure your microcontroller is accesible via USB. Take note of the port which it is available on (usually COM1, COM2, or COM3 on windows). For an Arduino UNO you'd flash the program as follows:
    ```shell
    tinygo flash -target=arduino ./tinygo/main.go
    ```

4. Run Simulation program specifying the port of the USB device

    ```shell
    go run . -ttl=COM3 -p=":8080"
    ```

5. Open `index.html` with a browser and you are set.