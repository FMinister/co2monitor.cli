# CO2 Monitor CLI - README

This script is a simple CLI tool developed in Go that fetches data from a CO2 sensor server (like [https://co2.leyrer.io](https://co2.leyrer.io/swagger/index.html)) and displays it in the console. It uses the [Charm-Bracelet](https://github.com/charmbracelet/bubbletea) library for the user interface and a loading spinner.

## Functionality

The script does the following:

1. Loads environment variables from a `.env` file, including a required API key (`X_API_KEY`) and the API URL (`API_URL`) if present.

2. Starts an infinite loop that periodically (every 100 milliseconds) sends an event to a channel to signal script activity.

3. Periodically (every 60 seconds), checks the CO2 sensor server by sending an HTTP GET request and verifies the server response status.

4. Updates the user interface with the CO2 data and temperature information retrieved from the server or displays an error message if an error occurs.

5. Allows the user to exit the program by pressing any key.

## Usage

To run this script:

1. Ensure you have a `.env` file in your working directory containing the `API_URL` and `X_API_KEY` required for authentication with the CO2 sensor server.

2. Execute the script by compiling and running the `main.go` file or by running it directly if you have Go installed on your system.

3. The program displays the current CO2 data, temperature, and the timestamp of the last measurement. If an error occurs, it will be displayed.

4. Press any key to exit the program.

## Dependencies

The script uses external Go modules that can be downloaded using `go get`:

- Charm-Bracelet for the user interface: `github.com/charmbracelet/bubbletea`
- Loading spinner: `github.com/charmbracelet/bubbles/spinner`
- Lipgloss for text formatting: `github.com/charmbracelet/lipgloss`
- Godotenv for loading environment variables from a `.env` file: `github.com/joho/godotenv`

Ensure that these dependencies are properly installed before running the script.

## Troubleshooting

If you encounter issues or receive error messages when running the script:

- Ensure you have a `.env` file with a valid `API_URL` and `X_API_KEY` in your working directory.
- Verify that the mentioned dependencies are correctly installed.
- Make sure your internet connection is active, and the CO2 sensor server is reachable.
