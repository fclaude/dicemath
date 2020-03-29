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
		pdfData, err := generator.GeneratePDF(name)
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

<!-- Global site tag (gtag.js) - Google Analytics -->
<script async src="https://www.googletagmanager.com/gtag/js?id=UA-162008850-1"></script>
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag(){dataLayer.push(arguments);}
  gtag('js', new Date());

  gtag('config', 'UA-162008850-1');
</script>

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
				<p class="lead">Generate a random dice addition worksheet for your child!</p>
			</div>
		</div>
	</div>
	<div class="container">
		<div class="row">
			<form action="/generate/" method="POST">
			  <div class="form-group">
				<label for="name">Child Name (Optional):</label>
				<input type="text" name="name" class="form-control" id="name" maxlength="10">
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

  </body>
</html>
`
