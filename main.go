// file: main.go
package main
import (
    "github.com/kataras/iris"
    "github.com/kataras/iris/context"
    "image"
    "os"
    "flag"
    "os/exec"
    "bufio"
    "strings"
    "time"
    "math/rand"
    "log"
    "image/png"
    "image/draw"
    "image/color"
    "golang.org/x/image/font"
    "golang.org/x/image/font/gofont/goregular"
    "github.com/golang/freetype"  
    "github.com/golang/freetype/truetype"  
)

var (
    dpi = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
    fontsize = flag.Float64("fontsize", 36, "font size in points")
    linespacing = flag.Float64("linespacing", 1.5, "line spacing (e.g. 2 means double spaced)")
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
    // https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func main() {
    flag.Parse()
    rand.Seed(time.Now().UnixNano())

    app := iris.New()
    // Load all templates from the "./templates" folder
    // where extension is ".html" and parse them
    // using the standard `html/template` package.
    app.RegisterView(iris.HTML("./templates", ".html"))

    // Method:    GET
    // Resource:  http://localhost:8080
    app.Get("/", func(ctx context.Context) {
        fImg, _ := os.Open("guidry.png")
        defer fImg.Close()

        img, _, _ := image.Decode(fImg)
        randomname := "/images/" + randSeq(24) + ".png"

        b := img.Bounds()
        rgba := image.NewRGBA(image.Rect(0, 0, 1594, 887))
        draw.Draw(rgba, rgba.Bounds(), img, b.Min, draw.Src)

        f, err := truetype.Parse(goregular.TTF)
        ftc := freetype.NewContext()
        ftc.SetDPI(*dpi)
        ftc.SetFont(f)
        ftc.SetFontSize(*fontsize)
        ftc.SetClip(rgba.Bounds())
        ftc.SetDst(rgba)
        ftc.SetSrc(image.NewUniform(color.RGBA{200, 100, 0, 255}))
        ftc.SetHinting(font.HintingFull)

        imageleftmargin := 15
        imagetopmargin := 15

        pt := freetype.Pt(imageleftmargin, imagetopmargin+int(ftc.PointToFixed(*fontsize)>>6))
               
        out, err := exec.Command("fortune").Output()
        lines := strings.Split(string(out), "\n")
        log.Println(lines)

        for _, line := range lines {
            _, err = ftc.DrawString(line, pt)
            if err != nil {
                return
            }
            pt.Y += ftc.PointToFixed(*fontsize * *linespacing)
        }
        
        destination, _ := os.Create("." + randomname)
        defer destination.Close()
        buff := bufio.NewWriter(destination)
        png.Encode(buff, rgba)
        buff.Flush()

        ctx.ViewData("image", randomname)
        // Render template file: ./templates/hello.html
        ctx.View("hello.html")
    })

    app.Get("/images/{imageName:string}", func(ctx context.Context) {
        file := "./images/" + ctx.Params().Get("imageName")
        ctx.SendFile(file, "guidry.png")
    }) 

    // Start the server using a network address and block.
    app.Run(iris.Addr(":8080"))
}