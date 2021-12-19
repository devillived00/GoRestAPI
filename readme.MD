App runs on port :8080

A propert GET operation is "/random/mean?requests={r}&length={l}" where "r" is an argument for number of requests and "l" is an argument for number of items in the list.

Docker instruction:

  1) Build an image - "docker build . -t ExampleName"
  2) Run app - "docker run -p 8080:8080 -tid ExampleName"
