const std = @import("std");
const mem = std.mem;
const http = std.http;
const fs = std.fs;

const print = std.debug.print;

//Route function signature
const RouteFnPtr = *const fn (response: *http.Server.Response, allocator: mem.Allocator) void;

pub const Router = struct {
    routes: std.StringHashMap(RouteFnPtr),

    pub fn init(allocator: mem.Allocator) Router {
        return .{ .routes = std.StringHashMap(RouteFnPtr).init(allocator) };
    }

    pub fn deinit(self: *Router) void {
        self.routes.deinit();
    }

    pub fn addRoute(self: *Router, target: []const u8, r: RouteFnPtr) !void {
        try self.routes.put(target, r);
    }

    pub fn route(self: *Router, response: *http.Server.Response, allocator: mem.Allocator) void {
        const request = response.request;
        const target = request.target;

        //Search on the target skipping the first '/'
        const i = mem.indexOf(u8, target[1..], "/") orelse target.len;

        //In the slice, we return the whole thing upto the index
        const slice = target[0..i];

        // Try to get the route and if not, load index as a fileserver. Failure to get index results in canned response
        const foundRoute: RouteFnPtr = self.routes.get(slice) orelse self.routes.get("/") orelse noDefault;

        foundRoute(response, allocator);
    }
};

fn noDefault(response: *http.Server.Response, allocator: mem.Allocator) void {
    _ = allocator;
    const noDefaultPage =
        \\<!doctype html>
        \\<html lang="en">
        \\  <head>
        \\    <title>No Default Route</title>
        \\  </head>
        \\  <body>
        \\    <main>
        \\       <div>
        \\          <h1>Uh Oh!</h1>
        \\          <p>
        \\            It looks like you have no default route on "/".. You are seeing this to indicate that you should add a default route!
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
