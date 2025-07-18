package middleware

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
)

// The recover with debug middleware is designed to manage the panic behavior of the framework by catching the panic sequence and restoring normal execution. When this takes place, the middleware will render a built-in go template that displays the panic message and related information. Please see the FrameworkTrace struct for details about what information is displayed in the user interface.
func (m *Middleware) RecovererWithDebug(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The handler is wrapped in a defer since we make a call to recover() to stop the panic sequence and restore normal execution.
		defer func() {

			// recover() will return the value passed to the panic() function (which is what's thrown when an exception occurs). If there was no panic, recover() returns nil.
			if rvr := recover(); rvr != nil {

				// following the mux pattern, do not recover a a http.ErrAbortHandler or log - just panic
				if rvr == http.ErrAbortHandler {
					panic(rvr)
				}

				// Use recover() (above) to catch a panic and then get information about the location where the panic occurred.
				// When calling recover(), it returns the value that was passed to panic(). However, to get more information about the panic, such as the file and line number where it occurred, we use the  runtime/debug package's Stack() function.
				buf := make([]byte, 4096)
				n := runtime.Stack(buf, false)

				trace := FrameworkTrace{
					AdeleVersion: m.FrameworkVersion,
					AppName:      os.Getenv("APP_NAME"),
					RootPath:     m.RootPath,
					StackRaw:     buf[:n],
					FrameCount:   0,
				}

				// Log the details of panic in an error log
				m.Log.Error(string(trace.StackRaw))

				// Get the build information
				build, ok := debug.ReadBuildInfo()
				if ok {
					trace.GoVersion = strings.Replace(build.GoVersion, "go", "", -1)
					trace.PackagePath = build.Path
					trace.MainPath = build.Main.Path
				}

				// get an interface that implements the error type or just get the string representation of the error, if it's not an error, return the error. Otherwise it will print just the error message.
				err, ok := rvr.(error)
				if !ok {
					trace.PanicMessage = rvr.(string)
				} else {
					message := regexp.MustCompile(`(.*)\:(.*)`)
					trace.PanicMessage = strings.TrimSpace(message.ReplaceAllString(err.Error(), "$2"))
					trace.PanicType = strings.TrimSpace(message.ReplaceAllString(err.Error(), "$1"))
				}

				// the 9th item in the stack contains the starting point where the panic occurred
				store := FrameworkTraceEntry{}
				for _, frame := range strings.Split(string(trace.StackRaw), "\n")[9:] {
					ext := regexp.MustCompile(`(.*)\.go:.*`)

					if trace.FrameCount%2 == 0 {
						store.Function = frame
					} else {
						// Line number
						lineNumber := strings.SplitAfter(frame, ":")
						if len(lineNumber) > 1 {
							store.Line = strings.Split(lineNumber[1], " ")[0]
						}

						// File with path
						path := strings.TrimSpace(ext.ReplaceAllString(frame, "$1.go"))
						store.File = path

						// Store
						trace.Stack = append(trace.Stack, store)
					}

					// OG stack formatted
					formattedFrame := strings.TrimSpace(ext.ReplaceAllString(frame, "$1.go"))
					trace.StackFormatted = append(trace.StackFormatted, formattedFrame+"\n")
					trace.FrameCount++
				}

				// Set line number for the stack
				if len(trace.Stack) > 0 {
					trace.PanicLine = trace.Stack[0].Line
				}

				// set the path and name of the file where the panic was triggered
				f := strings.Split(string(trace.StackRaw), "\n")[10]
				fmt.Println("f: " + f)
				trace.FilePath = strings.TrimSpace(strings.Split(f, ":")[0])
				name := regexp.MustCompile(`([^/]+$)`)
				match := name.FindStringSubmatch(f)
				if len(match) > 0 {
					trace.FileName = strings.TrimSpace(strings.Split(match[1], ":")[0])
				}

				fmt.Println("trace.FileNam: " + trace.FileName)
				fmt.Println("trace.FilePath: " + trace.FilePath)

				// We are making an assumption that the file where the panic is thrown is always going to be found and returning if not found. If this is the case we should just return the stack and sidestep printing the source.
				skipSource := false
				if _, err := os.Stat(trace.FilePath); os.IsNotExist(err) {
					skipSource = true
				}

				if !skipSource {

					source, err := os.ReadFile(trace.FilePath)
					if err != nil {
						m.Log.Error(err)
						return
					}

					trace.SourceRaw = string(source)

					for index, formatted := range strings.Split(trace.SourceRaw, "\n") {
						lineNumber := strings.TrimSpace(strconv.Itoa(index + 1))
						trace.SourceFormatted = append(trace.SourceFormatted, lineNumber+" "+formatted+"\n")
					}

					trace.SourceHighlight = strings.Join(trace.SourceFormatted[:], "")
				}

				// Build the HTML to return to the client
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)

				view := getRecoverHTML()
				tmpl := template.Must(template.New("example").Parse(view))

				err = tmpl.Execute(w, trace)
				if err != nil {
					m.Log.Error(err)
					return
				}
			}
		}()
		next.ServeHTTP(w, r)
	})

}

