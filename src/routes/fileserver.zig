const std = @import("std");
const http = std.http;
const fmt = std.fmt;
const fs = std.fs;

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
    var fileStat = cwd.statFile(filePath) catch unreachable;

    const kind = fs.File.Kind;

    switch (fileStat.kind) {
        kind.File => {
            serveFile(response) catch |err| switch (err) {
                // 503
                else => unreachable,
            };
        },
        kind.Directory => {
            serveFromIndex() catch |err| switch (err) {
                // 404
                fs.File.OpenError.FileNotFound => unreachable,
                // 503
                else => unreachable,
            };
        },
        else => return,
    }
}

fn serveFromIndex() fs.File.OpenError!void {}

fn serveFile(response: *http.Server.Response) fs.File.OpenError!void {
    _ = response;
}

fn serve404(response: *http.Server.Response) !void {
    _ = response;
}
