package cmdlets

import (
	"log"
	"github.com/spf13/cobra"

	"github.com/the-maldridge/yurt-tools/internal/http"
	"github.com/the-maldridge/yurt-tools/internal/http/hello"
)

var (
	httpHelloCmd = &cobra.Command{
		Use: "hello",
		Short: "hello provides a static webserver that generates a hello world and 204 response",
		Long: httpHelloCmdLongDocs,
		Run: httpHelloCmdRun,
	}
	httpHelloCmdLongDocs = `It is often useful to generate a static HTTP response to test load
balancers, associated routing infrastructure, and other components of
your network.  The hello cmdlet launches a webserver with two
endpoints available:

  * /hello/ - Returns the static text "Hello World!"
  * /hello/generate_204 - Returns an empty response with status 204
    (No Content)

Additional endpoints can be added on request.`
)

func init() {
	httpCmd.AddCommand(httpHelloCmd)
}

func httpHelloCmdRun(c *cobra.Command, args []string) {
	srv, err := http.New()
	if err != nil {
		log.Printf("Could not initialize webserver: %s", err)
		return
	}

	h := hello.New()

	srv.Mount("/hello", h.HTTPEntry())

	srv.Serve(":8080")
}
