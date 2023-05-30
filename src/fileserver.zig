const std = @import("std");
const http = std.http;
const fmt = std.fmt;
const fs = std.fs;
const io = std.io;
const heap = std.heap;
const mem = std.mem;

pub fn fileserver(response: *http.Server.Response, allocator: mem.Allocator) void {
    const target = response.request.target;

    //TODO: handle this better
    if (target.len > 128) return;
    const fsDir = "html{s}";

    //Format the target together with the filedir
    var buffer: [128]u8 = undefined;
    const filePath = fmt.bufPrint(&buffer, fsDir, .{target}) catch @panic("URL format buffer too small");
    std.debug.print("{s}\n", .{filePath});

    const cwd = fs.cwd();
    const fileStat = cwd.statFile(filePath) catch |err| switch (err) {
        error.FileNotFound => {
            error404(response, allocator) catch unreachable;
            return;
        },
        else => @panic("Unhandled error"),
    };

    const kind = fs.File.Kind;

    switch (fileStat.kind) {
        kind.file => {
            serveFile(response, filePath, http.Status.ok, allocator) catch |err| {
                std.debug.print("Unhandled error: {}", .{err});
                unreachable;
            };
        },
        kind.directory => {
            std.debug.print("TODO", .{});
        },
        else => {
            error404(response, allocator) catch return;
            return;
        },
    }
}

fn serveFile(response: *http.Server.Response, filePath: []const u8, status: http.Status, allocator: mem.Allocator) !void {
    const cwd = fs.cwd();
    const file = try cwd.openFile(filePath, .{});
    const stat = try file.stat();
    const size = stat.size;

    response.transfer_encoding = .{ .content_length = size };
    try serverStatus(response, status);

    const reader = file.reader();
    const content = try reader.readAllAlloc(allocator, size);

    try response.writer().writeAll(content);

    try response.finish();
}

fn serverStatus(response: *http.Server.Response, status: http.Status) !void {
    response.status = status;
    try response.do();
}

fn error404(response: *http.Server.Response, allocator: mem.Allocator) !void {
    // See if the requester accpets html back, and if so, serve the 404 page!
    const listOptional = response.request.headers.getEntries(allocator, "Accept") catch @panic("Unable to allocate memory to get headers");
    if (listOptional != null) {
        const list = listOptional.?;
        const acceptHTML = mem.containsAtLeast(u8, list[0].value, 1, "text/html");
        if (acceptHTML) {
            serveFile(response, "html/404.html", http.Status.not_found, allocator) catch |erri| switch (erri) {
                error.FileNotFound => {
                    noDefault404(response);
                    return;
                },
                else => {
                    std.debug.print("Unhandled error: {}", .{erri});
                    unreachable;
                },
            };
        }
    } else {
        try serverStatus(response, http.Status.not_found);
    }
}

fn noDefault404(response: *http.Server.Response) void {
    const noDefaultPage =
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

    response.status = http.Status.ok;
    response.transfer_encoding = .{ .content_length = noDefaultPage.len };
    response.headers.append("connection", "close") catch return;
    response.do() catch return;

    response.writeAll(noDefaultPage) catch return;
    response.finish() catch return;
}
