# Example Application

### Description

- Application examples to showcase how to make an application for distributed go web server
- Will be using Go in particular due to the rest of the code base being in the Go language
- Will potentially make a C++ example as well
- Most languages should be able to be used as most langugaes provide support for socket programming 

- For an application to connect to a leader, you must somehow configure application to connect to the leader via the leader's ip
- Depending on how you settup the application server, the application can connect to 1 or more leaders concurrently

### AIM

- showcase how to make an application
- test my code

### Reqiurements

- Any language capable of utilizing network TCP (most languages fulfill this role)
- Use of HTMX (will be communicationg application changes via html)


### Usage

- .html file must include:
  - must be contained within div with attributes of id="app" hx-swap-oob="outerHTML" (can have more)
  - must have an appname

- You can decide whether to connect to a single leader or open requests for leaders to connect to app 
  - must still implement ability for leaders to connect to app

- When responding to request, most usual response is to send back HTML as a string

- Still experimenting with clicks and other forms of inputs other than form requests, but as seen in the aim-trainer example, you can settup an HTTP server within you main app.html and get responses through there. Through this method, the user would never be sending requests through the leader but directly to the application, but the application can still respond through the leader or directly