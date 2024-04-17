const std = @import("std");
const http = std.http;
const fmt = std.fmt;
const fs = std.fs;
const io = std.io;
const heap = std.heap;
const mem = std.mem;

const html = @import("html.zig");
const config = @import("config.zig");

pub fn fileserver(response: *http.Server.Response, allocator: mem.Allocator) void {
    const target = response.request.target;

    var buffer: [128]u8 = undefined;
    const filePath = fmt.bufPrint(&buffer, "{s}{s}", .{ config.WWWPATH, target }) catch unreachable;

    std.log.debug("filepath: {s}", .{filePath});

    const cwd = fs.cwd();
    const fileStat = cwd.statFile(filePath) catch |err| switch (err) {
        error.FileNotFound => {
            error404(response, allocator) catch unreachable;
            return;
        },
        else => {
            std.debug.print("Unhandled error: {}", .{err});
            unreachable;
        },
    };

    const kind = fs.File.Kind;

    switch (fileStat.kind) {
        kind.file => {
            serveFile(response, filePath, http.Status.ok, allocator) catch |err| {
                switch (err) {
                    error.ConnectionResetByPeer => {
                        std.log.warn("Connection reset by peer: {}", .{err});
                    },
                    else => {
                        std.debug.print("Unhandled error: {}", .{err});
                        unreachable;
                    },
                }
            };
        },
        kind.directory => {
            serveIndexOrDirectory(response, filePath, allocator);
        },
        else => {
            error404(response, allocator) catch return;
            return;
        },
    }
}

fn serveIndexOrDirectory(response: *http.Server.Response, filePath: []const u8, allocator: mem.Allocator) void {
    var buffer: [128]u8 = undefined;
    const indexPath = fmt.bufPrint(&buffer, "{s}/index.html", .{filePath}) catch unreachable;

    const cwd = fs.cwd();
    _ = cwd.statFile(indexPath) catch |err| switch (err) {
        error.FileNotFound => {
            serveDirectory(response, filePath, allocator) catch |erri| {
                std.log.err("Error Serving Directory: {}", .{erri});
            };
            return;
        },
        else => unreachable,
    };

    serveFile(response, indexPath, http.Status.ok, allocator) catch |err| {
        std.log.err("Error serving file: {}", .{err});
    };
}

fn serveDirectory(response: *http.Server.Response, filePath: []const u8, allocator: mem.Allocator) !void {
    var target = response.request.target;

    // Remove the leading '/' from the path
    target = target[1..];

    // Define the title for the page as root if target doesn't exist.
    var title: []const u8 = undefined;
    if (target.len == 0) {
        title = "root";
    } else {
        title = target;
    }

    std.log.debug("TARGET {s}", .{target});

    // Build the html to be sent for the directory listings
    var buffer: [html.HTMLHEAD.len + 128]u8 = undefined;
    var htmlBuilder = fmt.bufPrint(&buffer, html.HTMLHEAD, .{ title, title }) catch unreachable;

    const cwd = fs.cwd();
    var dirIter = cwd.openIterableDir(filePath, .{}) catch unreachable;
    defer dirIter.close();

    var iterator = dirIter.iterate();

    while (try iterator.next()) |item| {
        // Just reuse the buffer
        var bufferi: [128]u8 = undefined;
        const li = fmt.bufPrint(&bufferi, html.LISTITEM, .{ target, item.name, item.name }) catch @panic("format buffer too small!");
        const slice = [_][]const u8{ htmlBuilder, li };
        htmlBuilder = mem.concat(allocator, u8, &slice) catch @panic("Cannot concat!");
    }

    const slice = [_][]const u8{ htmlBuilder, html.HTMLFOOT };
    htmlBuilder = mem.concat(allocator, u8, &slice) catch @panic("Cannot concat!");

    response.status = http.Status.ok;
    response.transfer_encoding = .{ .content_length = htmlBuilder.len };
    try response.headers.append("connection", "close");
    try response.do();

    try response.writeAll(htmlBuilder);
    try response.finish();
}

fn serveFile(response: *http.Server.Response, filePath: []const u8, status: http.Status, allocator: mem.Allocator) !void {
    const cwd = fs.cwd();
    const file = try cwd.openFile(filePath, .{});
    defer file.close();

    const stat = try file.stat();
    const size = stat.size;

    response.transfer_encoding = .{ .content_length = size };
    try serverStatus(response, status);

    const reader = file.reader();

    const content = try reader.readAllAlloc(allocator, size);
    defer allocator.free(content);

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
    response.status = http.Status.ok;
    response.transfer_encoding = .{ .content_length = html.NODEFAULT.len };
    response.headers.append("connection", "close") catch return;
    response.do() catch return;

    response.writeAll(html.NODEFAULT) catch return;
    response.finish() catch return;
}
