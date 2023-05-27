const std = @import("std");
const http = std.http;
const fmt = std.fmt;
const fs = std.fs;
const io = std.io;
const heap = std.heap;

pub fn fileserver(response: *http.Server.Response) void {
    const target = response.request.target;

    //TODO: handle this better
    if (target.len > 128) return;
    const fsDir = "html{s}";

    //Format the target together with the filedir
    var buffer: [128]u8 = undefined;
    const filePath = fmt.bufPrint(&buffer, fsDir, .{target}) catch @panic("URL format buffer too small");
    std.debug.print("{s}\n", .{filePath});

    const cwd = fs.cwd();
    //404 me later
    const fileStat = cwd.statFile(filePath) catch |err| switch (err) {
        error.FileNotFound => {
            serve404(response) catch unreachable;
            return;
        },
        else => @panic("Unhandled error"),
    };

    const kind = fs.File.Kind;

    switch (fileStat.kind) {
        kind.file => {
            serveFile(response, filePath) catch {
                serve503(response) catch unreachable;
                return;
            };
        },
        kind.directory => {
            unreachable;
        },
        else => serve503(response) catch unreachable,
    }
}

fn serveFile(response: *http.Server.Response, filePath: []const u8) fs.File.OpenError!void {
    const cwd = fs.cwd();
    const file = cwd.openFile(filePath, .{}) catch unreachable;
    const stat = file.stat() catch unreachable;
    const size = stat.size;

    response.status = http.Status.ok;
    response.transfer_encoding = .{ .content_length = size };
    response.headers.append("connection", "close") catch unreachable;
    response.do() catch unreachable;

    var aAlloc = heap.ArenaAllocator.init(heap.page_allocator);

    const reader = file.reader();
    const content = reader.readAllAlloc(aAlloc.allocator(), size) catch unreachable;

    response.writer().writeAll(content) catch unreachable;

    response.finish() catch unreachable;
}

fn serve503(response: *http.Server.Response) !void {
    response.status = http.Status.internal_server_error;
    try response.headers.append("connection", "close");
    try response.do();
    try response.finish();
}

fn serve404(response: *http.Server.Response) !void {
    response.status = http.Status.not_found;
    try response.headers.append("connection", "close");
    try response.do();
    try response.finish();
}
