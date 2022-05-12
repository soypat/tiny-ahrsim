# tiny-ahrsim
TinyGo attitude estimation simulation applet. [Here's a youtube video](https://www.youtube.com/watch?v=M0_s6UW86cs&ab_channel=PatricioWhittingslow) of the app in action with a Raspberry Pi Pico.


## Instructions

### Requirements
* Go installed ([golang.org](https://golang.org/))
* git installed ([git-scm.com](https://git-scm.com/downloads))
* TinyGo installed ([tinygo.org](https://tinygo.org/getting-started/install/))
* [gopherjs](https://github.com/gopherjs/gopherjs) installed: once Go installed, run `go install github.com/gopherjs/gopherjs@latest` in console

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

## ⚠️ Notice! ⚠️
The [gopherjs bindings](https://github.com/soypat/gthree) for `three.js` have been archived!

If you are starting a new project consider using the [WASM bindings](https://github.com/soypat/three)!