func getRecoverHTML() string {

	return `<html>
        <head>
			<link rel="stylesheet" href="//fonts.googleapis.com/css2?family=Roboto:ital,wght@0,100;0,300;0,400;0,500;0,700;0,900;1,100;1,300;1,400;1,500;1,700;1,900&display=swap">

			<style type="text/css">
				/*!
					Theme: Default
					Description: Original highlight.js style
					Author: (c) Ivan Sagalaev <maniac@softwaremaniacs.org>
					Maintainer: @highlightjs/core-team
					Website: https://highlightjs.org/
					License: see project LICENSE
					Touched: 2021
					*/pre code.hljs{display:block;overflow-x:auto;padding:1em}code.hljs{padding:3px 5px}.hljs{background:#f3f3f3;color:#444}.hljs-comment{color:#697070}.hljs-punctuation,.hljs-tag{color:#444a}.hljs-tag .hljs-attr,.hljs-tag .hljs-name{color:#444}.hljs-attribute,.hljs-doctag,.hljs-keyword,.hljs-meta .hljs-keyword,.hljs-name,.hljs-selector-tag{font-weight:700}.hljs-deletion,.hljs-number,.hljs-quote,.hljs-selector-class,.hljs-selector-id,.hljs-string,.hljs-template-tag,.hljs-type{color:#800}.hljs-section,.hljs-title{color:#800;font-weight:700}.hljs-link,.hljs-operator,.hljs-regexp,.hljs-selector-attr,.hljs-selector-pseudo,.hljs-symbol,.hljs-template-variable,.hljs-variable{color:#ab5656}.hljs-literal{color:#695}.hljs-addition,.hljs-built_in,.hljs-bullet,.hljs-code{color:#397300}.hljs-meta{color:#1f7199}.hljs-meta .hljs-string{color:#38a}.hljs-emphasis{font-style:italic}.hljs-strong{font-weight:700}
			</style>

			<style type="text/css">
				pre > code {
    				font-family: Roboto, ui-sans-serif, system-ui, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
				}

				.hljs{
					background-color: transparent;
				}

				.hljs-number{
					color: #F7B6C2;
					padding-right: 14px;
				}
				.hljs-number.active{
					color: #EB4765;
					font-weight: 700;
				}

				pre{
					position: relative;
    				z-index: 2;
				}

				.spotlight {
					height: 19px;
					width: 100%;
					background: #FDEDF0;
					z-index: -1;
					position: absolute;
					top: 0;
				}
				body {
					background-color: #FDEDF0;
					font-family: Roboto, ui-sans-serif, system-ui, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
  					font-feature-settings: normal;
  					font-variation-settings: normal;
				}
				.container {
					margin: 0 auto;
					width: 90%;
					padding: 15px;
				}
				.col{
					flex-direction: column;
					flex-grow: 1;
					flex-shrink: 1;
				}
				.block__message{
					display: flex;
					flex-wrap: wrap;
					background-color: #EB4765;
					border-radius: 0.25rem;
					height: 125px;
					margin-bottom: 25px;
					padding: 10px 30px;
				}
				.block__stack {
					display: flex;
					height: 83%;
				}
				.block__frame {
					background-color: #FBDAE0;
					border-radius: 0.25rem 0 0 0.25rem;
					flex-basis: 25%;
					padding: 10px 30px;
					text-align: right;
					max-width: 240px;
				}
				.block__code {
					background-color: #FFFFFF;
					border-radius: 0 0.25rem 0.25rem 0;
					padding: 10px 30px;
				}
				.block__code p {
					color: #5C0A19;
					font-size: 16px;
					font-weight: 500;
					text-align: right;
				}
				.block__message h1 {
					flex-grow: 1;
					flex-shrink: 0;
					flex-basis: 70%;
				}
				.block__message p:first-of-type{
					flex-basis: 10%;
				}
				h1 {
					font-size: 1.25rem;
					line-height: 1.75rem;
					font-weight: 900;
					text-transform: uppercase;
					letter-spacing: 0.1em;
					color: #FDEDF0;
					margin-bottom: 25px;
				}
				h2 {
					color: #A5122D;
					font-size: 21px;
					fonts-weight: 700;
					margin: 0;
				}
				p {
		        	font-size: 18px;
					font-weight: 400;
					--tw-text-opacity: 1;
					color: #FDEDF0;
				}
				ul {
					list-style: none;
					padding: 0;
					margin: 0;
					word-break: break-all;
				}
				ul li{
					color: #F7B6C2;
					font-size: 14px;
					font-weight: 400;
					margin: 5px 0 17px 0;
				}
				ul li.active{
					color: #A5122D;
				}
				ul li .function {
					display: block;
					font-weight: 300;
					margin-top: 3px;
				}
				ul li .function .line{
					font-weight: 400;
				}
				.text-capitalize:first-letter {
					text-transform: capitalize;
				}
				.inline-block{
					display: inline-block;
				}
				.text-tiny {
					font-size: 14px;
					font-weight: 300;
				}
				.scroll{
					overflow: hidden;
					overflow-y: scroll;
				}
			</style>
        </head>
        <body>
			<div class="container">
				<div class="block__message">
					<h1>{{ .PanicType }}</h1>
					<p class="text-tiny">Go {{.GoVersion}}</p>
					<p class="text-tiny">Adele {{.AdeleVersion}}</p>
					<p><span class="text-capitalize inline-block">{{.FileName}}</span>: {{ .PanicMessage }}</p>

				</div>
				<div class="block__stack">
					<div class="block__frame col scroll">
						<h2>{{ .FrameCount }} Frames</h2>
						<ul>
						{{ range $index, $frame := .Stack }}
							<li id="frame_{{ len (printf "a%*s" $index "") }}" {{if not $index }}class="active"{{end}}>{{$frame.Function}}
								<span class="function">{{$frame.File}}<span class="line">:{{$frame.Line}}</span></span>
							</li>
						{{ end }}
						</ul>
					</div>
					<div class="block__code col scroll">
						<p id="filePath">{{.FilePath}}</p>
						<pre><code class="language-go">{{.SourceHighlight}}</code>
						<div class="spotlight"></div>
						</pre>

					</div>
				</div>

			</div>

			<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"></script>
            <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/languages/go.min.js"></script>
			<script>
				hljs.addPlugin({
					'after:highlightElement': ({el, result, text}) => {
					console.log(el)
						const elements = document.getElementsByClassName('hljs-number')
						const testDivs = Array.prototype.filter.call(
							elements,
							(e) => {
								if(e.innerHTML == "{{ .PanicLine }}"){
									e.classList.add('active')

									let l = document.getElementsByClassName('hljs-number active')[0]
									let line = l.getBoundingClientRect();
									let sl = document.getElementsByClassName('spotlight')[0]
									let spotlight = sl.getBoundingClientRect()

									document.getElementsByClassName('spotlight')[0].style.top = line.top - spotlight.top +"px"

									e.scrollIntoView({  block: "center"})
								}
							},
						);

					}
				})

			</script>
            <script>hljs.highlightAll();</script>
        </body>
    </html>`
}
