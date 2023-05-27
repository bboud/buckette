const std = @import("std");
const mem = std.mem;
const http = std.http;

const print = std.debug.print;

//Route function signature
const RouteFnPtr = *const fn (response: *http.Server.Response) void;

//16 bytes
pub const Route = struct {
    method: http.Method,
    route: RouteFnPtr,
    fileserver: bool,
};

pub const Router = struct {
    routes: std.StringHashMap(Route),

    pub fn init(allocator: mem.Allocator) Router {
        return .{ .routes = std.StringHashMap(Route).init(allocator) };
    }

    pub fn deinit(self: *Router) void {
        self.deinit();
    }

    pub fn addRoute(self: *Router, target: []const u8, r: Route) !void {
        try self.routes.put(target, r);
    }

    pub fn route(self: *Router, response: *http.Server.Response) !void {
        const request = response.request;
        const target = request.target;

        //Search on the target skipping the first '/'
        const i = mem.indexOf(u8, target[1..], "/") orelse target.len;

        //In the slice, we return the whole thing upto the index
        const slice = target[0..i];

        // Try to get the route and if not, load index as a fileserver. Failure to get index is a panic.
        const foundRoute: Route = self.routes.get(slice) orelse self.routes.get("/").?;

        if (foundRoute.method == request.method) {
            foundRoute.route(response);
        } else {
            CannedResponses.failed404(response);
        }
    }
};

//TODO: Handle errors correctly
pub const CannedResponses = struct {
    // Returns generic 200 with no data
    pub fn default(response: *http.Server.Response) void {
        response.status = http.Status.ok;
        response.headers.append("connection", "close") catch unreachable;
        response.do() catch unreachable;
        response.finish() catch unreachable;
    }

    pub fn failed404(response: *http.Server.Response) void {
        response.status = http.Status.not_found;
        response.headers.append("connection", "close") catch unreachable;
        response.do() catch unreachable;
        response.finish() catch unreachable;
    }
};
