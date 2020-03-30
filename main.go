package main

import (
	"github.com/fclaude/dicemath/generator"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func ap(c string) string {
	if c == "s" {
		return "'"
	}
	return "'s"
}

const MaxName = 20

func cleanName(name string) string {
	if len(name) > MaxName {
		log.Printf("Got a name that is too long: %d", len(name))
		name = name[:MaxName]
	}
	res := ""
	for _, c := range name {
		if (c == ' ') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			res += string(c)
		}
	}
	if res == "" {
		return "My"
	}
	return res + ap(string(res[len(res) - 1]))
}

func getOperation(op string) string {
	switch op {
	case "addition": return "+"
	case "multiplication": return "$\\times$"
	case "subtraction": return "-"
	case "division": return "/"
	default: return "?"
	}
}

func main() {
	mode := gin.ReleaseMode
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}
	gin.SetMode(mode)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.Data(200, "text/html", []byte(page))
	})

	r.POST("/generate/", func(c *gin.Context) {
		name := cleanName(c.PostForm("name"))
		operation := getOperation(c.PostForm("operation"))
		pdfData, err := generator.GeneratePDF(name, operation)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.Header("Content-Disposition", `attachment; filename="worksheet.pdf"`)
		c.Data(200, "application/pdf", pdfData)
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

const page = `
<!doctype html>
<html lang="en">
  <head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=yes">

    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">

	<style>
        body {
            padding-top: 50px;
        }

        .title {
            padding: 40px 15px;
            text-align: center;
        }
    </style>
    <title>Dice-math worksheet generator</title>
  </head>
  <body>
   
	<div class="container">
		<div class="row">
			<div class="title">
				<h1>Dice-math worksheet generator</h1>	
				<p class="lead">Generate a random dice-math worksheet for your child!</p>
			</div>
		</div>
	</div>
	<div class="container">
		<div class="row">
			<form action="/generate/" method="POST" id="generate-form">
			  <div class="form-group">
				<label for="name">Child Name (Optional):</label>
				<input type="text" name="name" class="form-control" id="name" maxlength="20">
              	<label for="operation">Operation:</label>
				<select id="operation" name="operation">
				  <option value="addition">Addition</option>
				  <option value="multiplication">Multiplication</option>
				  <option value="subtraction">Subtraction</option>
				  <option value="division">Division</option>
				</select> 
			  </div>
			  <button type="submit" class="btn btn-default btn-primary">Generate</button>
			</form> 
		</div>
	</div>
	<div class="container">
		<div class="row"><p></p></div>
		<div class="row">
				<p>Built to keep my children busy, hope it helps you too :)</p>
		</div>
		<div class="row">
				<p>For more random stuff, you can check out <a href="https://fclaude.recoded.cl/">my blog</a>!</p>
		</div>
	</div>

<script>
// From Google Analytics documentation
(function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
(i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
})(window,document,'script','https://www.google-analytics.com/analytics.js','ga');

ga('create', 'UA-162008850-1', 'auto');
ga('send', 'pageview');

ga('send', 'event', 'Index visit', 'submit', {});

function createFunctionWithTimeout(callback, opt_timeout) {
  var called = false;
  function fn() {
    if (!called) {
      called = true;
      callback();
    }
  }
  setTimeout(fn, opt_timeout || 1000);
  return fn;
}

// Gets a reference to the form element, assuming
// it contains the id attribute "generate-form".
var form = document.getElementById('generate-form');

// Adds a listener for the "submit" event.
form.addEventListener('submit', function(event) {

  // Prevents the browser from submitting the form
  // and thus unloading the current page.
  event.preventDefault();

  // Sends the event to Google Analytics and
  // resubmits the form once the hit is done.
  ga('send', 'event', 'Generate Sheet', document.getElementById('operation').value, {
    hitCallback: createFunctionWithTimeout(function() {
      form.submit();
    })
  });
});

</script>
  </body>
</html>
`
