*** Settings ***
Library			Browser

*** Variables ***
# There is a conflict between browser names used by Selenium (using "chrome") and Browser (using "chromium")
# This is why the browser is currently fixed for the Browser Library.
${BROWSER}		%{BROWSER}

*** Test Cases ***
Visit Amazon DE
	IF	"${BROWSER}" == 'chrome'
		New Browser		chromium
	ELSE
		New Browser		${BROWSER}
	END
	New Page			https://amazon.de
	Take Screenshot
