<h2>Under contruction</h2>
This is an experimental website/web server written in Go (golang). It employs the Tiedot database system (which itself was wtitten in Go) and the Martini web package.

# Dependencies

<b>go get github.com/codegangsta/martini</b>

<b>go get github.com/HouzuoGuo/tiedot</b>

#Configuration

Change the DATABASE_DIR setting in <b>/src/tiedotmartini1/config.go</b>, and make sure the new directory exists.

#Building and Running

Add this git's root directory to the GOPATH search path environment variable, then issue <b>go build</b>. This will generate an executable (TiedotMartini1.exe in Windows).

This website is pointed to <b>http://localhost:3000</b>. Set the PORT environment variable to change the port number.    