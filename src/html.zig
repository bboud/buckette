pub const HTMLHEAD =
    \\<!DOCTYPE html>
    \\<html lang="en">
    \\<head>
    \\    <meta charset="UTF-8">
    \\    <title>{s}</title>
    \\</head>
    \\<body>
    \\    ------------------ <br>
    \\    {s} <br>
    \\    ------------------ <br>
    \\    <ul> 
    \\    
    \\
;

pub const LISTITEM =
    \\  <li>
    \\      <a href="{s}/{s}"> {s} </a>
    \\  </li>
;

pub const HTMLFOOT =
    \\    </ul>
    \\</body>
    \\</html>
;

pub const NODEFAULT =
    \\<!doctype html>
    \\<html lang="en">
    \\  <head>
    \\    <title>No Default Route</title>
    \\  </head>
    \\  <body>
    \\    <main>
    \\       <div>
    \\          <h1>404 Not Found</h1>
    \\          <p>
    \\            The file or directory could not be found.
    \\          </p>
    \\      </div>
    \\    </main>
    \\  </body>
    \\</html>
;
