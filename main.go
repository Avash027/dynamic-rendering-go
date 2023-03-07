package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
)

var reactPort = "4173"

func dynamicRenderer() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request is from a bot
		isBot := false

		// If request is from a bot, use Puppeteer to render the React app
		if isBot {
			// Start the React app in the frontend folder in a separate port

			// Connect to Puppeteer
			ctx, cancel := chromedp.NewContext(context.Background())
			defer cancel()

			// Navigate to the page and wait for it to load
			url := "http://localhost:" + reactPort + c.Request.URL.Path
			var html string
			err := chromedp.Run(ctx,
				chromedp.Navigate(url),
				chromedp.InnerHTML("html", &html, chromedp.NodeVisible, chromedp.ByQuery),
			)
			if err != nil {
				log.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			html = fmt.Sprintf("<html>%s</html>", html)

			// Send back the rendered HTML
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
			c.AbortWithStatus(http.StatusOK)
		}

		// If not a bot, continue to serve the React app as usual
		c.Next()
	}
}

func main() {
	r := gin.Default() // Create a new instance of the gin framework with default middleware already included.

	r.Use(dynamicRenderer()) // Use the dynamicRenderer middleware to render HTML templates using data passed to the endpoint.

	r.Static("/", "./frontend/dist") // Serve static files from the "./frontend/dist" directory when the root endpoint is accessed.

	stdChan := make(chan bool) // Create a new channel to communicate between the main and go routines.

	cmd := exec.Command("npm", "run", "preview") // Create a new command to execute an npm script named "preview".

	cmd.Dir = "./frontend" // Set the working directory for the command to "./frontend".

	cmd.Stderr = os.Stderr // Set the standard error output of the command to the standard error output of the current process.

	stdout, err := cmd.StdoutPipe() // Get a pipe that will be used to read the standard output of the command.

	if err != nil {
		log.Fatal(err)
		return
	}

	if err := cmd.Start(); err != nil { // Start the execution of the command in a new process.
		log.Fatal(err)
		return
	}

	reader := bufio.NewReader(stdout) // Create a new buffered reader for the standard output pipe of the command.

	go func() { // Start a new go routine to read the standard output of the command and signal when the preview is ready.
		for {
			line, _, _ := reader.ReadLine()
			if strings.Contains(string(line), fmt.Sprintf("http://localhost:%s/", reactPort)) {
				stdChan <- true
				break
			}
		}
	}()

	<-stdChan // Wait for the preview to be ready.

	if err := r.Run(":3000"); err != nil { // Start the server and listen for incoming requests.
		log.Fatal(err)
	}

}
