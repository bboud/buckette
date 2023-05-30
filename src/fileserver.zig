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
            serverStatus(response, http.Status.not_found) catch return;
            response.finish() catch return;
            return;
        },
        else => @panic("Unhandled error"),
    };

    const kind = fs.File.Kind;

    switch (fileStat.kind) {
        kind.file => {
            serveFile(response, filePath, http.Status.ok, allocator) catch {
                serverStatus(response, http.Status.internal_server_error) catch return;
                response.finish() catch return;
                return;
            };
        },
        kind.directory => {
            std.debug.print("TODO", .{});
        },
        else => {
            serveFile(response, "html/404.html", http.Status.not_found, allocator) catch return;
            response.finish() catch return;
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
