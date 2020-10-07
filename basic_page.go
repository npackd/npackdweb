package main

import (
	"net/http"
	"context"
)

// Basisklasse
type BasicPage struct {
	// Kontext
	c context.Context

	// Anfrage
	r *http.Request

	// Antwort
	w http.ResponseWriter

	// Ausgabe
	b *HTMLBuffer

	// Fehlermeldung
  pageError string
  
  info string

  params map[string]string
}

func (page *BasicPage) setNoCacheHeaders() {
	page.w.Header()["Cache-Control"] = []string{"no-cache, must-revalidate"}
	page.w.Header()["Expires"] = []string{"Fri, 01 Jan 1990 00:00:00 GMT"}
}

func (page *BasicPage) beforeBody(title string) {
	b := page.b
	b.B("html")
	b.B("head")
	b.B("title", title)
	b.E("title")
	b.Unencoded(`<!-- Latest compiled and minified CSS -->
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">

<meta name="viewport" content="width=device-width, initial-scale=1">

<!-- Optional theme -->
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap-theme.min.css" integrity="sha384-rHyoN1iRsVXV4nD0JutlnGaslCJuC7uwjduW9SVrLvRYooPp2bWYgmgJQIXwl/Sp" crossorigin="anonymous">

<link rel="manifest" href="/static/manifest.json">

<!-- iOS still needs this, 27, May 2018 -->
<link rel="apple-touch-icon" sizes="192x192" href="/static/apple-touch-icon180.png">

<!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.12.4/jquery.min.js"></script>

<script src="/static/autosize.min.js"></script>

<!-- Latest compiled and minified JavaScript -->
<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js" integrity="sha384-Tc5IQib027qvyjSMfHjOMaLkfuWVxZxUPnCJA7l2mCWNIpG9mGCD8wGNIcPD7Txa" crossorigin="anonymous"></script>
<LINK REL="stylesheet" TYPE="text/css" HREF="/static/Style.css" MEDIA="all" />
<link rel="stylesheet" type="text/css" href="/static/jquery.cookiebar.css" />
<link rel="icon" type="image/png" href="/static/logo.png" />
`)
	b.E("head")
	b.B("body")

	page.menu()

	b.B("div", "class", "content")
}

func (page *BasicPage) afterBody() {
	page.b.E("div") // class=content
	page.footer()
	page.b.Unencoded(`  
    <script type="text/javascript" src="/static/jquery.cookiebar.js"></script>
                
	
    <script type="text/javascript">
        $(document).ready(function(){
                $.cookieBar({
                });
        });
        autosize($('textarea'));
    </script>
    <script type="text/javascript">
        var deferredPrompt = null;
        window.addEventListener('beforeinstallprompt', function(e) {
            e.preventDefault();
            deferredPrompt = e;
            document.getElementById("promptInstallLink").style.display = "inline";
        });

        function promptInstall() {
            if (deferredPrompt !== null) {
                deferredPrompt.prompt();
                deferredPrompt.userChoice.then(function(choiceResult) {
                    if (choiceResult.outcome === "accepted") {
                        
                    }
                    deferredPrompt = null;
                });
            } else {
                alert("Creating a desktop icon is not supported on your system.");
            }
        }
    </script>
</body>
		</html>`)
}

func (page *BasicPage) footer() {
	page.b.Unencoded(`<div class="nw-footer">
<a href="https://github.com/qedus/nds">nds</a> was used and is released under the Apache 3.0 license<br>
<a href="https://github.com/captaincodeman/datastore-mapper">Datastore Mapper</a> was used<br>
<a href="http://www.glyphicons.com">GLYPHICONS FREE</a> were used and are released under the Creative Commons Attribution 3.0 Unported License (CC BY 3.0)<br>
<a href="http://www.fatcow.com/free-icons">Farm Fresh Icons by Fatcow Web Hosting</a> were used and are released under the Creative Commons Attribution 4.0 (CC BY 4.0)<br>
<a href="http://getbootstrap.com/">Bootstrap</a> was used and is released under MIT License<br>
<a href="http://www.primebox.co.uk/projects/jquery-cookiebar/">jQuery CookieBar Plugin</a> was used and is released under the Creative Commons Attribution 3.0 Unported License (CC BY 3.0)<br>
<a href="https://github.com/mmcdole/gofeed">gofeed</a> was used and is released under the MIT license<br>

All data is stored on Google servers. See their <a href="https://www.google.com/intl/en/policies/terms/">Terms of Service</a><br>
This program was created by Tim Lebedkov (Tim dot Lebedkov at g m a i l dot com)
</div>`)
}

func (page *BasicPage) menu() {
  menuTemplate.Execute(page.w, "test")
  
  /* TODO
  return NWUtils.tmpl("Menu.html", "admin",
  NWUtils.isAdminLoggedIn() ? "true" : null, "login",
  MyPage.getLoginHeader(request), "searchForm",
  needsSearchFormInTheMenu() ? "true" : null);	
  */
}


