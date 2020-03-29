package main

import (
	"github.com/fclaude/dicesum/generator"
	"github.com/gin-gonic/gin"
	"os"
)

func ap(c string) string {
	if c == "s" {
		return "'"
	}
	return "'s"
}

func cleanName(name string) string {
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
				<p class="lead">Generate a random dice addition worksheet for your kid!</p>
			</div>
		</div>
	</div>
	<div class="container">
		<div class="row">
			<form action="/generate/" method="POST">
			  <div class="form-group">
				<label for="name">Child Name:</label>
				<input type="text" name="name" class="form-control" id="text">
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
				<p>For more random stuff, you can check <a href="https://fclaude.recoded.cl/">my blog</a> out!</p>
		</div>
	</div>

  </body>
</html>
`