*** Settings ***
Library			    Browser
Library             OperatingSystem

Suite Setup         Task Setup
Suite Teardown      Close Browser  ALL

*** Variables ***
${BROWSER}		%{BROWSER}

*** Test Cases ***
Visit No Hits Page
	New Page			{{ .PublicNoHitsLink }}

	${file}=  Download With Selector   html
	Move File   ${file.saveAs}  ${OUTPUT DIR}${/}downloads${/}no_hits.html

	Take Screenshot     EMBED  fileType=jpeg    quality=25

Visit Multiple Hits Page
	New Page			{{ .PublicMultipleHitsLink }}

	${file}=  Download With Selector   html
	Move File   ${file.saveAs}  ${OUTPUT DIR}${/}downloads${/}multiple_hits.html

	Take Screenshot     EMBED  fileType=jpeg    quality=25

Visit Direct Hit Page
	New Page			{{ .PublicDirectHitLink }}

	${file}=  Download With Selector   html
	Move File   ${file.saveAs}  ${OUTPUT DIR}${/}downloads${/}direct_hit.html

	Take Screenshot     EMBED  fileType=jpeg    quality=25

*** Keywords ***
Task Setup
    Open Browser
	Set Browser Timeout   150s

Open Browser
	New Browser	    chromium  headless=True   downloadsPath=${OUTPUT DIR}
    New Context     acceptDownloads=True

Download With Selector
    [Arguments]    ${selector}=html
	${elem}=          Get Element   ${selector}
    ${file_object}=   Download  ${elem}
    [Return]    ${file_object}
